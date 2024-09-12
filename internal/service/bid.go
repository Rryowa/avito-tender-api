package service

import (
	"errors"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/storage"
)

type BidService struct {
	storage storage.Storage
}

func NewBidService(s storage.Storage) *BidService {
	return &BidService{storage: s}
}

func (bs *BidService) CreateBid(r *http.Request, bid *models.Bid) (models.Bid, error) {
	if bid.AuthorType != models.OrganizationAuthorType {
		return models.Bid{}, errors.New("предложения можно создавать только от имени организации")
	}

	var emptyBid models.Bid
	err := bs.storage.CheckTenderExists(r.Context(), bid.TenderID.String())
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserByIDExists(r.Context(), bid.AuthorID.String())
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesUserID(r.Context(), bid.AuthorID.String(), bid.TenderID.String())
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.CreateBid(r.Context(), bid)
}

func (bs *BidService) GetUserBids(r *http.Request, offset, limit int32, username string) ([]models.Bid, error) {
	err := bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return nil, err
	}

	return bs.storage.GetUserBids(r.Context(), offset, limit, username)
}

func (bs *BidService) GetBidsForTender(r *http.Request, tenderID string, offset, limit int32, username string) ([]models.Bid, error) {
	err := bs.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return nil, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return nil, err
	}

	err = bs.storage.ValidateUsersPrivileges(r.Context(), tenderID, username)
	if err != nil {
		return nil, err
	}

	return bs.storage.GetBidsForTender(r.Context(), tenderID, offset, limit)
}

func (bs *BidService) GetBidStatus(r *http.Request, bidID, username string) (string, error) {
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return "", err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return "", err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return "", err
	}

	return bs.storage.GetBidStatus(r.Context(), bidID, username)
}

func (bs *BidService) UpdateBidStatus(r *http.Request, bidID, status, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.UpdateBidStatus(r.Context(), bidID, status, username)
}

func (bs *BidService) EditBid(r *http.Request, bid *models.Bid, bidID, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.EditBid(r.Context(), bid, bidID, username)
}

func (bs *BidService) SubmitBidDecision(r *http.Request, bidID, decision, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.SubmitBidDecision(r.Context(), bidID, decision, username)
}

func (bs *BidService) SubmitBidFeedback(r *http.Request, bidID, bidFeedback, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.SubmitBidFeedback(r.Context(), bidID, bidFeedback, username)
}

func (bs *BidService) GetBidReviews(r *http.Request, tenderID, authorUsername, requesterUsername string, offset, limit int32) ([]models.Review, error) {
	err := bs.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return nil, err
	}

	err = bs.storage.CheckUserExists(r.Context(), requesterUsername)
	if err != nil {
		return nil, err
	}

	err = bs.storage.ValidateUsersPrivileges(r.Context(), tenderID, requesterUsername)
	if err != nil {
		return nil, err
	}

	return bs.storage.GetBidReviews(r.Context(), authorUsername, offset, limit)
}

func (bs *BidService) RollbackBid(r *http.Request, bidID string, version int32, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.CheckBidExists(r.Context(), bidID)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidID(r.Context(), bidID, username)
	if err != nil {
		return emptyBid, err
	}

	err = bs.storage.CheckBidVersionExists(r.Context(), bidID, version)
	if err != nil {
		return emptyBid, err
	}

	return bs.storage.RollbackBid(r.Context(), bidID, version, username)
}