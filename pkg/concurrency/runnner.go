package concurrency

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrManagerAlreadyStarted = errors.New("runner manager already started")

type Runner func(context.Context) error

type RunnerManager struct {
	mu      sync.Mutex
	runners []Runner
	running atomic.Bool
}

// NewRunnerManager creates a new RunnerManager.
func NewRunnerManager(runners ...Runner) *RunnerManager {
	return &RunnerManager{
		runners: runners,
	}
}

func (r *RunnerManager) Add(runners ...Runner) error {
	if r.running.Load() {
		return ErrManagerAlreadyStarted
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runners = append(r.runners, runners...)

	return nil
}

func (r *RunnerManager) Run(ctx context.Context) error {
	if !r.running.CompareAndSwap(false, true) {
		return ErrManagerAlreadyStarted
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)
	for _, runner := range r.runners {
		go func(runner Runner) {
			defer cancel()
			rErr := runner(ctx)
			if rErr != nil && !errors.Is(rErr, context.Canceled) {
				errCh <- rErr
				return
			}
			errCh <- nil
		}(runner)
	}

	errObjs := make([]error, len(r.runners))
	for i := 0; i < len(r.runners); i++ {
		errObjs[i] = <-errCh
	}

	return errors.Join(errObjs...)
}
