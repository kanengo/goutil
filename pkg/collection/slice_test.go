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

	// s.Remove(s.FindIndex(func(ele int) bool {
	// 	return ele == 1000
	// }))

	s.Append(123)

	fmt.Println(s, cap(s), len(s))

	si := s.Map(func(val int) int {
		fmt.Println("m", val)
		return val + 1
	}).Filter(func(val int) bool {
		fmt.Println("f", val)
		if val == 1001 || val == 4 {
			return true
		}
		return false
	})

	fmt.Println(si.Slice())
}
