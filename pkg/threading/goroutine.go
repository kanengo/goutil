package threading

import (
	"context"
	"fmt"
	"sync"

	"github.com/kanengo/goutil/pkg/threading/gopool"
	"github.com/kanengo/goutil/pkg/utils"
)

func Go(fn func()) {
	go func() {
		utils.CheckGoPanic(context.Background(), nil)
		fn()
	}()
}

func GoSafe(ctx context.Context, fn func()) {
	go func() {
		utils.CheckGoPanic(ctx, nil)
		fn()
	}()
}

func GoSafeWithPanicHandler(ctx context.Context, fn func(), panicHandler func(context.Context)) {
	go func() {
		utils.CheckGoPanic(ctx, panicHandler)
		fn()
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
