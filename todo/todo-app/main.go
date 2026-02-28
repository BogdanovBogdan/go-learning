package main

import (
	"log"
	"todo-app/server"
	"todo-app/todo"
)

func main() {
	todoList := todo.NewTodoList()
	httpHandlers := server.NewHTTPHandlers(todoList)
	httpServer := server.NewHTTPServer(httpHandlers)

	if err := httpServer.StartHTTPServer(); err != nil {
		log.Fatal(err)
	}
}
