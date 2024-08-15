package aysync

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrTaskTimeout could be removed in future
	ErrTaskTimeout          = errors.New("task took too long to complete")
	ErrTaskContextCancelled = errors.New("context cancelled")
	ErrNilValueEncountered  = errors.New("null value encountered")
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
}
