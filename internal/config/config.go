package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP     HTTP
		Database Database
	}

	HTTP struct {
		PORT string `env:"HTTP_PORT,required"`
	}

	Database struct {
		ConnectionString string `env:"DATABASE_URL,required"`
	}
)

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
