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
	return ts.storage.CreateTender(r.Context(), tender)
}

func (ts *TenderService) GetTenders(r *http.Request, offset, limit int32, serviceTypes []string) ([]models.Tender, error) {
	return ts.storage.GetTenders(r.Context(), offset, limit, serviceTypes)
}

func (ts *TenderService) GetUserTenders(r *http.Request, offset, limit int32, username string) ([]models.Tender, error) {
	return ts.storage.GetUserTenders(r.Context(), offset, limit, username)
}

func (ts *TenderService) GetTenderStatus(r *http.Request, tenderId, username string) (string, error) {
	return ts.storage.GetTenderStatus(r.Context(), tenderId, username)
}

func (ts *TenderService) UpdateTenderStatus(r *http.Request, tenderId, status, username string) (models.Tender, error) {
	return ts.storage.UpdateTenderStatus(r.Context(), tenderId, status, username)
}

func (ts *TenderService) EditTender(r *http.Request, tender *models.Tender, tenderId, username string) (models.Tender, error) {
	return ts.storage.EditTender(r.Context(), tender, tenderId, username)
}

func (ts *TenderService) RollbackTender(r *http.Request, tenderId string, version int32, username string) (models.Tender, error) {
	return ts.storage.RollbackTender(r.Context(), tenderId, version, username)
}