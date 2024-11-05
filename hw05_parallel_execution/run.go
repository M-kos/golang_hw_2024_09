package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task)

	var errCounter int64
	var err error
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasksCh, &errCounter)
		}()
	}

	for i := 0; i < len(tasks); i++ {
		if atomic.LoadInt64(&errCounter) >= int64(m) {
			err = ErrErrorsLimitExceeded
			break
		}

		tasksCh <- tasks[i]
	}

	close(tasksCh)
	wg.Wait()

	return err
}

func worker(tasksCh chan Task, errCounter *int64) {
	for task := range tasksCh {
		if err := task(); err != nil {
			atomic.AddInt64(errCounter, 1)
		}
	}
}
