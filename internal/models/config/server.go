package config

import "time"

type ServerConfig struct {
	ServerAddr      string        `env:"SERVER_ADDRESS"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"`
	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT"`
}