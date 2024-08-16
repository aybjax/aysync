package adt

import "sync"

type Slice[V any] struct {
	s   []V
	mtx sync.RWMutex
}

func NewSlice[V any](s ...[]V) *Slice[V] {
	result := &Slice[V]{}

	if len(s) > 0 {
		result.s = s[0]
	}

	return result
}

func (s *Slice[V]) Get(index int) (val V) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.s[index]
}

func (s *Slice[V]) Len() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return len(s.s)
}

func (s *Slice[V]) Cap() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return cap(s.s)
}

func (s *Slice[V]) Slice(slices ...int) *Slice[V] {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	l := len(slices)

	if l == 0 {
		return s
	} else if l == 1 {
		return NewSlice(s.s[slices[0]:])
	} else if l == 2 {
		return NewSlice(s.s[slices[0]:slices[1]])
	}

	return nil
}

func (s *Slice[V]) Append(val ...V) *Slice[V] {
	// Lock
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	newSlice := append(s.s, val...)

	return NewSlice(newSlice)
}

func (s *Slice[V]) Set(index int, val V) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.s[index] = val
}

func (s *Slice[V]) Clear() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	clear(s.s)
}

func (s *Slice[V]) Mutate(f func(m []V) []V) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.s = f(s.s)
}
