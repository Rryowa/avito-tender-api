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

var ServiceTypeMap = map[ServiceType]struct{}{
	Construction: {},
	Delivery:     {},
	Manufacture:  {},
}

type Tender struct {
	ID              uuid.UUID    `db:"id" json:"id"`
	Name            string       `db:"name" json:"name"`
	Description     string       `db:"description" json:"description,omitempty"`
	ServiceType     ServiceType  `db:"service_type" json:"serviceType"`
	Status          TenderStatus `db:"status" json:"status"`
	Version         int          `db:"version" json:"version"`
	OrganizationID  uuid.UUID    `db:"organization_id" json:"organizationId,omitempty"`
	CreatorUsername string       `db:"creator_username" json:"creatorUsername"`
	CreatedAt       *time.Time   `db:"created_at"`
	UpdatedAt       *time.Time   `db:"updated_at,omitempty"`
}

type TenderHistory struct {
	ID              uuid.UUID    `db:"id"`
	TenderID        uuid.UUID    `db:"tender_id"`
	Name            string       `db:"name" json:"name"`
	Description     string       `db:"description" json:"description,omitempty"`
	ServiceType     ServiceType  `db:"service_type" json:"serviceType"`
	Status          TenderStatus `db:"status" json:"status"`
	Version         int          `db:"version" json:"version"`
	OrganizationID  uuid.UUID    `db:"organization_id" json:"organizationId,omitempty"`
	CreatorUsername string       `db:"creator_username" json:"creatorUsername"`
	CreatedAt       *time.Time   `db:"created_at"`
	UpdatedAt       *time.Time   `db:"updated_at,omitempty"`
}