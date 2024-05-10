package threading

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kanengo/goutil/pkg/gls"
	"github.com/kanengo/goutil/pkg/threading/gopool"
)

func testPanicFunc() {
	panic("test")
}

func TestPool(t *testing.T) {
	p := gopool.NewPool("test", 100, gopool.NewConfig())
	var n int32
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&n, 1)
		}, nil)
	}
	wg.Wait()
	if n != 2000 {
		t.Error(n)
	}
}

const benchmarkTimes = 10000

func TestPoolPanic(t *testing.T) {

	p := gopool.NewPool("test", 100, gopool.NewConfig())
	// p.SetPanicHandler(func(ctx context.Context) {
	// 	fmt.Println("catch panic poll default")
	// })
	p.Go(testPanicFunc, nil)
	// p.Go(testPanicFunc, func(ctx context.Context) {
	// 	fmt.Println("catch panic task default")
	// })
	time.Sleep(time.Second)
}

func DoCopyStack(a, b int) int {
	if b < 100 {
		return DoCopyStack(0, b+1)
	}
	return 0
}

func testFunc() {
	DoCopyStack(0, 0)
}

func BenchmarkPool(b *testing.B) {
	config := gopool.NewConfig()
	config.ScaleThreshold = 1
	p := gopool.NewPool("benchmark", int32(runtime.GOMAXPROCS(0)), config)
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			p.Go(func() {
				testFunc()
				wg.Done()
			}, nil)
		}
		wg.Wait()
	}
}

func BenchmarkGo(b *testing.B) {
	var wg sync.WaitGroup
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			go func() {
				testFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TestGoInheritable(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	gls.InheritableThreadLocalSet[int64]("go-1", 10)
	Go(func() {
		defer wg.Done()
		v, _ := gls.InheritableThreadLocalGet[int64]("go-1")
		fmt.Println("1-", v)
		gls.InheritableThreadLocalSet[int64]("go-2", 20)
		Go(func() {
			defer wg.Done()
			v1, _ := gls.InheritableThreadLocalGet[int64]("go-1")
			v2, _ := gls.InheritableThreadLocalGet[int64]("go-2")
			fmt.Println("2-", v1, v2)
		})
	})

	wg.Wait()
}
