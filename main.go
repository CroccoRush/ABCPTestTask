package main

import (
	"ABCPTestTask/task"
	"context"
	"sync"
	"time"
)

func main() {

	// Context to break goroutines
	ctx, cancel := context.WithCancel(context.Background())

	// WaitGroup to wait until the task stops
	wg := sync.WaitGroup{}

	// Starting the task
	wg.Add(1)
	go func() {
		defer wg.Done()
		task.Start(ctx)
	}()

	// Run for 10 seconds
	time.Sleep(10 * time.Second)
	cancel()

	wg.Wait()
}
