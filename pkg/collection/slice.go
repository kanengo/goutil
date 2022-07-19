package collection

import (
	"fmt"
	"sync"
)

type Slice[T any] []T

func (sp *Slice[T]) Append(ele ...T) {
	if len(*sp)+len(ele) > cap(*sp) { //扩容
		newCap := len(*sp) + len(ele)
		newCap += newCap >> 1
		newSlice := make(Slice[T], 0, newCap)
		*sp = append(newSlice, (*sp)...)
	}
	*sp = append(*sp, ele...)
}

func (s Slice[T]) checkBound(index int) bool {
	if index < 0 || index >= len(s) {
		return false
	}
	return true
}

func (s Slice[T]) Update(index int, ele T) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	s[index] = ele
	return nil
}

func (sp *Slice[T]) Insert(index int, ele ...T) error {
	if !sp.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	if cap(*sp) >= len(*sp)+len(ele) {
		tmp := (*sp)[len(*sp)-1]
		eleLen := len(ele)
		for i := 0; i < eleLen; i++ {
			*sp = append(*sp, tmp)
		}
		for i := len(*sp) - eleLen - 1; i >= index; i-- {
			(*sp)[i+eleLen] = (*sp)[i]
		}
		elei := 0
		for i := index; i < index+eleLen; i++ {
			(*sp)[i] = ele[elei]
			elei += 1
		}
	} else {
		part1 := (*sp)[:index]
		part2 := (*sp)[index:]
		newCap := len(*sp) + len(ele)
		newCap += newCap >> 1
		newSlice := make(Slice[T], 0, newCap)
		*sp = append(newSlice, part1...)
		*sp = append(*sp, ele...)
		*sp = append(*sp, part2...)
	}

	return nil
}

func (sp *Slice[T]) Remove(index int) error {
	if !sp.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}

	for i := index; i < len(*sp)-1; i++ {
		(*sp)[i] = (*sp)[i+1]
	}

	(*sp) = (*sp)[:len(*sp)-1]

	return nil
}

func (sp *Slice[T]) FindIndex(fn func(ele T) bool) int {
	for i, ele := range *sp {
		if fn(ele) {
			return i
		}
	}

	return -1
}

func (sp *Slice[T]) Length() int {
	return len(*sp)
}

func (sp *Slice[T]) Clear() {
	(*sp) = (*sp)[:0]
}

func NewSlice[T any](args ...int) Slice[T] {
	sl := 0
	sc := 0
	if len(args) >= 2 {
		sl = args[0]
		sc = args[1]
	} else if len(args) >= 1 {
		sl = args[0]
	}
	return make(Slice[T], sl, sc)
}

// type SliceMapIter[T any] struct {
// 	slice *Slice[T]
// 	fn    func(T) T
// }

// func (smi *SliceMapIter[T]) Range(fn func(val T) bool) {
// 	for _, ele := range *smi.slice {
// 		if !fn(smi.fn(ele)) {
// 			break
// 		}
// 	}
// }

// func (smi *SliceMapIter[T]) Slice() Slice[T] {
// 	slice := NewSlice[T](0, len(*smi.slice))
// 	smi.Range(func(val T) bool {
// 		slice.Append(val)
// 		return true
// 	})

// 	return slice
// }

// type SliceFilterIter[T any] struct {
// 	slice  *Slice[T]
// 	filter func(T) bool
// }

// func (smi *SliceFilterIter[T]) Range(fn func(val T) bool) {
// 	for _, ele := range *smi.slice {
// 		if smi.filter != nil && smi.filter(ele) {
// 			continue
// 		}
// 	}
// }

// func (smi *SliceFilterIter[T]) Slice() Slice[T] {
// 	slice := NewSlice[T](0, len(*smi.slice))
// 	smi.Range(func(val T) bool {
// 		slice.Append(val)
// 		return true
// 	})

// 	return slice
// }

// func (smi *SliceFilterIter[T]) Map(f func(val T) T) Iterator[T] {

// }

// func (sp *Slice[T]) Filter(fn func(T) bool) Iterator[T] {
// 	return &SliceFilterIter[T]{
// 		slice:  sp,
// 		filter: fn,
// 	}
// }
// func (sp *Slice[T]) Map(fn func(T) T) Iterator[T] {
// 	return &SliceMapIter[T]{
// 		slice: sp,
// 		fn:    fn,
// 	}
// }

type SlicePool[T any] struct {
	*sync.Pool
}

func (p *SlicePool[T]) Put(s *Slice[T]) {
	s.Clear()
	p.Pool.Put(p)
}

func (p *SlicePool[T]) Get() (s *Slice[T]) {
	return p.Pool.Get().(*Slice[T])
}

func NewSlicePool[T any]() SlicePool[T] {
	return SlicePool[T]{
		&sync.Pool{
			New: func() any {
				slice := NewSlice[T]()
				return &slice
			},
		},
	}
}

func GetSliceInPool[T any]() *Slice[T] {
	return nil
}
