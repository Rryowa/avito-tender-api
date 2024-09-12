package storage

import (
	"context"
	"zadanie-6105/internal/models"
)

type Storage interface {
	Tender
	Bid
	Checker
	Validator
}

type Tender interface {
	CreateTender(ctx context.Context, tender *models.Tender) (models.Tender, error)
	GetTenders(ctx context.Context, offset, limit int32, serviceTypes []string) ([]models.Tender, error)
	GetUserTenders(ctx context.Context, offset, limit int32, username string) ([]models.Tender, error)
	GetTenderStatus(ctx context.Context, tenderID, username string) (string, error)
	UpdateTenderStatus(ctx context.Context, tenderID, status, username string) (models.Tender, error)
	EditTender(ctx context.Context, tender *models.Tender, tenderID, username string) (models.Tender, error)
	RollbackTender(ctx context.Context, tenderID string, version int32, username string) (models.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, bid *models.Bid) (models.Bid, error)
	GetUserBids(ctx context.Context, offset, limit int32, username string) ([]models.Bid, error)
	GetBidsForTender(ctx context.Context, tenderID string, offset, limit int32, status ...string) ([]models.Bid, error)
	GetBidStatus(ctx context.Context, bidID, username string) (string, error)
	UpdateBidStatus(ctx context.Context, bidID, status, username string) (models.Bid, error)
	EditBid(ctx context.Context, bid *models.Bid, bidID, username string) (models.Bid, error)
	SubmitBidDecision(ctx context.Context, bidID, decision, username string) (models.Bid, error)
	SubmitBidFeedback(ctx context.Context, bidID, bidFeedback, username string) (models.Bid, error)
	GetBidReviews(ctx context.Context, authorUsername string, offset, limit int32) ([]models.Review, error)
	RollbackBid(ctx context.Context, bidID string, version int32, username string) (models.Bid, error)
}

type Checker interface {
	CheckUserExists(ctx context.Context, username string) error
	CheckUserByIDExists(ctx context.Context, id string) error
	CheckTenderExists(ctx context.Context, tenderID string) error
	CheckBidExists(ctx context.Context, bidID string) error
	CheckUserBidAuthor(ctx context.Context, bidID, requestedUser string) error
	CheckBidVersionExists(ctx context.Context, bidID string, version int32) error
	CheckTenderVersionExists(ctx context.Context, tenderID string, version int32) error
}

type Validator interface {
	ValidateUserResponsible(ctx context.Context, tenderID, requestedUser string) error
	ValidateUserResponsibleUserID(ctx context.Context, userID, tenderID string) error
	ValidateUserResponsibleOrgID(ctx context.Context, orgID, username string) error
	ValidateUserResponsibleBidID(ctx context.Context, bidID, username string) error
}