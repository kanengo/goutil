package collection

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap[int, string]()

	m[1] = "1"
	m[2] = "2"
	m[3] = "3"
	m[4] = "4"
	m[5] = "5"
	m.Add(6, "6")

	m.Remove(3)

	m.Update(1, "111")

	fmt.Println(m.Keys(), m.Values())

	// fmt.Println(m.Get(1), m.Get(2), m.Get(3, "default"), m.Get(6), m.Get(7))

}

func TestSyncMap(t *testing.T) {
	m := NewSyncMap[int, string]()
	m.Add(1, "111")
	m.Add(2, "222")
	m.Add(3, "333")

	m.Remove(3)

	m.Add(4, "444")

	a := m
	a.Range(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})
	fmt.Println(a.Keys(), a.Values())
}
