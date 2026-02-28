package todo

import (
	"sync"
)

type List struct {
	tasks map[int]Task
	mutex sync.RWMutex
}

func NewTodoList() *List {
	return &List{
		tasks: make(map[int]Task),
	}
}

func (list *List) AddTask(task *Task) {
	list.mutex.Lock()
	list.tasks[task.ID] = *task
	list.mutex.Unlock()
}

func (list *List) GetTask(ID int) (Task, error) {
	list.mutex.RLock()
	defer list.mutex.RUnlock()

	task, ok := list.tasks[ID]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (list *List) GetAllTasks() []Task {
	list.mutex.RLock()
	defer list.mutex.RUnlock()

	response := make([]Task, 0, len(list.tasks))

	for _, task := range list.tasks {
		response = append(response, task)
	}

	return response
}

func (list *List) GetFilteredTasks(isCompleted bool) []Task {
	list.mutex.RLock()
	defer list.mutex.RUnlock()

	filteredTasks := make([]Task, 0)

	for _, task := range list.tasks {
		if task.Completed == isCompleted {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks
}

func (list *List) ToggleTask(ID int, isCompleted bool) (Task, error) {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	task, ok := list.tasks[ID]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	if isCompleted {
		task.CompleteTask()
	} else {
		task.UncompleteTask()
	}

	list.tasks[ID] = task

	return task, nil
}

func (list *List) DeleteTask(ID int) error {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	if _, ok := list.tasks[ID]; !ok {
		return ErrTaskNotFound
	}

	delete(list.tasks, ID)

	return nil
}
