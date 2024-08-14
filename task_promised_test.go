package aysync

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMap(t *testing.T) {
	result := 1
	err := errors.New("i am error")
	cancelledCtx, cancel := context.WithCancel(context.TODO())
	cancel()

	testCases := []struct {
		name          string
		ctx           context.Context
		taskGenerator func() (int, error)
		mapper        func(int) (any, error)
		expected      any
		err           error
	}{
		{
			name: "test simple function result",
			taskGenerator: func() (int, error) {
				return result, nil
			},
			mapper: func(i int) (any, error) {
				return i + 1, nil
			},
			expected: result + 1,
			err:      nil,
		},
		{
			name: "test simple function result (type change)",
			taskGenerator: func() (int, error) {
				return result, nil
			},
			mapper: func(i int) (any, error) {
				return fmt.Sprintf("%d", i), nil
			},
			expected: fmt.Sprintf("%d", result),
			err:      nil,
		},
		{
			name: "test simple function error",
			taskGenerator: func() (int, error) {
				return 100, err
			},
			mapper: func(i int) (any, error) {
				return fmt.Sprintf("%d", i), nil
			},
			expected: nil,
			err:      err,
		},
		{
			name: "test simple function error",
			taskGenerator: func() (int, error) {
				return 100, nil
			},
			mapper: func(i int) (any, error) {
				return nil, err
			},
			expected: nil,
			err:      err,
		},
		{
			ctx:  cancelledCtx,
			name: "test result with timeout error",
			taskGenerator: func() (int, error) {
				return 100, nil
			},
			mapper: func(i int) (any, error) {
				return fmt.Sprintf("%d", i), err
			},
			expected: nil,
			err:      ErrTaskContextCancelled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(tc.ctx, tc.taskGenerator)
			promised := Map(tc.ctx, task, tc.mapper)
			result, err := promised.Await()
			require.Equal(t, result, tc.expected)
			require.Equal(t, err, tc.err)
		})
	}
}
