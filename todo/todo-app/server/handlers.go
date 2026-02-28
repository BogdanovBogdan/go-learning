package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todo-app/todo"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	taskStore TaskStore
}

func NewHTTPHandlers(taskStore TaskStore) *HTTPHandlers {
	return &HTTPHandlers{
		taskStore: taskStore,
	}
}

func parseIDFromRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	IDStr := mux.Vars(r)["id"]

	IDInt, err := strconv.Atoi(IDStr)
	if err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		http.Error(w, response, http.StatusBadRequest)
		return 0, err
	}

	return IDInt, nil
}

/*
pattern: /tasks
method:  POST
info:    JSON in HTTP request body

successful:
  - status code:   201 Created
  - response body: JSON represent created task

failed:
  - status code:   400, 500
  - response body: JSON with error and status
*/
func (httpHandler *HTTPHandlers) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var dto TaskCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		http.Error(w, response, http.StatusBadRequest)
		return
	}

	if err := dto.Validate(); err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		http.Error(w, response, http.StatusBadRequest)
		return
	}

	newTask := todo.NewTask(dto.Title, dto.Description)
	httpHandler.taskStore.AddTask(&newTask)

	response, _ := ResponseDTO{
		Success: true,
		Result:  newTask,
	}.ToJSON()

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(response); err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
pattern: /tasks
method:  GET
info:    -

successful:
  - status code:   200
  - response body: list of tasks in JSON

failed:
  - status code:   -
  - response body: -
*/
func (httpHandler *HTTPHandlers) handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := httpHandler.taskStore.GetAllTasks()

	response, _ := ResponseDTO{
		Success: true,
		Result:  tasks,
	}.ToJSON()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
pattern: /tasks/{id}
method:  GET
info:    ID in pattern

successful:
  - status code:   200 OK
  - response body: task in JSON

failed:
  - status code:   400, 404, 500
  - response body: JSON with error and status
*/
func (httpHandler *HTTPHandlers) handleGetTask(w http.ResponseWriter, r *http.Request) {
	ID, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}

	task, err := httpHandler.taskStore.GetTask(ID)
	if err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, response, http.StatusNotFound)
		} else {
			http.Error(w, response, http.StatusInternalServerError)
		}
		return
	}

	response, _ := ResponseDTO{
		Success: true,
		Result:  task,
	}.ToJSON()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
pattern: /tasks?complete=true|false
method:  GET
info:    query param 'complete'

successful:
  - status code:   200 OK
  - response body: filtered tasks in JSON

failed:
  - status code:   400, 500
  - response body: JSON with error and status
*/
func (httpHandler *HTTPHandlers) handleGetFilteredTasks(w http.ResponseWriter, r *http.Request) {
	isCompletedStr := r.URL.Query().Get("completed")

	isCompletedBool, err := strconv.ParseBool(isCompletedStr)
	if err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		http.Error(w, response, http.StatusBadRequest)
		return
	}

	filteredTasks := httpHandler.taskStore.GetFilteredTasks(isCompletedBool)

	response, _ := ResponseDTO{
		Success: true,
		Result:  filteredTasks,
	}.ToJSON()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
pattern: /tasks/{id}
method:  PATCH
info:    ID in pattern, 'complete' param in body

successful:
  - status code:   200 ok
  - response body: JSON of changed task

failed:
  - status code:   400, 404, 500
  - response body: JSON with error and status
*/
func (httpHandler *HTTPHandlers) handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	ID, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}

	var dto TaskCompletedDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		http.Error(w, response, http.StatusBadRequest)
		return
	}

	task, err := httpHandler.taskStore.ToggleTask(ID, dto.Complete)
	if err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, response, http.StatusNotFound)
		} else {
			http.Error(w, response, http.StatusBadRequest)
		}
		return
	}

	response, _ := ResponseDTO{
		Success: true,
		Result:  task,
	}.ToJSON()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		fmt.Println(err.Error())
		return
	}
}

/*
pattern: /tasks/{id}
method:  DELETE
info     ID in pattern

successful:
  - status code:   204 No Content
  - response body: -

failed:
  - status code:   400, 404
  - response body: JSON with error and status
*/
func (httpHandler *HTTPHandlers) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	ID, err := parseIDFromRequest(w, r)
	if err != nil {
		return
	}

	if err := httpHandler.taskStore.DeleteTask(ID); err != nil {
		response, _ := ResponseDTO{
			Success: false,
			Message: err.Error(),
		}.ToString()
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, response, http.StatusNotFound)
		} else {
			http.Error(w, response, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
