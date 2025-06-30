package httpserver

import (
	"fmt"
	"in-mem-io-task-manager/internal/domain/service"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type Configuration struct {
	L                      *logger.Logger
	RestAddress            string
	ReadTimeoutSeconds     uint
	WriteTimeoutSeconds    uint
	ExposeAPISpecification bool
	TaskManagerService     *service.TaskManagerService
}

func (conf Configuration) validate() error {
	if conf.RestAddress == "" {
		return fmt.Errorf("http port is required")
	}
	if conf.L == nil {
		return fmt.Errorf("logger is not set")
	}
	if conf.TaskManagerService == nil {
		return fmt.Errorf("task manager service is not set")
	}

	return nil
}
