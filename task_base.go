package aysync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type task[T any] struct {
	retriever func() (T, error)
	ctx       context.Context
}

func NewTask[T any](ctx context.Context, f func() (T, error)) Task[T] {
	if ctx == nil {
		ctx = context.TODO()
	}
	var cancelFnx context.CancelFunc
	ctx, cancelFnx = context.WithCancel(ctx)

	once := task[T]{}.createOnceFunc(ctx, cancelFnx, f)
	go once()

	return &task[T]{
		retriever: once,
		ctx:       ctx,
	}
}

func (t task[T]) createOnceFunc(ctx context.Context, cancelFnx context.CancelFunc, f func() (T, error)) func() (T, error) {
	once := sync.OnceValues(func() (T, error) {
		funcRes := make(chan *taskResult[T])
		doneRes := make(chan *taskResult[T])

		go func() {
			var (
				data T
				err  error
			)
			defer func() {
				if excp := recover(); excp != nil {
					switch v := excp.(type) {
					case string:
						err = fmt.Errorf("task panic'd: %s", v)
					case error:
						err = fmt.Errorf("task panic'd: %w", v)
					default:
						err = errors.New(fmt.Sprint(v))
					}
				}

				funcRes <- &taskResult[T]{
					data: data,
					err:  err,
				}
			}()

			data, err = f()
		}()

		go func() {
			select {
			case result := <-funcRes:
				if result.err != nil {
					cancelFnx()
				}
				doneRes <- result
			case <-time.After(defaultTimeout):
				doneRes <- &taskResult[T]{
					err: ErrTaskTimeout,
				}
			case <-ctx.Done():
				doneRes <- &taskResult[T]{
					err: ErrTaskContextCancelled,
				}
			}
		}()

		result := <-doneRes

		return result.data, result.err
	})

	return once
}

func (t *task[T]) Await() (T, error) {
	return t.retriever()
}

func (t *task[T]) Subscribe(cb func(data T, err error)) {
	go cb(t.Await())
}

func (t *task[T]) GetContext() context.Context {
	return t.ctx
}
