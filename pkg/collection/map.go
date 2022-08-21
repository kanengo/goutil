package collection

import (
	"sync"
)

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

func (m Map[K, V]) Get(k K, def ...V) (V, bool) {
	v, ok := m[k]
	if !ok && len(def) > 0 {
		return def[0], false
	}
	return v, ok
}

func (m Map[K, V]) Update(k K, v V) (old V) {
	old = m[k]
	m[k] = v
	return
}

func (m Map[K, V]) Keys() *Slice[K] {
	keys := NewSlice[K](0, len(m))
	for key := range m {
		keys.Append(key)
	}

	return keys
}

func (m Map[K, V]) Values() []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}

	return values
}

func (m Map[K, V]) Clear() {
	for k := range m {
		delete(m, k)
	}
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

func NewMap[K comparable, V comparable]() Map[K, V] {
	return make(Map[K, V])
}

// SyncMap 包装sync.Map
type SyncMap[K comparable, V any] struct {
	sm *sync.Map
}

func (m SyncMap[K, V]) Add(k K, v V) {
	m.sm.Store(k, v)
}

func (m SyncMap[K, V]) Remove(k K) V {
	v, loaded := m.sm.LoadAndDelete(k)
	if !loaded {
		var zero V
		return zero
	}
	return v.(V)
}

func (m SyncMap[K, V]) Get(k K) (V, bool) {
	val, ok := m.sm.Load(k)
	if !ok {
		var zero V
		return zero, ok
	}

	return val.(V), true
}

func (m SyncMap[K, V]) Update(k K, v V) {
	m.sm.Store(k, v)
}

func (m SyncMap[K, V]) Keys() *Slice[K] {
	keys := NewSlice[K]()
	m.sm.Range(func(key, value any) bool {
		keys.Append(key.(K))
		return true
	})

	return keys
}

func (m SyncMap[K, V]) Values() []V {
	values := make([]V, 0)
	m.sm.Range(func(key, value any) bool {
		values = append(values, value.(V))
		return true
	})

	return values
}

func (m SyncMap[K, V]) Clear() {
	m.sm.Range(func(key, value any) bool {
		m.sm.Delete(key)
		return true
	})
}

func (m SyncMap[K, V]) Range(fn func(k K, v V) bool) {
	m.sm.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}

func NewSyncMap[K comparable, V comparable]() SyncMap[K, V] {
	return SyncMap[K, V]{
		sm: &sync.Map{},
	}
}
