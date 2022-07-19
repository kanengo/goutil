package collection

import (
	"fmt"
)

type Slice[T comparable] []T

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

func (sp *Slice[T]) RemoveEle(remveEle T) {
	flag := false
	for i := 0; i < len(*sp)-1; i++ {
		if (*sp)[i] == remveEle {
			flag = true
		}
		if flag {
			(*sp)[i] = (*sp)[i+1]
		}
	}
	(*sp) = (*sp)[:len(*sp)-1]
}

func NewSlice[T comparable](args ...int) Slice[T] {
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
