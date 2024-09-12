//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml openapi.yaml
package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/util"
)

// GetTenders (GET /api/tenders/).
func (c *Controller) GetTenders(ctx echo.Context, params GetTendersParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	var serviceTypes []string
	if params.ServiceType != nil {
		for _, st := range *params.ServiceType {
			serviceTypes = append(serviceTypes, string(st))
		}
	}

	tenders, err := c.tenderService.GetTenders(ctx.Request(), offset, limit, serviceTypes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, tenders)
	return nil
}

// CreateTender (POST /tenders/new).
func (c *Controller) CreateTender(ctx echo.Context) error {
	var tender models.Tender

	if err := util.DecodeJSONBody(ctx.Request(), &tender); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequestError
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newTender, err := c.tenderService.CreateTender(ctx.Request(), &tender)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, newTender)
	return nil
}

// GetUserTenders (GET /tenders/my).
func (c *Controller) GetUserTenders(ctx echo.Context, params GetUserTendersParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	tenders, err := c.tenderService.GetUserTenders(ctx.Request(), offset, limit, *params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, tenders)
	return nil
}

// GetTenderStatus (GET /tender/{tenderId}/status).
func (c *Controller) GetTenderStatus(ctx echo.Context, tenderID TenderId, params GetTenderStatusParams) error {
	tender, err := c.tenderService.GetTenderStatus(ctx.Request(), tenderID, *params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, tender)
	return nil
}

// UpdateTenderStatus (PUT /tender/{tenderId}/status).
func (c *Controller) UpdateTenderStatus(ctx echo.Context, tenderID TenderId, params UpdateTenderStatusParams) error {
	status, err := c.tenderService.UpdateTenderStatus(ctx.Request(), tenderID, string(params.Status), params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

// EditTender (PATCH /tenders/{tenderId}/edit).
func (c *Controller) EditTender(ctx echo.Context, tenderID TenderId, params EditTenderParams) error {
	var tender models.Tender
	if err := util.DecodeJSONBody(ctx.Request(), &tender); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequestError
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newTender, err := c.tenderService.EditTender(ctx.Request(), &tender, tenderID, params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, newTender)
	return nil
}

// RollbackTender (PUT /tenders/{tenderId}/rollback/{version}).
func (c *Controller) RollbackTender(ctx echo.Context, tenderID TenderId, version int32, params RollbackTenderParams) error {
	newTender, err := c.tenderService.RollbackTender(ctx.Request(), tenderID, version, params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, newTender)
	return nil
}