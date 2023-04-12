package matcher

import (
	"sort"
	"strings"
)

type Matcher[T any] interface {
	Add(selector string, ms ...T)
	Match(operation string) []T
	Use(ms ...T)
}

type matcher[T any] struct {
	prefix   []string
	defaults []T
	matches  map[string][]T
}

func New[T any]() Matcher[T] {
	return &matcher[T]{}
}

func (m *matcher[T]) Use(ms ...T) {
	m.defaults = ms
}

func (m *matcher[T]) Add(selector string, ms ...T) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.prefix = append(m.prefix, selector)
		sort.Slice(m.prefix, func(i, j int) bool {
			return m.prefix[i] > m.prefix[j]
		})
	}
	m.matches[selector] = append(m.matches[selector], ms...)
}

func (m *matcher[T]) Match(operation string) []T {
	ms := make([]T, 0, len(m.defaults))
	if len(m.defaults) > 0 {
		ms = append(ms, m.defaults...)
	}
	if next, ok := m.matches[operation]; ok {
		return append(ms, next...)
	}

	for _, prefix := range m.prefix {
		if strings.HasPrefix(operation, prefix) {
			return append(ms, m.matches[prefix]...)
		}
	}

	return ms
}
