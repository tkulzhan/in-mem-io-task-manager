package repository

import (
	"in-mem-io-task-manager/internal/infrastructure/logger"
	"in-mem-io-task-manager/internal/infrastructure/storage"
	"in-mem-io-task-manager/internal/domain/service"
)

type TaskRepository struct {
	l           *logger.Logger
	taskStorage *storage.TaskStorage
}

func NewTaskRepository(l *logger.Logger, taskStorage *storage.TaskStorage) *TaskRepository {
	if l == nil {
		panic("logger is not set")
	}
	if taskStorage == nil {
		panic("task storage is not set")
	}

	return &TaskRepository{
		l:           l,
		taskStorage: taskStorage,
	}
}

func (r *TaskRepository) CreateTask(task service.Task) error {
	if err := r.taskStorage.CreateTask(task); err != nil {
		r.l.Error("Failed to create task", map[string]interface{}{
			"task_id": task.GetID(),
			"error":   err.Error(),
		})
		return err
	}
	return nil
}

func (r *TaskRepository) DeleteTask(id string) error {
	if err := r.taskStorage.DeleteTask(id); err != nil {
		r.l.Error("Failed to delete task", map[string]interface{}{
			"task_id": id,
			"error":   err.Error(),
		})
		return err
	}
	return nil
}

func (r *TaskRepository) GetTaskById(id string) (service.Task, error) {
	task, err := r.taskStorage.GetTaskById(id)
	if err != nil {
		r.l.Error("Failed to get task by ID", map[string]interface{}{
			"task_id": id,
			"error":   err.Error(),
		})
		return nil, err
	}
	return task, nil
}
