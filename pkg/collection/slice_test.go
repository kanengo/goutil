package collection

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	s := NewSlice[int]()
	sb := s
	s.Append(1, 2, 3, 4, 5, 6, 7)

	s.Insert(2, 1000, 1001, 1002)
	s.Insert(2, 99)

	// s.Remove(3)

	// s.Remove(s.FindIndex(func(ele int) bool {
	// 	return ele == 1000
	// }))

	s.Append(123)

	fmt.Println(s.Source(), cap(*s.s), len(*s.s))

	si := s.Map(func(val int) int {
		// fmt.Println("m", val)
		return val + 1
	}).Filter(func(val int) bool {
		// fmt.Println("f", val)
		if val == 1001 || val == 4 {
			return true
		}
		return false
	})

	fmt.Println(si.Slice().Source())

	sum := si.Reduce(func(previousValue any, val int) any {
		sum := previousValue.(int)
		sum += val
		return sum
	}, 0)

	_ = sum

	fmt.Println(sb.Source())
}
