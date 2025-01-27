package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidRoutineCount = errors.New("invalid routine count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 means ignore errors.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrInvalidRoutineCount
	}
	if len(tasks) < n {
		// start less or equal goroutines than tasks
		n = len(tasks)
	}
	errorCount := &atomic.Int64{}
	errorCountMax := int64(m)
	processingQueue := make(chan Task)
	wg := new(sync.WaitGroup)
	// start workers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range processingQueue {
				if m <= 0 || errorCount.Load() < errorCountMax {
					if err := task(); err != nil {
						errorCount.Add(1)
					}
				}
			}
		}()
	}
	// fill tasks queue
	for _, task := range tasks {
		processingQueue <- task
	}
	close(processingQueue)
	wg.Wait()

	if m > 0 && errorCount.Load() >= errorCountMax {
		return ErrErrorsLimitExceeded
	}
	return nil
}
