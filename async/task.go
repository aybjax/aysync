package async

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrTaskTimeout could be removed in future
	ErrParentTaskErrored    = errors.New("parent task returned error")
	ErrTaskTimeout          = errors.New("task took too long to complete")
	ErrTaskContextCancelled = context.Canceled
	ErrNilValueEncountered  = errors.New("null value encountered")
	ErrNilFuncEncountered   = errors.New("null value encountered")
)

const (
	// having timeout may be incorrect, maybe removed in the future
	defaultTimeout = time.Hour
)

type taskResult[T any] struct {
	data T
	err  error
}

type Task[T any] interface {
	Await() (T, error)
	Subscribe(cb func(data T, err error))
	GetContext() context.Context
	GetError() error
}
