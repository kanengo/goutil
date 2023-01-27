package collection

import (
	"fmt"
)

type Slice[T any] struct {
	s *[]T
	*defaultIterator[T]
}

func (s *Slice[T]) At(index int) (ret T, err error) {
	if !s.checkBound(index) {
		return ret, fmt.Errorf("the index is out of the slice bound")
	}
	ret = (*s.s)[index]

	return
}

func (s *Slice[T]) Append(ele ...T) {
	if len(*s.s)+len(ele) > cap(*s.s) { //扩容
		newCap := len(*s.s) + len(ele)
		newCap += newCap >> 1
		newSlice := make([]T, 0, newCap)
		*s.s = append(newSlice, *s.s...)
	}
	*s.s = append(*s.s, ele...)
}

func (s *Slice[T]) checkBound(index int) bool {
	if index < 0 || index >= len(*s.s) {
		return false
	}
	return true
}

func (s *Slice[T]) Update(index int, ele T) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	(*s.s)[index] = ele
	return nil
}

func (s *Slice[T]) Insert(index int, ele ...T) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	if cap(*s.s) >= len(*s.s)+len(ele) {
		tmp := (*s.s)[len(*s.s)-1]
		eleLen := len(ele)
		for i := 0; i < eleLen; i++ {
			*s.s = append(*s.s, tmp)
		}
		for i := len(*s.s) - eleLen - 1; i >= index; i-- {
			(*s.s)[i+eleLen] = (*s.s)[i]
		}
		elei := 0
		for i := index; i < index+eleLen; i++ {
			(*s.s)[i] = ele[elei]
			elei += 1
		}
	} else {
		part1 := (*s.s)[:index]
		part2 := (*s.s)[index:]
		newCap := len(*s.s) + len(ele)
		newCap += newCap >> 1
		newSlice := make([]T, 0, newCap)
		*s.s = append(newSlice, part1...)
		*s.s = append(*s.s, ele...)
		*s.s = append(*s.s, part2...)
	}

	return nil
}

func (s *Slice[T]) Remove(index int) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}

	for i := index; i < len(*s.s)-1; i++ {
		(*s.s)[i] = (*s.s)[i+1]
	}

	*s.s = (*s.s)[:len(*s.s)-1]

	return nil
}

func (s *Slice[T]) FindIndex(fn func(ele T) bool) int {
	for i, ele := range *s.s {
		if fn(ele) {
			return i
		}
	}

	return -1
}

func (s *Slice[T]) Find(fn func(ele T) bool) bool {
	return s.FindIndex(fn) >= 0
}

func (s *Slice[T]) Length() int {
	return len(*s.s)
}

func (s *Slice[T]) Clear() {
	*s.s = (*s.s)[:0]
}

func (s *Slice[T]) Source() []T {
	return *s.s
}

func (s *Slice[T]) Range(fn func(val T) bool) {
	for _, ele := range *s.s {
		if !fn(ele) {
			break
		}
	}
}

func NewSlice[T any](args ...int) *Slice[T] {
	sl := 0
	sc := 0
	if len(args) >= 2 {
		sl = args[0]
		sc = args[1]
	} else if len(args) >= 1 {
		sl = args[0]
	}
	s := make([]T, sl, sc)

	slice := &Slice[T]{
		s: &s,
	}

	slice.defaultIterator = &defaultIterator[T]{
		Ranger: slice,
	}

	return slice
}

func WrapSlice[T any](s []T) *Slice[T] {
	slice := NewSlice[T]()
	slice.s = &s

	return slice
}
