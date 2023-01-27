package collection

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap[int, string]()

	m.Add(1, "1")
	m.Add(2, "2")
	m.Add(3, "3")
	m.Add(4, "4")
	m.Add(5, "5")
	m.Add(6, "6")

	m.Remove(3)

	m.Update(1, "111")

	fmt.Println(m.Keys().Source(), m.Values())
	s := m.Filter(func(val Pairs[int, string]) bool {
		return val.Key() != 4
	}).Map(func(val Pairs[int, string]) Pairs[int, string] {
		val.SetValue(val.Value() + "0")
		return val
	}).Slice().Source()
	fmt.Println(s)
	// fmt.Println(m.Get(1), m.Get(2), m.Get(3, "default"), m.Get(6), m.Get(7))

}
