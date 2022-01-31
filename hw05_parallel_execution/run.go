package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded      = errors.New("errors limit exceeded")
	ErrWrongExpectedErrorNumber = errors.New("expected errors number must be positive")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrWrongExpectedErrorNumber
	}

	var errorsCount int32
	ch := make(chan Task)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range ch {
				if err := task(); err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorsCount) >= int32(m) {
			break
		}
		ch <- task
	}

	close(ch)

	wg.Wait()

	if errorsCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
