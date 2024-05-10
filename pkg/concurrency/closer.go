package concurrency

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var ErrManagerAlreadyClosed = errors.New("runner manager already closed")

type RunnerCloserManager struct {
	mgr *RunnerManager

	closers []func() error

	retErr error

	fatalShutdownFn func()

	closeFatalShutdown chan struct{}

	running atomic.Bool
	closing atomic.Bool
	closed  atomic.Bool
	closeCh chan struct{}
	stopped chan struct{}
}

func (c *RunnerCloserManager) WithFatalShutdown(fn func()) {
	c.fatalShutdownFn = fn
}

func NewRunnerCloserManager(gracePeriod time.Duration, runners ...Runner) *RunnerCloserManager {
	c := &RunnerCloserManager{
		mgr:                NewRunnerManager(runners...),
		stopped:            make(chan struct{}),
		closeCh:            make(chan struct{}),
		closeFatalShutdown: make(chan struct{}),
	}

	if gracePeriod == 0 {
		return c
	}

	c.fatalShutdownFn = func() {
		panic("\"Graceful shutdown timeout exceeded, forcing shutdown\"")
	}

	_ = c.AddCloser(func() {
		// log.Debugf("Graceful shutdown timeout: %s", *gracePeriod)

		t := time.NewTimer(gracePeriod)
		defer t.Stop()

		select {
		case <-t.C:
			c.fatalShutdownFn()
		case <-c.closeFatalShutdown:
		}
	})

	return c
}

// Add implements RunnerManager.Add.
func (c *RunnerCloserManager) Add(runner ...Runner) error {
	if c.running.Load() {
		return ErrManagerAlreadyStarted
	}

	return c.mgr.Add(runner...)
}

// AddCloser adds a closer to the list of closers to be closed once the main
// runners are done.
func (c *RunnerCloserManager) AddCloser(closers ...any) error {
	if c.closing.Load() {
		return ErrManagerAlreadyClosed
	}

	c.mgr.mu.Lock()
	defer c.mgr.mu.Unlock()

	var errs []error
	for _, cl := range closers {
		switch v := cl.(type) {
		case io.Closer:
			c.closers = append(c.closers, v.Close)
		case func(context.Context) error:
			c.closers = append(c.closers, func() error {
				// We use a background context here since the fatalShutdownFn will kill
				// the program if the grace period is exceeded.
				return v(context.Background())
			})
		case func() error:
			c.closers = append(c.closers, v)
		case func():
			c.closers = append(c.closers, func() error {
				v()
				return nil
			})
		default:
			errs = append(errs, fmt.Errorf("unsupported closer type: %T", v))
		}
	}

	return errors.Join(errs...)
}

func (c *RunnerCloserManager) Run(ctx context.Context) error {
	if !c.running.CompareAndSwap(false, true) {
		return ErrManagerAlreadyStarted
	}

	// Signal the manager is stopped.
	defer close(c.stopped)

	// If the main runner has at least one runner, add a closer that will
	// close the context once Close() is called.
	if len(c.mgr.runners) > 0 {
		c.mgr.Add(func(ctx context.Context) error {
			select {
			case <-ctx.Done():
			case <-c.closeCh:
			}
			return nil
		})
	}

	errCh := make(chan error, len(c.closers))
	go func() {
		errCh <- c.mgr.Run(ctx)
	}()

	rErr := <-errCh

	c.mgr.mu.Lock()
	defer c.mgr.mu.Unlock()
	c.closing.Store(true)

	errs := make([]error, len(c.closers)+1)
	errs[0] = rErr

	for _, closer := range c.closers {
		go func(closer func() error) {
			errCh <- closer()
		}(closer)
	}

	// Wait for all closers to be done.
	for i := 1; i < len(c.closers)+1; i++ {
		// Close the fatal shutdown goroutine if all closers are done. This is a
		// no-op if the fatal go routine is not defined.
		if i == len(c.closers) {
			close(c.closeFatalShutdown)
		}
		errs[i] = <-errCh
	}

	c.retErr = errors.Join(errs...)

	return c.retErr
}

// Close will close the main runners and then the closers.
func (c *RunnerCloserManager) Close() error {
	if c.closed.CompareAndSwap(false, true) {
		close(c.closeCh)
	}
	// If the manager is not running yet, we stop immediately.
	if c.running.CompareAndSwap(false, true) {
		close(c.stopped)
	}
	c.WaitUntilShutdown()
	return c.retErr
}

// WaitUntilShutdown will block until the main runners and closers are done.
func (c *RunnerCloserManager) WaitUntilShutdown() {
	<-c.stopped
}
