package httpserver

import "net/http"

func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a task
}

func (s *Server) GetTaskById(w http.ResponseWriter, r *http.Request, taskId string) {
	// Implementation for getting a task by ID
}

func (s *Server) DeleteTaskById(w http.ResponseWriter, r *http.Request, taskId string) {
	// Implementation for deleting a task by ID
}
