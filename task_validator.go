package aysync

import "context"
import "golang.org/x/sync/errgroup"

type taskAny interface {
	GetError() error
}

func AreValid(ctx context.Context, tasks ...taskAny) error {
	if len(tasks) == 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.TODO()
	}
	g, gCtx := errgroup.WithContext(ctx)
	for _, task := range tasks {
		task := task
		g.Go(func() error {
			done := make(chan error)
			go func() {
				err := task.GetError()
				done <- err
			}()
			select {
			case <-gCtx.Done():
				return nil
			case err := <-done:
				return err
			}

		})
	}
	err := g.Wait()

	return err
}

func MapOnValid[T any](ctx context.Context, generator func() (T, error), tasks ...taskAny) Task[T] {
	if generator == nil {
		return newErrTask[T](ctx, ErrNilFuncEncountered)
	}

	err := AreValid(ctx, tasks...)
	if err != nil {
		return newErrTask[T](ctx, err)
	}

	return NewTask(ctx, generator)
}
