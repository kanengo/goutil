package threading

import (
	"context"
	"fmt"
	"sync"

	"github.com/kanengo/goutil/pkg/threading/gopool"
	"github.com/kanengo/goutil/pkg/utils"
)

type Runnable func()

func Go(fn Runnable) {
	GoSafe(context.Background(), fn)
}

func GoSafe(ctx context.Context, fn Runnable) {
	GoSafeWithPanicHandler(ctx, fn, nil)
}

func GoSafeWithPanicHandler(ctx context.Context, fn func(), panicHandler func(context.Context)) {
	copied := createInheritedMap()

	go func() {
		defer func() {
			utils.CheckGoPanic(ctx, panicHandler)
		}()
		t := currentThread(copied != nil)
		if t == nil {
			defer func() {
				t = currentThread(false)
				if t != nil {
					t.threadLocals = nil
					t.inheritableThreadLocals = nil
				}
			}()
			fn()
		} else {
			backup := t.inheritableThreadLocals
			defer func() {
				t.threadLocals = nil
				t.inheritableThreadLocals = backup
			}()
			t.threadLocals = nil
			t.inheritableThreadLocals = copied
			fn()
		}
	}()
}

var defaultPool gopool.Pool
var poolMap sync.Map

func init() {
	defaultPool = gopool.NewPool("gopool.DefaultPool", 10000, &gopool.Config{
		ScaleThreshold: 1,
	})
}

func GoPool(fn func()) {
	defaultPool.Go(fn, nil)
}

func GoPoolWithPanicHandler(fn func(), panicHandler func(context.Context)) {
	defaultPool.Go(fn, panicHandler)
}

func GoCtxPool(ctx context.Context, fn func()) {
	defaultPool.GoCtx(ctx, fn, nil)
}

func GoCtxPoolWithPanicHandler(ctx context.Context, fn func(), panicHandler func(context.Context)) {
	defaultPool.GoCtx(ctx, fn, panicHandler)
}

func SetDafultPoolPanicHandler(panicHandler func(context.Context)) {
	defaultPool.SetPanicHandler(panicHandler)
}

func RegisterPool(p gopool.Pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if loaded {
		return fmt.Errorf("name: %s already registerd", p.Name())
	}
	return nil
}

func GetPool(name string) gopool.Pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}

	return p.(gopool.Pool)
}
