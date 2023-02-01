package gls

import (
	"fmt"
	"sync"
	"testing"
)

func TestThreadLocalGetSet(t *testing.T) {
	ThreadLocalSet[string]("go-1", "hello-1")
	ThreadLocalSet[string]("go-1", "world-1")
	ret, ok := ThreadLocalGet[string]("go-1")
	fmt.Println(ret, ok)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		ThreadLocalSet[string]("go-2", "hello-2")
		ThreadLocalSet[string]("go-2", "world-2")
		ret, ok := ThreadLocalGet[string]("go-1")
		fmt.Println(ret, ok)
	}()

	go func() {
		defer wg.Done()
		ThreadLocalSet[string]("go-3", "hello-3")
		ThreadLocalSet[string]("go-3", "world-3")
		ret, ok := ThreadLocalGet[string]("go-3")
		fmt.Println(ret, ok)
	}()

	wg.Wait()
}
