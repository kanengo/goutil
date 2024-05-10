package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

var workerPool sync.Pool

func init() {
	workerPool.New = newWorker
}

type worker struct {
	pool *pool
}

func newWorker() any {
	return &worker{}
}

func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}
			if t == nil {
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}
			w.pool.taskLock.Unlock()
			func() {
				var panicHandler func(context.Context)
				if t.panicHandler != nil {
					panicHandler = t.panicHandler
				} else if w.pool.panicHandler != nil {
					panicHandler = w.pool.panicHandler
				}
				defer func() {
					if r := recover(); r != nil {
						panicHandler(context.TODO())
					}
				}()
				t.fn()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) close() {
	w.pool.decWorkerCount()
}

func (w *worker) clear() {
	w.pool = nil
}

func (w *worker) Recycle() {
	w.clear()
	workerPool.Put(w)
}
