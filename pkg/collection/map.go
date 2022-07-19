package collection

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Add(k K, v V) (success bool) {
	if _, ok := m[k]; !ok {
		m[k] = v
		success = true
	}

	return
}

func (m Map[K, V]) Remove(k K) V {
	v := m[k]
	delete(m, k)
	return v
}

func (m Map[K, V]) Get(k K, def ...V) V {
	v, ok := m[k]
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

func (m Map[K, V]) Update(k K, v V) (old V) {
	old = m[k]
	m[k] = v
	return
}

func (m Map[K, V]) Combine(other Map[K, V]) (success int) {
	for k, v := range other {
		if _, ok := m[k]; !ok {
			m[k] = v
			success += 1
		}
	}
	return
}

func NewMap[K comparable, V any]() Map[K, V] {
	return make(Map[K, V])
}
