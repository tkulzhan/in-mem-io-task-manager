package storage

import (
	"in-mem-io-task-manager/internal/domain/service"
	"in-mem-io-task-manager/internal/infrastructure/errors"
	"sync"
)

type TaskStorage struct {
	tasks map[string]service.RunningTask
	mu    *sync.Mutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]service.RunningTask),
		mu:    &sync.Mutex{},
	}
}

func (ts *TaskStorage) CreateTask(rt service.RunningTask) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[rt.Task.GetID()]; exists {
		return errors.NewBadRequestError("task with this ID already exists")
	}

	ts.tasks[rt.Task.GetID()] = rt
	return nil
}

func (ts *TaskStorage) DeleteTask(id string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, exists := ts.tasks[id]
	if !exists {
		return errors.NewNotFoundError("task not found")
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStorage) GetTaskById(id string) (service.RunningTask, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	rt, exists := ts.tasks[id]
	if !exists {
		return service.RunningTask{}, errors.NewNotFoundError("task not found")
	}

	return rt, nil
}
