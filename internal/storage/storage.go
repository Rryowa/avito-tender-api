package storage

import (
	"context"
	"zadanie-6105/internal/models"
)

type Storage interface {
	CreateTender(ctx context.Context, tender *models.Tender) (models.Tender, error)
	GetTenders(ctx context.Context, offset, limit int32, serviceTypes []string) ([]models.Tender, error)
	GetUserTenders(ctx context.Context, offset, limit int32, username string) ([]models.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId, username string) (string, error)
	UpdateTenderStatus(ctx context.Context, tenderId, status, username string) (models.Tender, error)
	EditTender(ctx context.Context, tender *models.Tender, tenderId, username string) (models.Tender, error)
	RollbackTender(ctx context.Context, tenderId string, version int32, username string) (models.Tender, error)
}