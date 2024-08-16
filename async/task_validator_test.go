package async

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTaskValidator_AreValid(t *testing.T) {
	err := errors.New("me, I am an error")
	testCases := []struct {
		name         string
		tasks        []func() taskAny
		runningRange [2]time.Duration
		err          error
	}{
		{
			name:         "test all ok",
			tasks:        []func() taskAny{},
			runningRange: [2]time.Duration{time.Millisecond * 0, time.Millisecond * 10},
			err:          nil,
		},
		{
			name: "test all ok",
			tasks: []func() taskAny{
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
			},
			runningRange: [2]time.Duration{time.Millisecond * 950, time.Millisecond * 1050},
			err:          nil,
		},
		{
			name: "test error",
			tasks: []func() taskAny{
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Millisecond * 50)
						return 0, err
					})
				},
			},
			runningRange: [2]time.Duration{time.Millisecond * 0, time.Millisecond * 100},
			err:          err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			tasks := make([]taskAny, len(tc.tasks))
			for i, f := range tc.tasks {
				tasks[i] = f()
			}
			err := AreValid(nil, tasks...)
			dur := time.Since(start)

			if dur.Milliseconds() > tc.runningRange[1].Milliseconds() ||
				dur.Milliseconds() < tc.runningRange[0].Milliseconds() {
				t.Errorf("running outside of range[%d:%d]: %d", tc.runningRange[0].Milliseconds(), tc.runningRange[1].Milliseconds(), dur.Milliseconds())
			}
			require.Equal(t, tc.err, err)
		})
	}
}

func TestMapOnValid(t *testing.T) {
	result := 1
	err := errors.New("me, I am an error")
	testCases := []struct {
		name         string
		tasks        []func() taskAny
		generator    func() (int, error)
		result       int
		runningRange [2]time.Duration
		err          error
	}{
		{
			name: "test all ok",
			generator: func() (int, error) {
				return result, nil
			},
			result: result,
			tasks: []func() taskAny{
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
			},
			runningRange: [2]time.Duration{time.Millisecond * 950, time.Millisecond * 1050},
			err:          nil,
		},
		{
			name: "test generator error",
			generator: func() (int, error) {
				return result, err
			},
			result: result,
			tasks: []func() taskAny{
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
			},
			runningRange: [2]time.Duration{time.Millisecond * 950, time.Millisecond * 1050},
			err:          err,
		},
		{
			name: "test error",
			tasks: []func() taskAny{
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Second)
						return 0, nil
					})
				},
				func() taskAny {
					return NewTask[int](nil, func() (int, error) {
						time.Sleep(time.Millisecond * 50)
						return 0, err
					})
				},
			},
			generator: func() (int, error) {
				return result, nil
			},
			result:       0,
			runningRange: [2]time.Duration{time.Millisecond * 0, time.Millisecond * 100},
			err:          err,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			tasks := make([]taskAny, len(tc.tasks))
			for i, f := range tc.tasks {
				tasks[i] = f()
			}
			task := MapOnValid(nil, tc.generator, tasks...)
			result, err := task.Await()
			dur := time.Since(start)

			if dur.Milliseconds() > tc.runningRange[1].Milliseconds() ||
				dur.Milliseconds() < tc.runningRange[0].Milliseconds() {
				t.Errorf("running outside of range: %d", dur.Milliseconds())
			}
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}
