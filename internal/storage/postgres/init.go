package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"zadanie-6105/internal/models/config"
	"zadanie-6105/internal/storage"
	"zadanie-6105/internal/util"
)

type Database struct {
	Pool      *pgxpool.Pool
	zapLogger *zap.SugaredLogger
}

func NewPostgresRepository(ctx context.Context, cfg *config.DBConfig, zap *zap.SugaredLogger) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			zap.Fatalln(err, "db connection error")
		}

		return nil
	}, cfg.Attempts, cfg.Timeout)

	if err != nil {
		zap.Fatalln(err, "DoWithTries error")
	}
	zap.Infoln("Connected to db")

	return &Database{
		Pool:      pool,
		zapLogger: zap,
	}
}