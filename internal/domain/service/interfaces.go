package service

import (
	"context"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type Task interface {
	Execute(ctx context.Context, l *logger.Logger) error
	GetID() string
}

type TaskRepository interface {
	CreateTask(ctx context.Context, rt RunningTask) error
	DeleteTask(ctx context.Context, id string) error
	GetTaskById(ctx context.Context, id string) (RunningTask, error)
}
