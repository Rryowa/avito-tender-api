package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"zadanie-6105/internal/models"
)

func (d *Database) CreateBid(ctx context.Context, bid *models.Bid) (models.Bid, error) {
	const op = "storage.CreateBid"

	query := `INSERT INTO bid (name, description, tender_id, author_type, author_id)
				VALUES ($1,	$2, $3, $4, $5)
				RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at;`

	rows, err := d.Pool.Query(ctx, query, bid.Name, bid.Description, bid.TenderId, bid.AuthorType, bid.AuthorID)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newBid models.Bid
	err = pgxscan.ScanOne(&newBid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newBid, nil
}

func (d *Database) GetUserBids(ctx context.Context, offset, limit int32, username string) ([]models.Bid, error) {
	const op = "storage.GetUserBids"

	query := `SELECT b.id, b.name, b.description, b.tender_id, b.status, b.decision, b.author_type, b.author_id, b.version, b.created_at, b.updated_at
				FROM bid b
				JOIN employee e ON (b.author_id = e.id)
				WHERE e.username = $1
				ORDER BY b.name
				OFFSET $2
				FETCH NEXT $3 ROWS ONLY;`

	rows, err := d.Pool.Query(ctx, query, username, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var bids []models.Bid
	if err := pgxscan.ScanAll(&bids, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return bids, err
}

func (d *Database) GetBidsForTender(ctx context.Context, tenderId string, offset, limit int32) ([]models.Bid, error) {
	const op = "storage.GetBidsForTender"

	query := `SELECT b.id, b.name, b.description, b.tender_id, b.status, b.decision, b.author_type, b.author_id, b.version, b.created_at, b.updated_at
				FROM bid b
				JOIN employee e ON (b.author_id = e.id)
				WHERE b.tender_id = $1
				ORDER BY b.name
				OFFSET $2
				FETCH NEXT $3 ROWS ONLY;`

	rows, err := d.Pool.Query(ctx, query, tenderId, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var bids []models.Bid
	if err := pgxscan.ScanAll(&bids, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	return bids, err
}

func (d *Database) GetBidStatus(ctx context.Context, bidId, username string) (string, error) {
	const op = "storage.GetBidStatus"

	query := `SELECT b.status
				FROM bid b
				JOIN employee e ON (b.author_id = e.id)
				WHERE b.id = $1 AND e.username = $2;`

	rows, err := d.Pool.Query(ctx, query, bidId, username)
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

func (d *Database) UpdateBidStatus(ctx context.Context, bidId, status, username string) (models.Bid, error) {
	const op = "storage.UpdateBidStatus"

	query := `UPDATE bid b
				SET status = $1
				FROM employee e
				WHERE b.author_id = e.id AND b.id = $2 AND e.username = $3
				RETURNING b.id, b.name, b.description, b.tender_id, b.status, b.decision, b.author_type, b.author_id, b.version, b.created_at, b.updated_at;`

	rows, err := d.Pool.Query(ctx, query, status, bidId, username)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var updatedBid models.Bid
	err = pgxscan.ScanOne(&updatedBid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return updatedBid, nil
}

func (d *Database) EditBid(ctx context.Context, bid *models.Bid, bidId, username string) (models.Bid, error) {
	const op = "storage.EditBid"

	query := `UPDATE bid b
				SET 
					name = COALESCE(NULLIF($1, ''), name), 
					description = COALESCE(NULLIF($2, ''), description) 
				FROM employee e
				WHERE b.author_id = e.id AND b.id = $3 AND e.username = $4
				RETURNING b.id, b.name, b.description, b.tender_id, b.status, b.decision, b.author_type, b.author_id, b.version, b.created_at, b.updated_at;`

	rows, err := d.Pool.Query(ctx, query, bid.Name, bid.Description, bidId, username)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newBid models.Bid
	err = pgxscan.ScanOne(&newBid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newBid, nil
}

func (d *Database) SubmitBidDecision(ctx context.Context, bidId, decision, username string) (models.Bid, error) {
	const op = "storage.SubmitBidDecision"

	query := `UPDATE bid b
				SET decision = $1
				FROM employee e
				WHERE b.author_id = e.id AND b.id = $2 AND e.username = $3
				RETURNING b.id, b.name, b.description, b.tender_id, b.status, b.decision, b.author_type, b.author_id, b.version, b.created_at, b.updated_at;`

	rows, err := d.Pool.Query(ctx, query, decision, bidId, username)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var updatedBid models.Bid
	err = pgxscan.ScanOne(&updatedBid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return updatedBid, nil
}

func (d *Database) SubmitBidFeedback(ctx context.Context, bidId, bidFeedback, username string) (models.Bid, error) {
	const op = "storage.SubmitBidFeedback"

	insertQuery := `INSERT INTO review (bid_id, author_username, description)
				VALUES ($1,	$2, $3);`
	_, err := d.Pool.Exec(ctx, insertQuery, bidId, username, bidFeedback)
	if err != nil {
		return models.Bid{}, err
	}

	query := `SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
				FROM bid
				WHERE id = $1;`

	rows, err := d.Pool.Query(ctx, query, bidId)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var bid models.Bid
	err = pgxscan.ScanOne(&bid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return bid, nil
}

func (d *Database) GetBidReviews(ctx context.Context, authorUsername string, offset, limit int32) ([]models.Review, error) {
	const op = "storage.GetBidReviews"

	query := `SELECT r.id, r.bid_id, r.author_username, r.description, r.created_at
				FROM review r
				WHERE r.author_username = $1
				OFFSET $2
				FETCH NEXT $3 ROWS ONLY;`

	rows, err := d.Pool.Query(ctx, query, authorUsername, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var reviews []models.Review
	err = pgxscan.ScanAll(&reviews, rows)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}

	fmt.Println(reviews)

	return reviews, nil
}

func (d *Database) RollbackBid(ctx context.Context, bidId string, version int32, username string) (models.Bid, error) {
	const op = "storage.RollbackBid"

	query := `SELECT * FROM rollback_bid_version($1::UUID, $2::INT, $3::VARCHAR)`
	rows, err := d.Pool.Query(ctx, query, bidId, version, username)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newBid models.Bid
	err = pgxscan.ScanOne(&newBid, rows)
	if err != nil {
		return models.Bid{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newBid, nil
}