package task

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

var (
	wgGlobal     sync.WaitGroup // Global wait group to synchronize goroutines
	succeedTasks sync.Map       // Map to store successful tasks
	failedTasks  sync.Map       // Map to store failed tasks
)

// taskCreator Starts task creation.
// Accepts a context for terminating work and a channel for transfer of new jobs
func taskCreator(ctx context.Context, newTaskChan chan<- TTask) {

	defer close(newTaskChan)

cycle:
	for {
		select {
		case <-ctx.Done():
			break cycle
		default:
			// Creating a task
			currentTime := time.Now().Format(time.RFC3339)
			if time.Now().Second()%2 > 0 { // Condition for erroneous tasks
				currentTime = "Some error occurred"
			}

			newTaskChan <- TTask{id: uuid.New(), creationTime: currentTime, state: NEW}
		}
	}
}

// taskWorker Processes tasks.
// Accepts a task for processing.
// Returns an error and changes the status of the task if it could not be processed.
func taskWorker(task TTask) (TTask, error) {

	taskTime, err := time.Parse(time.RFC3339, task.creationTime)
	if err != nil || taskTime.Before(time.Now().Add(-20*time.Second)) {
		task.result = []byte("something went wrong")
		task.state = FAILURE
		err = errors.Wrap(err, "task failed")
	} else {
		task.result = []byte("task has been succeeded")
		task.state = SUCCESS
	}

	task.finishTime = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)

	return task, err
}

// taskProcessor Starts task processing.
// Accepts a channel for receiving tasks and channels for transmitting successful and failed tasks.
func taskProcessor(newTaskChan <-chan TTask, successTaskChan, failureTaskChan chan<- TTask) {

	defer func() {
		close(successTaskChan)
		close(failureTaskChan)
	}()

	for task := range newTaskChan {
		processedTask, err := taskWorker(task)
		if err == nil {
			successTaskChan <- processedTask
		} else {
			failureTaskChan <- processedTask
		}
	}
}

// writer Writes the task to the log.
// Accepts context for terminating work.
func writer(ctx context.Context) {

	job := func() {
		log.Printf("Succeed tasks:\n")
		succeedTasks.Range(func(key, value interface{}) bool {
			log.Println(value.(TTask))
			return true
		})
		log.Printf("Failed tasks:\n")
		failedTasks.Range(func(key, value interface{}) bool {
			log.Println(value.(TTask))
			return true
		})
	}

	// Print results every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

cycle:
	for {
		select {
		case <-ticker.C:
			job()
		case <-ctx.Done():
			_ = <-ticker.C // Waiting for the ticker so as not to change the printing frequency
			job()
			break cycle
		}
	}
}

// resultCollector Collect tasks.
// Accepts channels for receiving successful and failed tasks.
func resultCollector(successTaskChan, failureTaskChan <-chan TTask) {

	wg := sync.WaitGroup{}

	// Collecting successful tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range successTaskChan {
			succeedTasks.Store(task.id, task)
		}
	}()

	// Collecting failed tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range failureTaskChan {
			failedTasks.Store(task.id, task)
		}
	}()

	wgWriter := sync.WaitGroup{}

	// Writing failed and succeed tasks
	wgWriter.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wgWriter.Done()
		writer(ctx)
	}()

	wg.Wait()
	cancel()

	wgWriter.Wait()
}

// Start Starts the task of generating and processing tasks.
// Accepts a context for terminating work.
func Start(ctx context.Context) {

	// Channel to receive tasks
	newTaskChan := make(chan TTask, 10)

	// Start task creator
	wgGlobal.Add(1)
	go func() {
		defer wgGlobal.Done()
		taskCreator(ctx, newTaskChan)
	}()

	// Channels for processed tasks
	successTaskChan := make(chan TTask)
	failureTaskChan := make(chan TTask)

	// Start processing tasks
	wgGlobal.Add(1)
	go func() {
		defer wgGlobal.Done()
		taskProcessor(newTaskChan, successTaskChan, failureTaskChan)
	}()

	// Collect results
	wgGlobal.Add(1)
	go func() {
		defer wgGlobal.Done()
		resultCollector(successTaskChan, failureTaskChan)
	}()

	// Wait for all goroutines to finish
	wgGlobal.Wait()
}
