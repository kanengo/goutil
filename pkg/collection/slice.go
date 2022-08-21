package collection

import (
	"fmt"
	"sync"
)

type Slice[T any] []T

func (s *Slice[T]) At(index int) (ret T, err error) {
	if !s.checkBound(index) {
		return ret, fmt.Errorf("the index is out of the slice bound")
	}
	ret = (*s)[index]

	return
}

func (s *Slice[T]) Append(ele ...T) {
	if len(*s)+len(ele) > cap(*s) { //扩容
		newCap := len(*s) + len(ele)
		newCap += newCap >> 1
		newSlice := make(Slice[T], 0, newCap)
		*s = append(newSlice, (*s)...)
	}
	*s = append(*s, ele...)
}

func (s *Slice[T]) checkBound(index int) bool {
	if index < 0 || index >= len(*s) {
		return false
	}
	return true
}

func (s *Slice[T]) Update(index int, ele T) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	(*s)[index] = ele
	return nil
}

func (s *Slice[T]) Insert(index int, ele ...T) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}
	if cap(*s) >= len(*s)+len(ele) {
		tmp := (*s)[len(*s)-1]
		eleLen := len(ele)
		for i := 0; i < eleLen; i++ {
			*s = append(*s, tmp)
		}
		for i := len(*s) - eleLen - 1; i >= index; i-- {
			(*s)[i+eleLen] = (*s)[i]
		}
		elei := 0
		for i := index; i < index+eleLen; i++ {
			(*s)[i] = ele[elei]
			elei += 1
		}
	} else {
		part1 := (*s)[:index]
		part2 := (*s)[index:]
		newCap := len(*s) + len(ele)
		newCap += newCap >> 1
		newSlice := make(Slice[T], 0, newCap)
		*s = append(newSlice, part1...)
		*s = append(*s, ele...)
		*s = append(*s, part2...)
	}

	return nil
}

func (s *Slice[T]) Remove(index int) error {
	if !s.checkBound(index) {
		return fmt.Errorf("the index is out of the slice bound")
	}

	for i := index; i < len(*s)-1; i++ {
		(*s)[i] = (*s)[i+1]
	}

	(*s) = (*s)[:len(*s)-1]

	return nil
}

func (s *Slice[T]) FindIndex(fn func(ele T) bool) int {
	for i, ele := range *s {
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
	return len(*s)
}

func (s *Slice[T]) Clear() {
	(*s) = (*s)[:0]
}

func (s *Slice[T]) Source() []T {
	return *s
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
	ret := make(Slice[T], sl, sc)

	return &ret
}

func (s *Slice[T]) Map(f func(val T) T) Iterator[T] {
	si := SliceIterator[T]{
		slice: s,
	}
	si.Map(f)
	return &si
}

func (s *Slice[T]) Filter(f func(val T) bool) Iterator[T] {
	si := SliceIterator[T]{
		slice: s,
	}
	si.Filter(f)
	return &si
}

func (s *Slice[T]) Reduce(f func(previousValue any, val T) any, initialValue any) any {
	si := SliceIterator[T]{
		slice: s,
	}
	return si.Reduce(f, initialValue)
}

func (s *Slice[T]) Range(fn func(val T, index int) bool) {
	for i, ele := range *s {
		if !fn(ele, i) {
			break
		}
	}
}

type SliceIterator[T any] struct {
	next  *iterTypeFunc[T]
	slice *Slice[T]
}

func (si *SliceIterator[T]) Range(fn func(val T) bool) {
	for _, ele := range *si.slice {
		iterF := si.next
		for iterF != nil {
			switch iterF.typ {
			case IterTypeFilter:
				if iterF.filterF(ele) {
					goto NEXT
				}
			case IterTypeMap:
				ele = iterF.mapF(ele)
			}
			iterF = iterF.next
		}
		if !fn(ele) {
			break
		}
	NEXT:
	}
}

func (si *SliceIterator[T]) Slice() *Slice[T] {
	ret := NewSlice[T]()
	si.Range(func(val T) bool {
		ret.Append(val)
		return true
	})
	return ret
}

func (si *SliceIterator[T]) insertFunc(itf *iterTypeFunc[T]) {
	if si.next == nil {
		si.next = itf
	} else {
		si.next.next = itf
	}
}

func (si *SliceIterator[T]) Map(f func(val T) T) Iterator[T] {
	itf := &iterTypeFunc[T]{
		typ:  IterTypeMap,
		mapF: f,
	}

	si.insertFunc(itf)

	return si
}

func (si *SliceIterator[T]) Filter(f func(val T) bool) Iterator[T] {
	itf := &iterTypeFunc[T]{
		typ:     IterTypeFilter,
		filterF: f,
	}
	si.insertFunc(itf)
	return si
}

func (si *SliceIterator[T]) Reduce(f func(previousValue any, val T) any, initialValue any) any {
	ret := initialValue
	si.Range(func(val T) bool {
		ret = f(ret, val)
		return true
	})
	return ret
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
