package service

import (
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/storage"
)

type TenderService struct {
	storage storage.Storage
}

func NewTenderService(s storage.Storage) *TenderService {
	return &TenderService{storage: s}
}

func (ts *TenderService) CreateTender(r *http.Request, tender *models.Tender) (models.Tender, error) {
	var emptyTender models.Tender
	if err := ts.storage.CheckUserExists(r.Context(), tender.CreatorUsername); err != nil {
		return emptyTender, err
	}

	if err := ts.storage.ValidateUsersPrivilegesOrgID(r.Context(), tender.OrganizationID.String(), tender.CreatorUsername); err != nil {
		return emptyTender, err
	}

	return ts.storage.CreateTender(r.Context(), tender)
}

func (ts *TenderService) GetTenders(r *http.Request, offset, limit int32, serviceTypes []string) ([]models.Tender, error) {
	return ts.storage.GetTenders(r.Context(), offset, limit, serviceTypes)
}

func (ts *TenderService) GetUserTenders(r *http.Request, offset, limit int32, username string) ([]models.Tender, error) {
	err := ts.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return nil, err
	}

	return ts.storage.GetUserTenders(r.Context(), offset, limit, username)
}

func (ts *TenderService) GetTenderStatus(r *http.Request, tenderID, username string) (string, error) {
	err := ts.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return "", err
	}

	err = ts.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return "", err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderID, username)
	if err != nil {
		return "", err
	}

	return ts.storage.GetTenderStatus(r.Context(), tenderID, username)
}

func (ts *TenderService) UpdateTenderStatus(r *http.Request, tenderID, status, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderID, username)
	if err != nil {
		return emptyTender, err
	}

	return ts.storage.UpdateTenderStatus(r.Context(), tenderID, status, username)
}

func (ts *TenderService) EditTender(r *http.Request, tender *models.Tender, tenderID, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderID, username)
	if err != nil {
		return emptyTender, err
	}

	return ts.storage.EditTender(r.Context(), tender, tenderID, username)
}

func (ts *TenderService) RollbackTender(r *http.Request, tenderID string, version int32, username string) (models.Tender, error) {
	var emptyTender models.Tender
	err := ts.storage.CheckTenderExists(r.Context(), tenderID)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.CheckUserExists(r.Context(), username)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.ValidateUsersPrivileges(r.Context(), tenderID, username)
	if err != nil {
		return emptyTender, err
	}

	err = ts.storage.CheckTenderVersionExists(r.Context(), tenderID, version)
	if err != nil {
		return emptyTender, err
	}

	return ts.storage.RollbackTender(r.Context(), tenderID, version, username)
}