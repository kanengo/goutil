package usync

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExclusiveCallDo(t *testing.T) {
	g := NewSingleFlight[string]()
	v, err := g.Do("key", func() (string, error) {
		return "bar", nil
	})
	if got, want := fmt.Sprintf("%v (%T)", v, v), "bar (string)"; got != want {
		t.Errorf("Do = %v; want %v", got, want)
	}
	if err != nil {
		t.Errorf("Do error = %v", err)
	}
}

func TestExclusiveCallDoErr(t *testing.T) {
	g := NewSingleFlight[string]()
	someErr := errors.New("some error")
	v, err := g.Do("key", func() (string, error) {
		return "", someErr
	})
	if err != someErr {
		t.Errorf("Do error = %v; want someErr", err)
	}
	if v != "" {
		t.Errorf("unexpected non-nil value %#v", v)
	}
}

func TestExclusiveCallDoDupSuppress(t *testing.T) {
	g := NewSingleFlight[string]()
	c := make(chan string)
	var calls int32
	fn := func() (string, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, err := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
}
