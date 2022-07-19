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

	fmt.Println(m.Get(1), m.Get(2), m.Get(3, "default"), m.Get(6), m.Get(7))

}
