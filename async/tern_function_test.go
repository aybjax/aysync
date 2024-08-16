package async

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTern(t *testing.T) {
	ctx := context.TODO()
	taskRes := 1
	otherwiseRes := 2
	err := errors.New("i am an error")
	testCases := []struct {
		name      string
		ctx       context.Context
		condition bool
		taskGen   func() (int, error)
		otherwise int
		result    int
		err       error
	}{
		{
			name:      "true case: task runs",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: otherwiseRes,
			result:    taskRes,
			err:       nil,
		},
		{
			name:      "false case: task runs",
			ctx:       ctx,
			condition: false,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: otherwiseRes,
			result:    otherwiseRes,
			err:       nil,
		},
		{
			name:      "true case: task errors",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, err
			},
			otherwise: otherwiseRes,
			result:    taskRes,
			err:       err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tsk := Tern(tc.ctx, tc.condition, tc.taskGen, tc.otherwise)

			result, err := tsk.Await()
			require.Equal(t, result, tc.result)
			require.Equal(t, err, tc.err)
		})
	}
}

func TestTernFunc(t *testing.T) {
	ctx := context.TODO()
	taskRes := 1
	otherwiseRes := 2
	err := errors.New("i am an error")
	testCases := []struct {
		name      string
		ctx       context.Context
		condition bool
		taskGen   func() (int, error)
		otherwise func() (int, error)
		result    int
		err       error
	}{
		{
			name:      "true case: task runs",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: taskRes,
			err:    nil,
		},
		{
			name:      "false case: task runs",
			ctx:       ctx,
			condition: false,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: otherwiseRes,
			err:    nil,
		},
		{
			name:      "true case: task errors",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, err
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: taskRes,
			err:    err,
		},
		{
			name:      "false case: otherwise errors",
			ctx:       ctx,
			condition: false,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, err
			},
			result: otherwiseRes,
			err:    err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tsk := TernFunc(tc.ctx, tc.condition, tc.taskGen, tc.otherwise)

			result, err := tsk.Await()

			require.Equal(t, result, tc.result)
			require.Equal(t, err, tc.err)
		})
	}
}

func TestTernTask(t *testing.T) {
	ctx := context.TODO()
	taskRes := 1
	otherwiseRes := 2
	err := errors.New("i am an error")
	testCases := []struct {
		name      string
		ctx       context.Context
		condition bool
		taskGen   func() (int, error)
		otherwise func() (int, error)
		result    int
		err       error
	}{
		{
			name:      "true case: task runs",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: taskRes,
			err:    nil,
		},
		{
			name:      "false case: task runs",
			ctx:       ctx,
			condition: false,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: otherwiseRes,
			err:    nil,
		},
		{
			name:      "true case: task errors",
			ctx:       ctx,
			condition: true,
			taskGen: func() (int, error) {
				return taskRes, err
			},
			otherwise: func() (int, error) {
				return otherwiseRes, nil
			},
			result: taskRes,
			err:    err,
		},
		{
			name:      "false case: otherwise errors",
			ctx:       ctx,
			condition: false,
			taskGen: func() (int, error) {
				return taskRes, nil
			},
			otherwise: func() (int, error) {
				return otherwiseRes, err
			},
			result: otherwiseRes,
			err:    err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tsk := TernTask(tc.ctx, tc.condition, tc.taskGen, tc.otherwise)

			result, err := tsk.Await()

			require.Equal(t, result, tc.result)
			require.Equal(t, err, tc.err)
		})
	}
}
