package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"
)

func RespondJSONSuccess(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func ConvertIDStrToInt(IDStr string) (int, error) {
	IDInt, err := strconv.Atoi(IDStr)
	if err != nil {
		return 0, err
	}
	return IDInt, nil
}

var taskID atomic.Int64

func GetNextTaskID() int {
	taskID.Add(1)
	return int(taskID.Load())
}
