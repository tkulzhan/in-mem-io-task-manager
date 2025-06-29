package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"in-mem-io-task-manager/internal/domain/service"
	"in-mem-io-task-manager/internal/infrastructure/http/server/generated"
	"net/http"
	"time"
)

const v1BaseURL string = "/api"

type Server struct {
	address            string
	srv                *http.Server
	taskManagerService service.TaskManagerService
}

func New(ctx context.Context, conf Configuration) (*Server, error) {
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	router := http.NewServeMux()

	if conf.ExposeAPISpecification {
		router.Handle(v1BaseURL+"/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			swagger, err := generated.GetSwagger()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			json.NewEncoder(w).Encode(swagger)
		}))
	}

	srv := &Server{
		address: conf.Address,
	}

	handler := generated.HandlerWithOptions(srv, generated.StdHTTPServerOptions{
		BaseURL:    v1BaseURL,
		BaseRouter: router,
	})

	server := &http.Server{
		Addr:         conf.Address,
		ReadTimeout:  time.Duration(conf.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeoutSeconds) * time.Second,
		Handler:      handler,
	}

	srv.srv = server

	return srv, nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("launching HTTP server error: %w", err)
		}
	}

	return nil
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("closing error: %w", err)
	}

	return nil
}
