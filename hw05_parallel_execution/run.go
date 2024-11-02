package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	done := make(chan struct{})

	var tasksCounter int64
	var errCounter int64
	var err error

	for i := 0; i < n; i++ {
		go worker(tasksCh, done, &tasksCounter, &errCounter)
	}

	go func() {
		for _, t := range tasks {
			tasksCh <- t
		}
	}()

	for {
		if errCounter >= int64(m) {
			err = ErrErrorsLimitExceeded
			break
		}

		if tasksCounter >= int64(len(tasks)) {
			break
		}
	}

	close(done)

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
