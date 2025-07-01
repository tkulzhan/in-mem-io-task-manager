package service

import (
	"context"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type Task interface {
	Execute(ctx context.Context, l *logger.Logger) error
	GetID() string
	MarshalJSON() ([]byte, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, rt RunningTask) error
	DeleteTask(ctx context.Context, id string) error
	GetTaskById(ctx context.Context, id string) (RunningTask, error)
}
