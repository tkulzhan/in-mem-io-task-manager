package service

import (
	"fmt"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type TaskManagerService struct {
	l              *logger.Logger
	taskRepository TaskRepository
}

func NewTaskManagerService(conf Configuration) (*TaskManagerService, error) {
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}
	return &TaskManagerService{
		l:              conf.L,
		taskRepository: conf.TaskRepository,
	}, nil
}

func (s *TaskManagerService) GetTaskById(id string) (Task, error) {
	return nil, nil
}

func (s *TaskManagerService) CreateTask(task Task) error {
	return nil
}

func (s *TaskManagerService) DeleteTask(id string) error {
	return nil
}
