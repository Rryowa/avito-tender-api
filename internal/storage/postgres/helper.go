package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func NewPostgresRepository(ctx context.Context, cfg *config.DbConfig, zap *zap.SugaredLogger) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
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

func (d *Database) IsUserExists(ctx context.Context, username string) error {
	const op = "storage.IsUserExists"

	query := `SELECT 1 
				FROM employee
				WHERE username = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) IsUserByIdExists(ctx context.Context, id string) error {
	const op = "storage.IsUserByIdExists"

	query := `SELECT 1 
				FROM employee
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, id).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) IsTenderExists(ctx context.Context, tenderId string) error {
	const op = "storage.IsTenderExists"

	query := `SELECT 1 
				FROM tender
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderId).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) IsBidExists(ctx context.Context, bidId string) error {
	const op = "storage.IsTenderExists"

	query := `SELECT 1 
				FROM bid
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidId).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ValidateUsersPrivileges Проверяет принадлежность пользователя организации, которая открыла тендер
func (d *Database) ValidateUsersPrivileges(ctx context.Context, tenderId, requestedUser string) error {
	const op = "storage.IsUserInOrganization"

	query := `SELECT 1
				FROM organization_responsible o
				JOIN employee e ON o.user_id = e.id
				JOIN tender t ON t.organization_id = o.organization_id
				WHERE t.id = $1 AND e.username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderId, requestedUser).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) ValidateUsersPrivilegesUserId(ctx context.Context, id, tenderId string) error {
	const op = "storage.ValidateUsersPrivilegesById"
	query := `SELECT 1
				FROM organization_responsible o
				JOIN tender t ON o.organization_id = t.organization_id
				WHERE o.user_id = $1 AND t.id = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, id, tenderId).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) ValidateUsersPrivilegesOrgId(ctx context.Context, orgId, username string) error {
	const op = "storage.ValidateUsersPrivilegesOrgId"
	query := `SELECT 1
				FROM organization_responsible o
				JOIN employee e ON o.user_id = e.id
				WHERE organization_id = $1 AND e.username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, orgId, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// TODO: add zap logging of errors
func (d *Database) ValidateUsersPrivilegesBidId(ctx context.Context, bidId, username string) error {
	const op = "storage.ValidateUsersPrivilegesBidId"

	query := `
		SELECT 1
		FROM bid b
		JOIN tender t ON b.tender_id = t.id
		JOIN organization_responsible o ON t.organization_id = o.organization_id
		WHERE b.id = $1 AND t.creator_username = $2;`

	var exists int
	err := d.Pool.QueryRow(ctx, query, bidId, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) IsBidVersionExists(ctx context.Context, bidId string, version int32) error {
	const op = "storage.IsBidVersionExists"

	query := `SELECT 1 
				FROM bid_history
				WHERE bid_id = $1 AND version = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidId, version).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) IsTenderVersionExists(ctx context.Context, tenderId string, version int32) error {
	const op = "storage.IsTenderVersionExists"

	query := `SELECT 1 
				FROM tender_history
				WHERE tender_id = $1 AND version = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderId, version).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}