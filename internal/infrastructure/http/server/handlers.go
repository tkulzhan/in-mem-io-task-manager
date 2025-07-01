package httpserver

import (
	"encoding/json"
	"in-mem-io-task-manager/internal/infrastructure/errors"
	"in-mem-io-task-manager/internal/infrastructure/http/server/generated"
	"net/http"
)

func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request generated.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errors.HandleError(w, errors.NewBadRequestError("invalid request body"))
		return
	}

	task, err := s.taskManagerService.ToTaskObject(ctx, request)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	if err := s.taskManagerService.CreateTask(ctx, task); err != nil {
		errors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data, err := task.MarshalJSON()
	if err != nil {
		errors.HandleError(w, errors.NewInternalError("failed to marshal response"))
		return
	}

	if _, err := w.Write(data); err != nil {
		errors.HandleError(w, errors.NewInternalError("failed to send response"))
	}
}

func (s *Server) GetTaskById(w http.ResponseWriter, r *http.Request, taskId string) {
	ctx := r.Context()

	task, err := s.taskManagerService.GetTaskById(ctx, taskId)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data, err := task.MarshalJSON()
	if err != nil {
		errors.HandleError(w, errors.NewInternalError("failed to marshal response"))
		return
	}

	if _, err := w.Write(data); err != nil {
		errors.HandleError(w, errors.NewInternalError("failed to send response"))
	}
}

func (s *Server) DeleteTaskById(w http.ResponseWriter, r *http.Request, taskId string) {
	ctx := r.Context()

	err := s.taskManagerService.DeleteTask(ctx, taskId)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
