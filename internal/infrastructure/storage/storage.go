package storage

import (
	"in-mem-io-task-manager/internal/domain/service"
	"in-mem-io-task-manager/internal/infrastructure/errors"
	"sync"
)

type TaskStorage struct {
	tasks map[string]service.Task
	mu    *sync.Mutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks: make(map[string]service.Task),
		mu:    &sync.Mutex{},
	}
}

func (ts *TaskStorage) CreateTask(task service.Task) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[task.GetID()]; exists {
		return errors.NewBadRequestError("task with this ID already exists")
	}
	ts.tasks[task.GetID()] = task
	return nil
}

func (ts *TaskStorage) DeleteTask(id string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		return errors.NewNotFoundError("task not found")
	}
	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStorage) GetTaskById(id string) (service.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, errors.NewNotFoundError("task not found")
	}
	return task, nil
}
