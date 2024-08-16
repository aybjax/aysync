package adt

import "sync"

type Map[K comparable, V any] struct {
	m   map[K]V
	mtx sync.RWMutex
}

func NewMap[K comparable, V any](m ...map[K]V) *Map[K, V] {
	result := &Map[K, V]{}

	if len(m) > 0 {
		result.m = m[0]
	} else {
		result.m = make(map[K]V)
	}

	return result
}

func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	val, ok = m.m[key]

	return val, ok
}

func (m *Map[K, V]) Set(key K, val V) (prev V, prevOK bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	prev, prevOK = m.m[key]
	m.m[key] = val

	return prev, prevOK
}

func (m *Map[K, V]) Del(key K) (prev V, prevOK bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	prev, prevOK = m.m[key]
	delete(m.m, key)

	return prev, prevOK
}

func (m *Map[K, V]) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	clear(m.m)
}

func (m *Map[K, V]) Mutate(f func(m map[K]V) map[K]V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.m = f(m.m)
}
