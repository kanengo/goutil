package gls

import (
	"github.com/timandy/routine"
)

type localStorage map[string]any

func (ls localStorage) Clone() any {
	if len(ls) == 0 {
		return nil
	}

	bak := make(localStorage, len(ls))
	for k, v := range ls {
		bak[k] = v
	}

	return bak
}

type ThreadLocal[T any] routine.ThreadLocal[T]

var threadLocal = routine.NewThreadLocal[localStorage]()
var inheritableThreadLocal = routine.NewInheritableThreadLocal[localStorage]()

func ThreadLocalSet[T any](key string, value T) {
	var storage localStorage
	v := threadLocal.Get()
	if v == nil {
		storage = make(localStorage)
	} else {
		storage = v
	}
	storage[key] = value
	threadLocal.Set(storage)
}

func ThreadLocalGet[T any](key string) (ret T, ok bool) {
	var storage localStorage
	tv := threadLocal.Get()
	if tv == nil {
		storage = make(localStorage)
	} else {
		storage = tv
	}

	v, ok := storage[key]
	if ok {
		ret = v.(T)
	}

	return ret, ok
}

func InheritableThreadLocalSet[T any](key string, value T) {
	var storage localStorage
	v := inheritableThreadLocal.Get()
	if v == nil {
		storage = make(localStorage)
	} else {
		storage = v
	}
	storage[key] = value
	inheritableThreadLocal.Set(storage)
}

func InheritableThreadLocalGet[T any](key string) (ret T, ok bool) {
	tv := inheritableThreadLocal.Get()
	if tv == nil {
		return ret, false
	}
	storage := tv

	v, ok := storage[key]
	if ok {
		ret = v.(T)
	}

	return ret, ok
}
