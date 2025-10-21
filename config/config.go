package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres Postgres
	Server   Server
	Logger   Logger
}

type Server struct {
	Addr string `env:"HTTP_ADDRESS" env-required:"true"`
}

type Postgres struct {
	DSN         string        `env:"PG_DSN" env-required:"true"`
	ConnTimeOut time.Duration `env:"PG_CONN_TIMEOUT" env-default:"5s"`
}

type Logger struct {
	Level string `env:"LOGGER_LEVEL" env-default:"info"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read envs: " + err.Error())
	}

	return &cfg
}
