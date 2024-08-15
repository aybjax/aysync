package aysync

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTask_Await(t *testing.T) {
	result := 1
	dflt := 0
	err := errors.New("i am error")
	cancelledCtx, cancel := context.WithCancel(context.TODO())
	cancel()

	testcases := []struct {
		name     string
		ctx      context.Context
		f        func() (int, error)
		expected int
		err      error
	}{
		{
			name: "test simple function result",
			f: func() (int, error) {
				return result, nil
			},
			expected: result,
			err:      nil,
		},
		{
			name: "test simple function error",
			f: func() (int, error) {
				return dflt, err
			},
			expected: dflt,
			err:      err,
		},
		{
			name: "test result with timeout",
			f: func() (int, error) {
				time.Sleep(5 * time.Millisecond)
				return result, nil
			},
			expected: result,
			err:      nil,
		},
		{
			name: "aybsync: test result with timeout error",
			f: func() (int, error) {
				time.Sleep(1100 * time.Millisecond)
				return result, nil
			},
			ctx:      cancelledCtx,
			expected: dflt,
			err:      ErrTaskContextCancelled,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(tc.ctx, tc.f)
			result, err := task.Await()
			require.Equal(t, result, tc.expected)
			require.Equal(t, err, tc.err)
		})
	}
}

func TestTask_Subscribe(t *testing.T) {
	result := 1
	dflt := 0
	err := errors.New("i am error")
	cancelledCtx, cancel := context.WithCancel(context.TODO())
	cancel()

	testcases := []struct {
		ctx      context.Context
		name     string
		f        func() (int, error)
		expected int
		err      error
	}{
		{
			name: "test simple function result",
			f: func() (int, error) {
				return result, nil
			},
			expected: result,
			err:      nil,
		},
		{
			name: "test simple function error",
			f: func() (int, error) {
				return dflt, err
			},
			expected: dflt,
			err:      err,
		},
		{
			name: "test result with timeout",
			f: func() (int, error) {
				time.Sleep(5 * time.Millisecond)
				return result, nil
			},
			expected: result,
			err:      nil,
		},
		{
			ctx:  cancelledCtx,
			name: "aybsync: test result with timeout error",
			f: func() (int, error) {
				return result, nil
			},
			expected: dflt,
			err:      ErrTaskContextCancelled,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(tc.ctx, tc.f)
			task.Subscribe(func(result int, err error) {
				require.Equal(t, result, tc.expected)
				require.Equal(t, err, tc.err)
			})
		})
	}
}

func TestTask_MultipleTasks(t *testing.T) {
	result := 1

	longFunction := func() (int, error) {
		time.Sleep(time.Millisecond * 500)
		return result, nil
	}

	tasks := append([]Task[int]{},
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
		NewTask(nil, longFunction),
	)

	start := time.Now()
	for _, task := range tasks {
		res, err := task.Await()

		require.Equal(t, err, nil)
		require.Equal(t, res, result)
	}
	dur := time.Since(start)
	if dur.Milliseconds() > 550 {
		t.Errorf("took longer than 500ms")
	}
}

func TestTask_ConcurrentCalls(t *testing.T) {
	result := 1

	longFunction := func() (int, error) {
		time.Sleep(time.Second * 1)
		return result, nil
	}

	results := [10]int{}
	task := NewTask(nil, longFunction)

	var wg sync.WaitGroup
	for i := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res, err := task.Await()
			if err != nil {
				t.Error(err)
			}
			results[i] = res
		}()
	}
	wg.Wait()

	for _, el := range results {
		if el != result {
			t.Error("value is not set")
		}
	}
}

func TestTask_GetContext(t *testing.T) {
	longFunction := func() (int, error) {
		return 0, errors.New("any error")
	}

	task := NewTask(nil, longFunction)
	ctx := task.GetContext()
	task2 := NewTask(ctx, longFunction)

	select {
	case <-ctx.Done():
		break
	case <-time.Tick(time.Second * 2):
		t.Error("context is not cancelled")
	}
	_, err := task2.Await()
	require.Error(t, err, ErrTaskContextCancelled.Error())
}
