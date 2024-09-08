package util

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strconv"
	"time"
	"zadanie-6105/internal/models/config"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
}

func NewServerConfig() *config.ServerConfig {
	return &config.ServerConfig{
		ServerAddr:    os.Getenv("SERVER_ADDRESS"),
		TelemetryAddr: os.Getenv("TELEMETRY_ADDRESS"),
	}
}

func NewDbConfig() *config.DbConfig {
	attempts, err := strconv.Atoi(os.Getenv("ATTEMPTS"))
	if err != nil {
		log.Fatalf("err converting ATTEMPTS: %v\n", err)
	}
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		log.Fatalf("Error parsing TIMEOUT: %v\n", err)
	}

	return &config.DbConfig{
		User:     os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DATABASE"),
		Attempts: attempts,
		Timeout:  timeout,
	}
}

func NewZapLogger() *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	sugar := logger.Sugar()
	sugar.Sync()
	return sugar
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}