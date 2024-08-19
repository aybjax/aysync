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

type ForIterMap[T comparable, U any] struct {
	Key  T
	Data U
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

func ForMap[T comparable, U any](_ context.Context, m map[T]U) <-chan ForIterMap[T, U] {
	ch := make(chan ForIterMap[T, U])

	go func() {
		for k, v := range m {
			ch <- ForIterMap[T, U]{
				Key:  k,
				Data: v,
			}
		}

		close(ch)
	}()

	return ch
}

func ForMapKey[T comparable, U any](_ context.Context, m map[T]U) <-chan T {
	ch := make(chan T)

	go func() {
		for k := range m {
			ch <- k
		}

		close(ch)
	}()

	return ch
}

func ForMapVal[T comparable, U any](_ context.Context, m map[T]U) <-chan U {
	ch := make(chan U)

	go func() {
		for _, v := range m {
			ch <- v
		}

		close(ch)
	}()

	return ch
}
