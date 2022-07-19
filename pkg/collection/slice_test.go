package collection

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	s := NewSlice[int]()
	s.Append(1, 2, 3)

	s.Insert(2, 1000, 1001, 1002)
	s.Insert(2, 99)

	// s.Remove(3)

	s.RemoveEle(1000)

	s.Append(123)

	fmt.Println(s, cap(s), len(s))
}
