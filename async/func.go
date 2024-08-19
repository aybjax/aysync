package async

import (
	"context"
	"runtime"
)

func Await(_ context.Context) {
	runtime.Gosched()
}

func ForItem[T any](_ context.Context, s []T) <-chan T {
	ch := make(chan T)

	go func() {
		for _, el := range s {
			ch <- el
		}

		close(ch)
	}()

	return ch
}

func ForIndex[T any](_ context.Context, s []T) <-chan int {
	ch := make(chan int)

	go func() {
		for ind := range s {
			ch <- ind
		}

		close(ch)
	}()

	return ch
}

type ForIter[T any] struct {
	Index int
	Data  T
}

func For[T any](_ context.Context, s []T) <-chan ForIter[T] {
	ch := make(chan ForIter[T])

	go func() {
		for ind, el := range s {
			ch <- ForIter[T]{
				Index: ind,
				Data:  el,
			}
		}

		close(ch)
	}()

	return ch
}
