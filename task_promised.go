package aysync

import "context"

type taskPromised[T any] struct {
	promised Task[T]
}

func Map[T, U any](ctx context.Context, tsk Task[T], mapper func(data T) (U, error)) Task[U] {
	promised := NewTask(ctx, func() (result U, err error) {
		var resolved T
		resolved, err = tsk.Await()
		if err != nil {
			return
		}
		result, err = mapper(resolved)

		return
	})

	return &taskPromised[U]{
		promised: promised,
	}
}

func (t *taskPromised[T]) Await() (res T, err error) {
	if t.promised == nil {
		err = ErrNilValueEncountered
		return
	}
	return t.promised.Await()
}

func (t *taskPromised[T]) Subscribe(cb func(data T, err error)) {
	go cb(t.Await())
}
