package storage

import (
	"context"
	"zadanie-6105/internal/models"
)

type Storage interface {
	Tender
	Bid
	Helper
}

type Tender interface {
	CreateTender(ctx context.Context, tender *models.Tender) (models.Tender, error)
	GetTenders(ctx context.Context, offset, limit int32, serviceTypes []string) ([]models.Tender, error)
	GetUserTenders(ctx context.Context, offset, limit int32, username string) ([]models.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId, username string) (string, error)
	UpdateTenderStatus(ctx context.Context, tenderId, status, username string) (models.Tender, error)
	EditTender(ctx context.Context, tender *models.Tender, tenderId, username string) (models.Tender, error)
	RollbackTender(ctx context.Context, tenderId string, version int32, username string) (models.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, bid *models.Bid) (models.Bid, error)
	GetUserBids(ctx context.Context, offset, limit int32, username string) ([]models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderId string, offset, limit int32) ([]models.Bid, error)
	GetBidStatus(ctx context.Context, bidId, username string) (string, error)
	UpdateBidStatus(ctx context.Context, bidId, status, username string) (models.Bid, error)
	EditBid(ctx context.Context, bid *models.Bid, bidId, username string) (models.Bid, error)
	SubmitBidDecision(ctx context.Context, bidId, decision, username string) (models.Bid, error)
	SubmitBidFeedback(ctx context.Context, bidId, bidFeedback, username string) (models.Bid, error)
	GetBidReviews(ctx context.Context, authorUsername string, offset, limit int32) ([]models.Review, error)
	RollbackBid(ctx context.Context, bidId string, version int32, username string) (models.Bid, error)
}

type Helper interface {
	IsUserExists(ctx context.Context, username string) error
	IsUserByIdExists(ctx context.Context, id string) error
	IsTenderExists(ctx context.Context, tenderId string) error
	IsBidExists(ctx context.Context, bidId string) error
	ValidateUsersPrivileges(ctx context.Context, tenderId, requestedUser string) error
	ValidateUsersPrivilegesUserId(ctx context.Context, id, tenderId string) error
	ValidateUsersPrivilegesOrgId(ctx context.Context, orgId, username string) error
	ValidateUsersPrivilegesBidId(ctx context.Context, bidId, username string) error
	IsBidVersionExists(ctx context.Context, bidId string, version int32) error
	IsTenderVersionExists(ctx context.Context, tenderId string, version int32) error
}