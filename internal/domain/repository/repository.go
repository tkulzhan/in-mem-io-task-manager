package repository

import (
	"context"
	"fmt"
	"in-mem-io-task-manager/internal/domain/service"
	"in-mem-io-task-manager/internal/infrastructure/logger"
	"in-mem-io-task-manager/internal/infrastructure/storage"
)

type TaskRepository struct {
	l           *logger.Logger
	taskStorage *storage.TaskStorage
}

func NewTaskRepository(l *logger.Logger, taskStorage *storage.TaskStorage) (*TaskRepository, error) {
	if l == nil {
		return &TaskRepository{}, fmt.Errorf("task repository logger is not set")
	}
	if taskStorage == nil {
		return &TaskRepository{}, fmt.Errorf("task repository task storage is not set")
	}

	return &TaskRepository{
		l:           l,
		taskStorage: taskStorage,
	}, nil
}

func (r *TaskRepository) CreateTask(ctx context.Context, rt service.RunningTask) error {
	if err := r.taskStorage.CreateTask(rt); err != nil {
		r.l.Error("Failed to create task", map[string]interface{}{
			"task_id": rt.Task.GetID(),
			"error":   err.Error(),
		})
		return err
	}
	return nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	if err := r.taskStorage.DeleteTask(id); err != nil {
		r.l.Error("Failed to delete task", map[string]interface{}{
			"task_id": id,
			"error":   err.Error(),
		})
		return err
	}
	return nil
}

func (r *TaskRepository) GetTaskById(ctx context.Context, id string) (service.RunningTask, error) {
	rt, err := r.taskStorage.GetTaskById(id)
	if err != nil {
		r.l.Error("Failed to get task", map[string]interface{}{
			"task_id": id,
			"error":   err.Error(),
		})
		return service.RunningTask{}, err
	}
	return rt, nil
}
