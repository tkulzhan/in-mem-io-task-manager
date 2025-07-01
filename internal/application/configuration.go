package application

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	configPath = "./config/.env"
)

type Configuration struct {
	LogLevel               string `env:"LOG_LEVEL"`
	RestAddress            string `env:"REST_ADDRESS"`
	ReadTimeoutSeconds     uint   `env:"HTTP_READ_TIME_OUT"`
	WriteTimeoutSeconds    uint   `env:"HTTP_WRITE_TIME_OUT"`
	ExposeAPISpecification bool   `env:"HTTP_EXPOSE_API_SPECIFICATION"`
	MaxConcurrentTasks     int    `env:"MAX_CONCURRENT_TASKS"`
}

func NewAppConfig() (*Configuration, error) {
	cfg := &Configuration{}

	if err := godotenv.Load(configPath); err != nil {
		var perr *fs.PathError
		if errors.As(err, &perr) {
			setDefaults(cfg)
			return cfg, nil
		} else {
			return nil, fmt.Errorf("loading configuration: %w", err)
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	applyDefaultsIfEmpty(cfg)

	return cfg, nil
}

func setDefaults(cfg *Configuration) {
	cfg.RestAddress = "0.0.0.0:8080"
	cfg.LogLevel = "debug"
	cfg.ExposeAPISpecification = true
	cfg.ReadTimeoutSeconds = 30
	cfg.WriteTimeoutSeconds = 30
	cfg.MaxConcurrentTasks = 1000
}

func applyDefaultsIfEmpty(cfg *Configuration) {
	if cfg.RestAddress == "" {
		cfg.RestAddress = "0.0.0.0:8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "debug"
	}
	if cfg.ReadTimeoutSeconds == 0 {
		cfg.ReadTimeoutSeconds = 30
	}
	if cfg.WriteTimeoutSeconds == 0 {
		cfg.WriteTimeoutSeconds = 30
	}
	if cfg.MaxConcurrentTasks == 0 {
		cfg.MaxConcurrentTasks = 1000
	}
}
