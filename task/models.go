package task

import (
	"fmt"
	"github.com/google/uuid"
)

// TTask represents a task with its metadata
type TTask struct {
	id           uuid.UUID
	creationTime string // creation time
	finishTime   string // finish time
	result       []byte
	state        EnumState
}

type EnumState string

const (
	SUCCESS EnumState = "success"
	FAILURE EnumState = "failure"
	NEW     EnumState = "new"
)

func (task TTask) String() string {
	if task.state == SUCCESS {
		return fmt.Sprintf("Task ID: %s --- Result: %s", task.id, task.result)
	} else if task.state == FAILURE {
		return fmt.Sprintf("Task ID: %s --- Creation time: %s --- Result: %s", task.id, task.creationTime, task.result)
	} else {
		return fmt.Sprintf("Task ID: %s --- Creation time: %s", task.id, task.creationTime)
	}
}
