package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func(ch chan int32) error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	totalCount := len(tasks)
	errorsLimit := int32(m)
	wg := sync.WaitGroup{}

	ch := make(chan int32)
	defer close(ch)

	var total int
	for total < totalCount {
		wg.Add(n)
		for i := 0; i < n; i++ {
			total++

			task := tasks[i]
			go task(ch)
			wg.Done()

			count := <-ch
			if errorsLimit != 0 && count >= errorsLimit {
				return ErrErrorsLimitExceeded
			}
		}
	}

	wg.Wait()

	return nil
}
