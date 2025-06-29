package httpserver

import (
	"fmt"
	"in-mem-io-task-manager/internal/infrastructure/logger"
)

type Configuration struct {
	L                      *logger.Logger
	Address                string
	ReadTimeoutSeconds     uint
	WriteTimeoutSeconds    uint
	ExposeAPISpecification bool
}

func (conf Configuration) validate() error {
	if conf.Address == "" {
		return fmt.Errorf("address is required")
	}
	if conf.L == nil {
		return fmt.Errorf("logger is not set")
	}

	return nil
}
