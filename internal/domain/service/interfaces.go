package service

import (
	"context"
	"in-mem-io-task-manager/internal/domain/entity"
)

type Task interface {
	Execute(ctx context.Context) error
	GetID() string
	GetStatus() entity.TaskStatus
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, id string) error
	GetTaskById(ctx context.Context, id string) (*Task, error)
}
