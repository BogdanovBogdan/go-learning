package server

import (
	"encoding/json"
	"errors"
	"strings"
	"todo-app/todo"
)

type TaskStore interface {
	AddTask(task *todo.Task)
	GetTask(ID int) (todo.Task, error)
	GetAllTasks() []todo.Task
	GetFilteredTasks(isCompleted bool) []todo.Task
	ToggleTask(ID int, isCompleted bool) (todo.Task, error)
	DeleteTask(ID int) error
}

type TaskCreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (task *TaskCreateDTO) Validate() error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("Title is required")
	}

	return nil
}

type TaskCompletedDTO struct {
	Complete bool `json:"complete"`
}

type ResponseDTO struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func (res ResponseDTO) ToJSON() ([]byte, error) {
	msg, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (res ResponseDTO) ToString() (string, error) {
	msg, err := res.ToJSON()
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
