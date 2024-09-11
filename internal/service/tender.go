package service

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/storage"
	"zadanie-6105/internal/util"
)

type TenderService struct {
	storage storage.Storage
}

func NewTenderService(s storage.Storage) *TenderService {
	return &TenderService{storage: s}
}

func (ts *TenderService) CreateTender(r *http.Request, tender *models.Tender) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.IsUserExists(r.Context(), tender.CreatorUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivilegesOrgId(r.Context(), tender.OrganizationID.String(), tender.CreatorUsername)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyTender, err
	}

	return ts.storage.CreateTender(r.Context(), tender)
}

func (ts *TenderService) GetTenders(r *http.Request, offset, limit int32, serviceTypes []string) ([]models.Tender, error) {
	return ts.storage.GetTenders(r.Context(), offset, limit, serviceTypes)
}

func (ts *TenderService) GetUserTenders(r *http.Request, offset, limit int32, username string) ([]models.Tender, error) {
	err := ts.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return nil, err
	}

	return ts.storage.GetUserTenders(r.Context(), offset, limit, username)
}

func (ts *TenderService) GetTenderStatus(r *http.Request, tenderId, username string) (string, error) {
	err := ts.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return "", err
	}

	err = ts.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return "", err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return "", err
	}

	return ts.storage.GetTenderStatus(r.Context(), tenderId, username)
}

func (ts *TenderService) UpdateTenderStatus(r *http.Request, tenderId, status, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyTender, err
	}

	err = ts.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyTender, err
	}

	return ts.storage.UpdateTenderStatus(r.Context(), tenderId, status, username)
}

func (ts *TenderService) EditTender(r *http.Request, tender *models.Tender, tenderId, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.NotFound}
		}
		return emptyTender, err
	}

	err = ts.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyTender, err
	}

	return ts.storage.EditTender(r.Context(), tender, tenderId, username)
}

func (ts *TenderService) RollbackTender(r *http.Request, tenderId string, version int32, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.IsTenderExists(r.Context(), tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.VersionNotFound}
		}
		return emptyTender, err
	}

	err = ts.storage.IsUserExists(r.Context(), username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusUnauthorized, Msg: util.Unathorized}
		}
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderId, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusForbidden, Msg: util.Forbidden}
		}
		return emptyTender, err
	}

	err = ts.storage.IsTenderVersionExists(r.Context(), tenderId, version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return emptyTender, util.MyErrorResponse{Status: http.StatusNotFound, Msg: util.VersionNotFound}
		}
		return emptyTender, err
	}

	return ts.storage.RollbackTender(r.Context(), tenderId, version, username)
}