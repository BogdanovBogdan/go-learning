package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandler *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandler,
	}
}

func (server HTTPServer) StartHTTPServer() error {
	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(server.httpHandlers.handleCreateTask)
	router.Path("/tasks").Methods("GET").Queries("completed", "{completed}").HandlerFunc(server.httpHandlers.handleGetFilteredTasks)
	router.Path("/tasks").Methods("GET").HandlerFunc(server.httpHandlers.handleGetAllTasks)
	router.Path("/tasks/{id}").Methods("GET").HandlerFunc(server.httpHandlers.handleGetTask)
	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(server.httpHandlers.handleCompleteTask)
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(server.httpHandlers.handleDeleteTask)

	fmt.Println("Server starts on port 8080...")
	return http.ListenAndServe(":8080", router)
}
