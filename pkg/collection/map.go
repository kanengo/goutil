package collection

type Map[K comparable, V any] struct {
	m map[K]V
	*defaultIterator[Pairs[K, V]]
}

func (m Map[K, V]) Add(k K, v V) (success bool) {
	if _, ok := m.m[k]; !ok {
		m.m[k] = v
		success = true
	}

	return
}

func (m Map[K, V]) Remove(k K) V {
	v := m.m[k]
	delete(m.m, k)
	return v
}

func (m Map[K, V]) Get(k K, def ...V) (V, bool) {
	v, ok := m.m[k]
	if !ok && len(def) > 0 {
		return def[0], false
	}
	return v, ok
}

func (m Map[K, V]) Update(k K, v V) (old V) {
	old = m.m[k]
	m.m[k] = v
	return
}

func (m Map[K, V]) Keys() *Slice[K] {
	keys := NewSlice[K](0, len(m.m))
	for key := range m.m {
		keys.Append(key)
	}

	return keys
}

func (m Map[K, V]) Values() []V {
	values := make([]V, 0, len(m.m))
	for _, value := range m.m {
		values = append(values, value)
	}

	return values
}

func (m Map[K, V]) Clear() {
	for k := range m.m {
		delete(m.m, k)
	}
}

func (m Map[K, V]) Combine(other Map[K, V]) (success int) {
	for k, v := range other.m {
		if _, ok := m.m[k]; !ok {
			m.m[k] = v
			success += 1
		}
	}
	return
}

func (m Map[K, V]) Range(fn func(p Pairs[K, V]) bool) {
	for key, value := range m.m {
		if !fn(Pairs[K, V]{key, value}) {
			break
		}
	}
}

func NewMap[K comparable, V comparable]() *Map[K, V] {
	m := &Map[K, V]{
		m: make(map[K]V, 10),
	}
	m.defaultIterator = &defaultIterator[Pairs[K, V]]{
		Ranger: m,
	}
	return m
}
