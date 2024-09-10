package models

import (
	"github.com/google/uuid"
	"time"
)

type TenderStatus string

const (
	Created   TenderStatus = "Created"
	Published TenderStatus = "Published"
	Closed    TenderStatus = "Closed"
)

type ServiceType string

const (
	Construction ServiceType = "Construction"
	Delivery     ServiceType = "Delivery"
	Manufacture  ServiceType = "Manufacture"
)

type Tender struct {
	ID              uuid.UUID    `db:"id" json:"id"`
	Name            string       `db:"name" json:"name"`
	Description     string       `db:"description" json:"description"`
	ServiceType     ServiceType  `db:"service_type" json:"serviceType"`
	Status          TenderStatus `db:"status" json:"status"`
	Version         int          `db:"version" json:"version"`
	OrganizationID  uuid.UUID    `db:"organization_id" json:"organizationId"`
	CreatorUsername string       `db:"creator_username" json:"creatorUsername"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
}

type TenderHistory struct {
	ID              uuid.UUID    `db:"id"`
	TenderID        uuid.UUID    `db:"tender_id"`
	Name            string       `db:"name"`
	Description     string       `db:"description"`
	ServiceType     ServiceType  `db:"service_type"`
	Status          TenderStatus `db:"status"`
	Version         int          `db:"version"`
	OrganizationID  uuid.UUID    `db:"organization_id"`
	CreatorUsername string       `db:"creator_username"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
}