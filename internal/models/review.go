package models

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	BidID          uuid.UUID  `db:"bid_id" json:"bidId"`
	AuthorUsername string     `db:"author_username" json:"authorUsername,omitempty"`
	Description    string     `db:"description" json:"description"`
	CreatedAt      *time.Time `db:"created_at"`
}