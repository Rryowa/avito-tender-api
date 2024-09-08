package config

import "time"

type DbConfig struct {
	User     string        `env:"POSTGRES_USERNAME"`
	Password string        `env:"POSTGRES_PASSWORD"`
	Host     string        `env:"POSTGRES_HOST"`
	Port     string        `env:"POSTGRES_PORT"`
	DBName   string        `env:"POSTGRES_DATABASE"`
	Attempts int           `env:"ATTEMPTS"`
	Timeout  time.Duration `env:"TIMEOUT"`
}