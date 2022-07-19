package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	Name() string
	// SetCap sets the goroutine capacity of the pool.
	SetCap(cap int32)
	// Go executes f.
	Go(fn func(), panicHandler func(context.Context))
	// GoCtx executes f and accepts the context.
	GoCtx(ctx context.Context, f func(), panicHandler func(context.Context))
	// SetPanicHandler sets the panic handler.
	SetPanicHandler(f func(context.Context))
	// WorkerCount returns the number of running workers
	WorkerCount() int32
}

var taskPool sync.Pool

func init() {
	taskPool.New = newTask
}

type task struct {
	ctx          context.Context
	fn           func()
	panicHandler func(context.Context)
	next         *task
}

func newTask() any {
	return &task{}
}

func (t *task) Clear() {
	t.ctx = nil
	t.fn = nil
	t.panicHandler = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.Clear()
	taskPool.Put(t)
}

type taskList struct {
	sync.Mutex
	head *task
	tail *task
}

type pool struct {
	name string
	// capacity of the pool, the maximum number of goroutines that are actually working
	cap int32

	config *Config

	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	workerCount int32

	panicHandler func(context.Context)
}

func NewPool(name string, cap int32, config *Config) Pool {
	p := &pool{
		name:   name,
		cap:    cap,
		config: config,
	}

	return p
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(fn func(), panicHandler func(context.Context)) {
	p.GoCtx(context.Background(), fn, panicHandler)
}

func (p *pool) GoCtx(ctx context.Context, fn func(), panicHandler func(context.Context)) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.fn = fn
	t.panicHandler = panicHandler
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()
	atomic.AddInt32(&p.taskCount, 1)

	if (atomic.LoadInt32(&p.taskCount) >= p.config.ScaleThreshold &&
		p.workerCount < atomic.LoadInt32(&p.cap)) || p.workerCount == 0 {
		p.incWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

// SetPanicHandler the func here will be called after the panic has been recovered.
func (p *pool) SetPanicHandler(f func(context.Context)) {
	p.panicHandler = f
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) decWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}
