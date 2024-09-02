package async

import (
	"context"
	"errors"
	"fmt"
	"runtime"
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

	var cancelFnx context.CancelCauseFunc
	ctx, cancelFnx = context.WithCancelCause(ctx)
	once := task[T]{}.createOnceFunc(ctx, cancelFnx, f)
	go once()
	runtime.Gosched()

	return &task[T]{
		retriever: once,
		ctx:       ctx,
	}
}

func (t task[T]) createOnceFunc(ctx context.Context, cancelFnx context.CancelCauseFunc, f func() (T, error)) func() (T, error) {
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
				doneRes <- result
				if result.err != nil {
					cancelFnx(result.err)
				}
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
	runtime.Gosched()
	return t.retriever()
}

func (t *task[T]) Subscribe(cb func(data T, err error)) {
	go cb(t.Await())
}

func (t *task[T]) GetContext() context.Context {
	return t.ctx
}

func (t *task[T]) GetError() error {
	_, err := t.retriever()
	return err
}
