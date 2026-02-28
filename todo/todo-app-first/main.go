package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	utils "todo-app/utils"
)

type Task struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Completed     bool       `json:"completed"`
	CreatedTime   time.Time  `json:"created_time"`
	CompletedTime *time.Time `json:"completed_time"`
}

var PATH = "/tasks/"

var mutex = sync.Mutex{}

var tasks = make(map[int]Task)

func createTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error during reading request body: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error deserializing  request body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		utils.RespondJSONError(w, "Title is required", http.StatusBadRequest)
		return
	}

	task.ID = utils.GetNextTaskID()
	task.CreatedTime = time.Now()
	task.CompletedTime = nil

	mutex.Lock()
	tasks[task.ID] = task
	mutex.Unlock()

	taskJSON, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("Error serialize task: %s", err.Error())
		utils.RespondJSONError(w, errMsg, http.StatusInternalServerError)
		return
	}

	utils.RespondJSONSuccess(w, http.StatusCreated, taskJSON)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	completeFilterStr := r.URL.Query().Get("completed")
	filteredTasks := []Task{}

	if completeFilterStr != "" {
		completeFilter, err := strconv.ParseBool(completeFilterStr)
		if err != nil {
			utils.RespondJSONError(w, fmt.Sprintf("Invalid '%s' parameter. Use boolean type", completeFilterStr), http.StatusBadRequest)
			return
		}

		mutex.Lock()
		for _, task := range tasks {
			if task.Completed == completeFilter {
				filteredTasks = append(filteredTasks, task)
			}
		}
		mutex.Unlock()

	} else {
		mutex.Lock()

		for _, task := range tasks {
			filteredTasks = append(filteredTasks, task)
		}
		mutex.Unlock()
	}

	tasksJSON, err := json.MarshalIndent(filteredTasks, "", "  ")
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error during serializing tasks: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	utils.RespondJSONSuccess(w, http.StatusOK, tasksJSON)
}

func getTask(w http.ResponseWriter, IDStr string) {
	ID, err := utils.ConvertIDStrToInt(IDStr)
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error converting ID to number: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	task, ok := tasks[ID]
	if !ok {
		utils.RespondJSONError(w, fmt.Sprintf("Task ID %d is not found", ID), http.StatusNotFound)
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	taskJSON, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error during the serialization of a found task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	utils.RespondJSONSuccess(w, http.StatusOK, taskJSON)
}

func completeTask(w http.ResponseWriter, IDStr string) {
	ID, err := utils.ConvertIDStrToInt(IDStr)
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error during convert ID to number: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	task, ok := tasks[ID]

	if ok {
		task.Completed = true
		timeNow := time.Now()
		task.CompletedTime = &timeNow
		tasks[ID] = task
		mutex.Unlock()
	} else {
		utils.RespondJSONError(w, fmt.Sprintf("Task ID %d is not found", ID), http.StatusNotFound)
		mutex.Unlock()
		return
	}

	taskJSON, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		utils.RespondJSONError(w, "Error in serializing the found task", http.StatusInternalServerError)
		return
	}
	utils.RespondJSONSuccess(w, http.StatusOK, taskJSON)
}

func deleteTask(w http.ResponseWriter, IDStr string) {
	ID, err := utils.ConvertIDStrToInt(IDStr)
	if err != nil {
		utils.RespondJSONError(w, fmt.Sprintf("Error in converting ID to number: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	delete(tasks, ID)
	mutex.Unlock()

	utils.RespondJSONSuccess(w, http.StatusNoContent, nil)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request",
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"content_type", r.Header.Get("Content-Type"),
		"content_length", r.Header.Get("Content-Length"),
		"time", time.Now(),
	)

	restPath := strings.TrimPrefix(r.URL.Path, PATH)
	IDString := strings.SplitN(restPath, "/", 2)[0]

	switch r.Method {
	case http.MethodGet:
		if IDString == "" {
			getTasks(w, r)
		} else {
			getTask(w, IDString)
		}
	case http.MethodPost:
		if IDString != "" {
			utils.RespondJSONError(w, "POST method does not accept additional information", http.StatusBadRequest)
			return
		}
		createTask(w, r)
	case http.MethodPatch:
		if IDString == "" {
			utils.RespondJSONError(w, "ID of task is required", http.StatusBadRequest)
			return
		}
		completeTask(w, IDString)
	case http.MethodDelete:
		if IDString == "" {
			utils.RespondJSONError(w, "ID of task is required", http.StatusBadRequest)
			return
		}

		deleteTask(w, IDString)
	default:
		utils.RespondJSONError(w, "Method not allowed", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc(PATH, taskHandler)

	fmt.Println("Start app on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}

}
