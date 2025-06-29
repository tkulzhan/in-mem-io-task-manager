package service

import (
	"fmt"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type Configuration struct {
	L              *logger.Logger
	TaskRepository TaskRepository
}

func (conf Configuration) validate() error {
	if conf.TaskRepository == nil {
		return fmt.Errorf("task repository is not set")
	}
	if conf.L == nil {
		return fmt.Errorf("logger is not set")
	}

	return nil
}
