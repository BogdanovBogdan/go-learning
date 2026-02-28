package todo

import (
	"sync/atomic"
	"time"
)

var taskID atomic.Int64

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func getNextTaskID() int {
	return int(taskID.Add(1))
}

func NewTask(title string, description string) Task {

	return Task{
		ID:          getNextTaskID(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (task *Task) CompleteTask() {
	timeNow := time.Now()

	task.Completed = true
	task.CompletedAt = &timeNow
}

func (task *Task) UncompleteTask() {
	task.Completed = false
	task.CompletedAt = nil
}
