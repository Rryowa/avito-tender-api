package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/http"
	"zadanie-6105/internal/util"
)

func (d *Database) CheckUserExists(ctx context.Context, username string) error {
	const op = "storage.IsUserExists"

	query := `SELECT 1 
				FROM employee
				WHERE username = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusUnauthorized, Msg: util.Unauthorized}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) CheckUserByIDExists(ctx context.Context, id string) error {
	const op = "storage.IsUserByIdExists"

	query := `SELECT 1 
				FROM employee
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, id).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusUnauthorized, Msg: util.Unauthorized}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) CheckTenderExists(ctx context.Context, tenderID string) error {
	const op = "storage.IsTenderExists"

	query := `SELECT 1 
				FROM tender
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) CheckBidExists(ctx context.Context, bidID string) error {
	const op = "storage.IsTenderExists"

	query := `SELECT 1 
				FROM bid
				WHERE id = $1;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// CheckUserBidAuthor Проверяет, что пользователь является автором Предложения.
func (d *Database) CheckUserBidAuthor(ctx context.Context, bidID, requestedUser string) error {
	const op = "storage.CheckUserBidAuthor"

	query := `SELECT 1
				FROM bid
				WHERE id = $1 AND author_username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidID, requestedUser).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ValidateUserResponsible Проверяет принадлежность пользователя организации, которая открыла тендер.
func (d *Database) ValidateUserResponsible(ctx context.Context, tenderID, username string) error {
	const op = "storage.ValidateUserResponsible"

	query := `SELECT 1
				FROM organization_responsible o
				JOIN employee e ON o.user_id = e.id
				JOIN tender t ON t.organization_id = o.organization_id
				WHERE t.id = $1 AND e.username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderID, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ValidateUserResponsibleBidID Проверяет принадлежность пользователя организации, которая открыла тендер.
func (d *Database) ValidateUserResponsibleBidID(ctx context.Context, bidID, username string) error {
	const op = "storage.ValidateUserResponsibleBidID"

	query := `SELECT 1
				FROM organization_responsible o
				JOIN employee e ON o.user_id = e.id
				JOIN bid b ON b.organization_id = o.organization_id
				WHERE b.id = $1 AND e.username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidID, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ValidateUserResponsibleUserID Проверяет принадлежность пользователя организации, которая открыла тендер.
func (d *Database) ValidateUserResponsibleUserID(ctx context.Context, userID, tenderID string) error {
	const op = "storage.ValidateUsersPrivilegesById"
	query := `SELECT 1
				FROM organization_responsible o
				JOIN tender t ON o.organization_id = t.organization_id
				WHERE o.user_id = $1 AND t.id = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, userID, tenderID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ValidateUserResponsibleOrgID Проверяет принадлежность пользователя организации, которая открыла тендер.
func (d *Database) ValidateUserResponsibleOrgID(ctx context.Context, orgID, username string) error {
	const op = "storage.ValidateUserResponsibleOrgID"
	query := `SELECT 1
				FROM organization_responsible o
				JOIN employee e ON o.user_id = e.id
				WHERE organization_id = $1 AND e.username = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, orgID, username).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) CheckBidVersionExists(ctx context.Context, bidID string, version int32) error {
	const op = "storage.IsBidVersionExists"

	query := `SELECT 1 
				FROM bid_history
				WHERE bid_id = $1 AND version = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, bidID, version).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusNotFound, Msg: util.VersionNotFound}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (d *Database) CheckTenderVersionExists(ctx context.Context, tenderID string, version int32) error {
	const op = "storage.IsTenderVersionExists"

	query := `SELECT 1 
				FROM tender_history
				WHERE tender_id = $1 AND version = $2;`

	var dummy int
	err := d.Pool.QueryRow(ctx, query, tenderID, version).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return util.MyResponseError{Status: http.StatusNotFound, Msg: util.VersionNotFound}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}