package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/saibaend/template-svc/internal/db"
	"time"
)

type Env struct {
	LogLevel          string        `env:"LOG_LEVEL,required"`
	HttpPort          string        `env:"HTTP_PORT,required"`
	HttpClientTimeout time.Duration `env:"HTTP_CLIENT_TIMEOUT,required"`
	DbConfig          db.Config
}

func New() (config Env, err error) {
	if err = env.Parse(&config); err != nil {
		return config, err
	}

	return config, nil
}
