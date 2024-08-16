package async

import "context"

func Tern[T any](ctx context.Context, cond bool, taskGen func() (T, error), otherwise T) Task[T] {
	if !cond {
		return &valueTask[T]{
			ctx:       ctx,
			otherwise: otherwise,
		}
	}

	return NewTask[T](ctx, taskGen)
}

func TernFunc[T any](ctx context.Context, cond bool, taskGen func() (T, error), otherwiseGen func() (T, error)) Task[T] {
	if !cond {
		otherwise, err := otherwiseGen()
		return &valueTask[T]{
			ctx:       ctx,
			otherwise: otherwise,
			err:       err,
		}
	}

	return NewTask[T](ctx, taskGen)
}

func TernTask[T any](ctx context.Context, cond bool, taskGen func() (T, error), otherwiseTaskGen func() (T, error)) Task[T] {
	if !cond {
		return NewTask[T](ctx, otherwiseTaskGen)
	}

	return NewTask[T](ctx, taskGen)
}
