package adt

import "sync"

type Container[T any] struct {
	data T
	mtx  sync.RWMutex
}

func (c *Container[T]) Get(f func(data T) any) any {
	c.mtx.RLock()
	c.mtx.RUnlock()

	return f(c.data)
}

func (c *Container[T]) Set(f func(data T) T) {
	c.mtx.Lock()
	c.mtx.Unlock()

	c.data = f(c.data)
}

func ContainerGet[C any, V any](c Container[C], f func(data C) V) V {
	resultAny := c.Get(func(data C) any {
		return f(data)
	})

	return resultAny.(V)
}
