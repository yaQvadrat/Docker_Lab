package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
	}

	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"VERSION"`
	}

	HTTP struct {
		Address string `env-required:"true" env:"SERVER_ADDRESS"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	PG struct {
		MaxPoolSize int    `env-required:"true" env:"MAX_POOL_SIZE"`
		URL         string `env-required:"true" env:"POSTGRES_CONN"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	// Reading config from .yaml file and enviroment (env variables more important)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("config - New - cleanenv.ReadEnv: %w", err)
	}

	return cfg, nil
}
