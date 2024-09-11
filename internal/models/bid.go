package models

import (
	"github.com/google/uuid"
	"time"
)

type BidStatus string

type BidDecision string

type AuthorType string

const (
	OrganizationAuthorType = "Organization"
	UserAuthorType         = "User"
)

type Bid struct {
	ID             uuid.UUID     `db:"id" json:"id"`
	Name           string        `db:"name" json:"name"`
	Description    string        `db:"description" json:"description"`
	Feedback       *string       `db:"feedback" json:"feedback,omitempty"`
	Status         BidStatus     `db:"status" json:"status"`
	TenderId       uuid.UUID     `db:"tender_id" json:"tenderId"`
	OrganizationID uuid.NullUUID `db:"organization_id" json:"organizationId,omitempty"`
	Decision       BidDecision   `db:"decision" json:"decision,omitempty"`
	AuthorID       uuid.UUID     `db:"author_id" json:"authorId"`
	AuthorUsername string        `db:"author_username" json:"authorUsername,omitempty"`
	AuthorType     AuthorType    `db:"author_type" json:"authorType"`
	Version        int           `db:"version" json:"version"`
	CreatedAt      *time.Time    `db:"created_at"`
	UpdatedAt      *time.Time    `db:"updated_at"`
}

type BidHistory struct {
	ID             uuid.UUID     `db:"id"`
	BidID          uuid.UUID     `db:"bid_id"`
	Name           string        `db:"name"`
	Description    string        `db:"description"`
	Feedback       *string       `db:"feedback"`
	Status         BidStatus     `db:"status"`
	TenderId       uuid.UUID     `db:"tender_id"`
	OrganizationID uuid.NullUUID `db:"organization_id"`
	Decision       BidDecision   `db:"decision"`
	AuthorID       uuid.UUID     `db:"author_id"`
	AuthorUsername string        `db:"author_username" json:"authorUsername"`
	AuthorType     AuthorType    `db:"author_type"`
	Version        int           `db:"version"`
	CreatedAt      *time.Time    `db:"created_at"`
	UpdatedAt      *time.Time    `db:"updated_at"`
}