package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/models/config"
	"zadanie-6105/internal/storage"
	"zadanie-6105/internal/util"
)

type Database struct {
	Pool      *pgxpool.Pool
	zapLogger *zap.SugaredLogger
}

func (d *Database) GetTenders(ctx context.Context, offset, limit int32, serviceTypes []string) ([]models.Tender, error) {
	const op = "storage.GetTenders"

	query := `SELECT id, name, description, service_type, status, version, organization_id, creator_username, created_at
				FROM tender
				WHERE status = $1
  				AND ($2::VARCHAR[] IS NULL OR service_type::VARCHAR = ANY($2::VARCHAR[]))
				ORDER BY name
				OFFSET $3
				FETCH NEXT $4 ROWS ONLY`

	if len(serviceTypes) == 0 {
		serviceTypes = nil
	}

	//TODO: replace with published
	rows, err := d.Pool.Query(ctx, query, models.Created, serviceTypes, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var tenders []models.Tender
	if err := pgxscan.ScanAll(&tenders, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return tenders, err
}

func (d *Database) CreateTender(ctx context.Context, tender *models.Tender) (models.Tender, error) {
	const op = "storage.CreateTender"

	query := `INSERT INTO tender (name, description, service_type, status, organization_id, creator_username)
				VALUES ($1,	$2, $3, $4, $5, $6)
				RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at;`

	rows, err := d.Pool.Query(ctx, query, tender.Name, tender.Description, tender.ServiceType, tender.Status, tender.OrganizationID, tender.CreatorUsername)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newTender models.Tender
	err = pgxscan.ScanOne(&newTender, rows)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newTender, nil
}

func (d *Database) GetUserTenders(ctx context.Context, offset, limit int32, username string) ([]models.Tender, error) {
	const op = "storage.GetUserTenders"

	query := `SELECT id, name, description, service_type, status, version, organization_id, creator_username, created_at
				FROM tender
				WHERE status = $1
  				AND creator_username = $2
				ORDER BY name
				OFFSET $3
				FETCH NEXT $4 ROWS ONLY`

	//TODO: replace with published
	rows, err := d.Pool.Query(ctx, query, models.Created, username, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var tenders []models.Tender
	if err := pgxscan.ScanAll(&tenders, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return tenders, err
}

func (d *Database) GetTenderStatus(ctx context.Context, tenderId, username string) (string, error) {
	const op = "storage.GetTenderStatus"

	query := `SELECT status
				FROM tender
				WHERE id = $1
  				AND creator_username = $2`

	//TODO: replace with published
	rows, err := d.Pool.Query(ctx, query, tenderId, username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var status string
	if err := pgxscan.ScanOne(&status, rows); err != nil {
		return "", fmt.Errorf("%s: %w", op2, err)
	}

	return status, err
}

func (d *Database) UpdateTenderStatus(ctx context.Context, tenderId, status, username string) (models.Tender, error) {
	const op = "storage.UpdateTenderStatus"

	query := `UPDATE tender
				SET status = $1
				WHERE id = $2 and creator_username = $3
				RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at, updated_at;
	`

	rows, err := d.Pool.Query(ctx, query, status, tenderId, username)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var updatedTender models.Tender
	err = pgxscan.ScanOne(&updatedTender, rows)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op2, err)
	}

	return updatedTender, nil
}

func (d *Database) EditTender(ctx context.Context, tender *models.Tender, tenderId, username string) (models.Tender, error) {
	const op = "storage.EditTender"

	query := `UPDATE tender 
				SET 
					name = COALESCE(NULLIF($1, ''), name), 
					description = COALESCE(NULLIF($2, ''), description), 
					service_type = CASE 
						WHEN $3 = '' THEN service_type
						ELSE $3::service_type
					END 
				WHERE id = $4 AND creator_username = $5
				RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at, updated_at;
		`

	rows, err := d.Pool.Query(ctx, query, tender.Name, tender.Description, tender.ServiceType, tenderId, username)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newTender models.Tender
	err = pgxscan.ScanOne(&newTender, rows)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newTender, nil
}

func (d *Database) RollbackTender(ctx context.Context, tenderId string, version int32, username string) (models.Tender, error) {
	const op = "storage.RollbackTender"

	query := `SELECT * FROM rollback_tender_version($1::UUID, $2::INT, $3::VARCHAR)`

	rows, err := d.Pool.Query(ctx, query, tenderId, version, username)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newTender models.Tender
	err = pgxscan.ScanOne(&newTender, rows)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newTender, nil
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