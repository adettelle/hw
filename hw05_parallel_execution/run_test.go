package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak" //nolint:depguard
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestMaxErrorsCount(t *testing.T) {
	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	tests := []struct {
		name           string
		workersCount   int
		maxErrorsCount int
		want           int
	}{
		{
			name:           "zero maxErrorsCount",
			workersCount:   5,
			maxErrorsCount: 0,
		},
		{
			name:           "negative maxErrorsCount",
			workersCount:   5,
			maxErrorsCount: -1,
		},
	}

	for _, test := range tests {
		var runTasksCount int32

		for i := range tasksCount {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				if i%2 == 0 {
					return fmt.Errorf("error in task")
				}
				return nil
			})
		}

		t.Run(test.name, func(t *testing.T) {
			err := Run(tasks, test.workersCount, test.maxErrorsCount)
			require.NoError(t, err)
			require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		})
	}
}
