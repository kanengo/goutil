package collection

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewMapSet[string]()
	s.Add("a")
	s.Add("b")
	s.Add("c")
	s.Add("d")
	s.Remove("c")

	s.Add("e")

	pop := s.Pop()
	r := s.Rand()
	fmt.Println(s.IsMember("a"), s.IsMember("b"), s.IsMember("c"), s.IsMember("d"),
		s.IsMember("e"), s.IsMember(pop), s.IsMember(r))

	fmt.Println(s.Members())
}
