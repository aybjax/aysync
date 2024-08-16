package async

import "context"

type valueTask[T any] struct {
	ctx       context.Context
	otherwise T
	err       error
}

func (t *valueTask[T]) Await() (T, error) {
	return t.otherwise, t.err
}

func (t *valueTask[T]) Subscribe(cb func(data T, err error)) {
	go cb(t.otherwise, nil)
}

func (t *valueTask[T]) GetContext() context.Context {
	return t.ctx
}

func (t *valueTask[T]) GetError() error {
	return t.err
}
