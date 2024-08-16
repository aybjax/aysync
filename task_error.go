package aysync

import "context"

type taskErr[T any] struct {
	ctx  context.Context
	dflt T
	err  error
}

func newErrTask[T any](ctx context.Context, err error) Task[T] {
	return &taskErr[T]{
		ctx: ctx,
		err: err,
	}
}

func (t *taskErr[T]) Await() (T, error) {
	return t.dflt, t.err
}

func (t *taskErr[T]) Subscribe(cb func(data T, err error)) {
	go cb(t.dflt, t.err)
}

func (t *taskErr[T]) GetContext() context.Context {
	return t.ctx
}

func (t *taskErr[T]) GetError() error {
	return t.err
}
