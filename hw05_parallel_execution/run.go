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
	done := make(chan struct{})

	var tasksCounter int64
	var errCounter int64
	var err error
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasksCh, done, &tasksCounter, &errCounter)
		}()
	}

	for i := 0; i < len(tasks); i++ {
		tasksCh <- tasks[i]

		if atomic.LoadInt64(&errCounter) >= int64(m) {
			err = ErrErrorsLimitExceeded
			break
		}

		if atomic.LoadInt64(&tasksCounter) >= int64(len(tasks)) {
			break
		}
	}

	close(done)
	wg.Wait()
	close(tasksCh)

	return err
}

func worker(tasksCh chan Task, done chan struct{}, tCounter, errCounter *int64) {
	for {
		select {
		case task := <-tasksCh:
			if err := task(); err != nil {
				atomic.AddInt64(errCounter, 1)
			}

			atomic.AddInt64(tCounter, 1)
		case <-done:
			return
		}
	}
}
