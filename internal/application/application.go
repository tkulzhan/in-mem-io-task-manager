package application

import (
	"context"
	"fmt"

	"in-mem-io-task-manager/internal/domain/repository"
	"in-mem-io-task-manager/internal/domain/service"
	httpserver "in-mem-io-task-manager/internal/infrastructure/http/server"
	"in-mem-io-task-manager/internal/infrastructure/logger"
	"in-mem-io-task-manager/internal/infrastructure/storage"
)

type Application struct {
	l *logger.Logger

	taskRepository *repository.TaskRepository

	taskManagerService *service.TaskManagerService

	httpServer *httpserver.Server
}

func New(ctx context.Context) (*Application, error) {
	conf, err := NewAppConfig()
	if err != nil {
		return nil, fmt.Errorf("new configuratuion: %w", err)
	}

	app := &Application{}

	if err := app.setLogger(conf); err != nil {
		return nil, fmt.Errorf("set logger: %w", err)
	}

	if err := app.setRepositories(); err != nil {
		return nil, fmt.Errorf("set repositories: %w", err)
	}

	if err := app.setService(); err != nil {
		return nil, fmt.Errorf("set service: %w", err)
	}

	if err := app.setServer(conf); err != nil {
		return nil, fmt.Errorf("set server: %w", err)
	}

	return app, nil
}

func (a *Application) setLogger(conf *Configuration) error {
	logger := logger.New(conf.LogLevel)

	a.l = logger

	return nil
}

func (a *Application) setRepositories() error {
	taskStorage := storage.NewTaskStorage()
	taskRepository, err := repository.NewTaskRepository(a.l, taskStorage)
	if err != nil {
		return fmt.Errorf("creating task repository: %w", err)
	}

	a.taskRepository = taskRepository

	return nil
}

func (a *Application) setService() error {
	taskManagerService, err := service.NewTaskManagerService(service.Configuration{
		L:              a.l,
		TaskRepository: a.taskRepository,
	})
	if err != nil {
		return fmt.Errorf("creating task Manager Service: %w", err)
	}

	a.taskManagerService = taskManagerService

	return nil
}

func (a *Application) setServer(conf *Configuration) error {
	// REST
	httpServer, err := httpserver.New(httpserver.Configuration{
		L:                      a.l,
		RestAddress:            conf.RestAddress,
		ReadTimeoutSeconds:     conf.ReadTimeoutSeconds,
		WriteTimeoutSeconds:    conf.WriteTimeoutSeconds,
		ExposeAPISpecification: conf.ExposeAPISpecification,
		TaskManagerService:     a.taskManagerService,
	})
	if err != nil {
		return fmt.Errorf("creating HTTP server: %w", err)
	}

	a.httpServer = httpServer

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	if err := a.httpServer.Start(); err != nil {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}

func (a *Application) Close(ctx context.Context) error {
	if a.httpServer != nil {
		if err := a.httpServer.Close(); err != nil {
			a.l.Error(fmt.Sprintf("closing HTTP server: %s", err.Error()), nil)
		}
	}

	return nil
}
