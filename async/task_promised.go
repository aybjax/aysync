package async

import (
	"context"
)

type taskPromised[T any] struct {
	promised Task[T]
}

func FMap[T, U any](ctx context.Context, tsk Task[T], mapper func(data T) (U, error)) Task[U] {
	promisedChan := make(chan Task[U])

	promised := NewTask(ctx, func() (result U, err error) {
		var resolved T
		resolved, err = tsk.Await()
		if err != nil {
			return
		}
		result, err = mapper(resolved)

		return
	})

	go func() {
		promisedChan <- promised
	}()

	for i := 0; i < 100; i++ {
		select {
		case promised := <-promisedChan:
			return &taskPromised[U]{
				promised: promised,
			}
		case <-tsk.GetContext().Done():
			var errCause error = ErrTaskContextCancelled
			if err := context.Cause(tsk.GetContext()); err != nil {
				errCause = err
			}
			return NewErrTask[U](tsk.GetContext(), errCause)
		}
	}
	return NewErrTask[U](tsk.GetContext(), ErrTaskContextCancelled)
}

func FMap2[T1, T2, U any](ctx context.Context, tsk1 Task[T1], tsk2 Task[T2], mapper func(data1 T1, data2 T2) (U, error)) Task[U] {
	promisedChan := make(chan Task[U])

	promised := NewTask(ctx, func() (result U, err error) {
		var (
			resolved1 T1
			resolved2 T2
		)
		resolved1, err = tsk1.Await()
		if err != nil {
			return
		}
		resolved2, err = tsk2.Await()
		if err != nil {
			return
		}
		result, err = mapper(resolved1, resolved2)

		return
	})

	go func() {
		promisedChan <- promised
	}()

	select {
	case promised := <-promisedChan:
		return &taskPromised[U]{
			promised: promised,
		}
	case <-tsk1.GetContext().Done():
		return NewErrTask[U](tsk1.GetContext(), ErrTaskContextCancelled)
	case <-tsk2.GetContext().Done():
		return NewErrTask[U](tsk2.GetContext(), ErrTaskContextCancelled)
	}
}

func FMap3[T1, T2, T3, U any](ctx context.Context, tsk1 Task[T1], tsk2 Task[T2], tsk3 Task[T3], mapper func(data1 T1, data2 T2, data3 T3) (U, error)) Task[U] {
	promisedChan := make(chan Task[U])

	promised := NewTask(ctx, func() (result U, err error) {
		var (
			resolved1 T1
			resolved2 T2
			resolved3 T3
		)
		resolved1, err = tsk1.Await()
		if err != nil {
			return
		}
		resolved2, err = tsk2.Await()
		if err != nil {
			return
		}
		resolved3, err = tsk3.Await()
		if err != nil {
			return
		}
		result, err = mapper(resolved1, resolved2, resolved3)

		return
	})

	go func() {
		promisedChan <- promised
	}()

	select {
	case promised := <-promisedChan:
		return &taskPromised[U]{
			promised: promised,
		}
	case <-tsk1.GetContext().Done():
		return NewErrTask[U](tsk1.GetContext(), ErrTaskContextCancelled)
	case <-tsk2.GetContext().Done():
		return NewErrTask[U](tsk2.GetContext(), ErrTaskContextCancelled)
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

func (t *taskPromised[T]) GetContext() context.Context {
	return t.promised.GetContext()
}

func (t *taskPromised[T]) GetError() error {
	_, err := t.Await()
	return err
}
