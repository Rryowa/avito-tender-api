package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"zadanie-6105/internal/models"
)

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

	rows, err := d.Pool.Query(ctx, query, models.Published, serviceTypes, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var tenders []models.Tender
	if err = pgxscan.ScanAll(&tenders, rows); err != nil {
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
				WHERE creator_username = $1
				ORDER BY name
				OFFSET $2
				FETCH NEXT $3 ROWS ONLY`

	rows, err := d.Pool.Query(ctx, query, username, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var tenders []models.Tender
	if err = pgxscan.ScanAll(&tenders, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return tenders, err
}

func (d *Database) GetTenderStatus(ctx context.Context, tenderID, username string) (string, error) {
	const op = "storage.GetTenderStatus"

	query := `SELECT status
				FROM tender
				WHERE id = $1
  				AND creator_username = $2`

	rows, err := d.Pool.Query(ctx, query, tenderID, username)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var status string
	if err = pgxscan.ScanOne(&status, rows); err != nil {
		return "", fmt.Errorf("%s: %w", op2, err)
	}

	return status, err
}

func (d *Database) UpdateTenderStatus(ctx context.Context, tenderID, status, username string) (models.Tender, error) {
	const op = "storage.UpdateTenderStatus"

	query := `UPDATE tender
				SET status = $1
				WHERE id = $2 and creator_username = $3
				RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at, updated_at;
	`

	rows, err := d.Pool.Query(ctx, query, status, tenderID, username)
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

func (d *Database) EditTender(ctx context.Context, tender *models.Tender, tenderID, username string) (models.Tender, error) {
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

	rows, err := d.Pool.Query(ctx, query, tender.Name, tender.Description, tender.ServiceType, tenderID, username)
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

func (d *Database) RollbackTender(ctx context.Context, tenderID string, version int32, username string) (models.Tender, error) {
	const op = "storage.RollbackTender"

	query := `SELECT * FROM rollback_tender_version($1::UUID, $2::INT, $3::VARCHAR)`

	rows, err := d.Pool.Query(ctx, query, tenderID, version, username)
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