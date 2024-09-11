package service

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/storage"
	"zadanie-6105/internal/util"
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
	err := bs.storage.IsTenderExists(r.Context(), bid.TenderId.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserByIdExists(r.Context(), bid.AuthorID.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesUserId(r.Context(), bid.AuthorID.String(), bid.TenderId.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden + err.Error()}
		}
		return emptyBid, err
	}

	return bs.storage.CreateBid(r.Context(), bid)
}

func (bs *BidService) GetUserBids(r *http.Request, offset, limit int32, username string) ([]models.Bid, error) {
	err := bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return nil, err
	}

	return bs.storage.GetUserBids(r.Context(), offset, limit, username)
}

func (bs *BidService) GetBidsForTender(r *http.Request, tenderId string, offset, limit int32, username string) ([]models.Bid, error) {
	err := bs.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return nil, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return nil, err
	}

	err = bs.storage.ValidateUsersPrivileges(r.Context(), tenderId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return nil, err
	}

	return bs.storage.GetBidsForTender(r.Context(), tenderId, offset, limit)
}

func (bs *BidService) GetBidStatus(r *http.Request, bidId, username string) (string, error) {
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return "", err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return "", err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return "", err
	}

	return bs.storage.GetBidStatus(r.Context(), bidId, username)
}

func (bs *BidService) UpdateBidStatus(r *http.Request, bidId, status, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyBid, err
	}

	return bs.storage.UpdateBidStatus(r.Context(), bidId, status, username)
}

func (bs *BidService) EditBid(r *http.Request, bid *models.Bid, bidId, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyBid, err
	}

	return bs.storage.EditBid(r.Context(), bid, bidId, username)
}

func (bs *BidService) SubmitBidDecision(r *http.Request, bidId, decision, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyBid, err
	}

	return bs.storage.SubmitBidDecision(r.Context(), bidId, decision, username)
}

func (bs *BidService) SubmitBidFeedback(r *http.Request, bidId, bidFeedback, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyBid, err
	}

	return bs.storage.SubmitBidFeedback(r.Context(), bidId, bidFeedback, username)
}

func (bs *BidService) GetBidReviews(r *http.Request, tenderId, authorUsername, requesterUsername string, offset, limit int32) ([]models.Review, error) {
	err := bs.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return nil, err
	}

	err = bs.storage.IsUserExists(r.Context(), requesterUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return nil, err
	}

	err = bs.storage.ValidateUsersPrivileges(r.Context(), tenderId, requesterUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return nil, err
	}

	return bs.storage.GetBidReviews(r.Context(), authorUsername, offset, limit)
}

func (bs *BidService) RollbackBid(r *http.Request, bidId string, version int32, username string) (models.Bid, error) {
	var emptyBid models.Bid
	err := bs.storage.IsBidExists(r.Context(), bidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyBid, err
	}

	err = bs.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyBid, err
	}

	err = bs.storage.ValidateUsersPrivilegesBidId(r.Context(), bidId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyBid, err
	}

	err = bs.storage.IsBidVersionExists(r.Context(), bidId, version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyBid, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.VersionNotFound}
		}
		return emptyBid, err
	}

	return bs.storage.RollbackBid(r.Context(), bidId, version, username)
}