package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in workersCount goroutines
// and stops its work when receiving maxErrosCount errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	wg := sync.WaitGroup{}

	taskChan := make(chan Task)
	stopChan := make(chan struct{})
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		startWorkers(workersCount, taskChan, stopChan, errChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskChan)

		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-stopChan:
				return
			}
		}
	}()

	errCount := 0
	var err error

	ignoreErrors := maxErrorsCount <= 0

	for {
		_, ok := <-errChan
		if !ok {
			break
		}
		if ignoreErrors {
			continue
		}
		errCount++

		if errCount >= maxErrorsCount {
			err = ErrErrorsLimitExceeded
			close(stopChan)
			break
		}
	}

	wg.Wait()
	return err
}

func startWorker(chTask <-chan Task, chStop <-chan struct{}, chErr chan<- error) {
	for {
		task, ok := <-chTask
		if !ok {
			return
		}
		err := task()
		select {
		case <-chStop:
			return

		default:
			if err != nil {
				chErr <- err
			}
		}
	}
}

func startWorkers(numWorkers int, chTask <-chan Task, chStop <-chan struct{}, chErr chan<- error) {
	defer close(chErr)

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			startWorker(chTask, chStop, chErr)
		}()
	}
	wg.Wait()
}
