package gls

import (
	"fmt"
	"sync"
	"testing"
)

func TestGoId(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			id := GoId()
			fmt.Println(fmt.Sprintf("%d:", i), id)
		}(i)
	}
	wg.Wait()
}
