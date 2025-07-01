package service

import (
	"context"
	"fmt"
	"in-mem-io-task-manager/internal/domain/entity"
	"in-mem-io-task-manager/internal/infrastructure/errors"
	"in-mem-io-task-manager/internal/infrastructure/http/server/generated"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type TaskManagerService struct {
	l              *logger.Logger
	taskRepository TaskRepository
	semaphore      chan struct{}
}

type RunningTask struct {
	Task   Task
	Cancel context.CancelFunc
}

func NewTaskManagerService(conf Configuration) (*TaskManagerService, error) {
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}
	return &TaskManagerService{
		l:              conf.L,
		taskRepository: conf.TaskRepository,
		semaphore:      make(chan struct{}, conf.MaxConcurrentTasks),
	}, nil
}

func (s *TaskManagerService) GetTaskById(ctx context.Context, id string) (Task, error) {
	task, err := s.taskRepository.GetTaskById(ctx, id)
	if err != nil {
		return nil, err
	}
	return task.Task, nil
}

func (s *TaskManagerService) CreateTask(ctx context.Context, task Task) error {
	execCtx, cancel := context.WithCancel(context.Background())

	rt := RunningTask{
		Task:   task,
		Cancel: cancel,
	}

	if err := s.taskRepository.CreateTask(ctx, rt); err != nil {
		cancel()
		return err
	}

	go func() {
		s.semaphore <- struct{}{}
		defer func() { <-s.semaphore }()

		err := task.Execute(execCtx, s.l)
		if err != nil {
			s.l.Error("task execution failed", map[string]interface{}{
				"task_id": task.GetID(),
				"error":   err.Error(),
			})
		} else {
			s.l.Debug("task executed successfully", map[string]interface{}{
				"task_id": task.GetID(),
			})
		}
	}()

	return nil
}

func (s *TaskManagerService) DeleteTask(ctx context.Context, id string) error {
	rt, err := s.taskRepository.GetTaskById(ctx, id)
	if err != nil {
		return err
	}

	rt.Cancel()

	if err := s.taskRepository.DeleteTask(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *TaskManagerService) ToTaskObject(ctx context.Context, req generated.CreateTaskRequest) (Task, error) {
	taskType := "default"
	if req.Type != nil && *req.Type != "" {
		taskType = *req.Type
	}

	switch taskType {
	case "default":
		input, err := parseDefaultTaskInput(req.Data)
		if err != nil {
			return nil, err
		}
		return entity.NewDefaultTask(input.Title, deref(input.Description)), nil

	default:
		return nil, errors.NewBadRequestError(fmt.Sprintf("unsupported task type: %s", taskType))
	}
}
