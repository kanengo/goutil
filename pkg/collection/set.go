package collection

type MapSet[T comparable] map[T]struct{}

func (s MapSet[T]) Add(member T) {
	s[member] = struct{}{}
}

func (s MapSet[T]) Remove(member T) {
	delete(s, member)
}

func (s MapSet[T]) IsMember(member T) (ok bool) {
	_, ok = s[member]
	return
}

func (s MapSet[T]) Len() int {
	return len(s)
}

func (s MapSet[T]) Pop() (ret T) {
	ret = s.Rand()
	delete(s, ret)
	return
}

func (s MapSet[T]) Rand() (ret T) {
	for k := range s {
		ret = k
	}
	return
}

func (s MapSet[T]) Members() *Slice[T] {
	members := NewSlice[T](0, len(s))
	for member := range s {
		members.Append(member)
	}
	return members
}

func (s MapSet[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s MapSet[T]) Combine(other MapSet[T]) (success int) {
	for k := range other {
		if _, ok := s[k]; !ok {
			success += 1
			s[k] = struct{}{}
		}
	}
	return
}

func NewMapSet[T comparable]() MapSet[T] {
	return make(MapSet[T])
}
