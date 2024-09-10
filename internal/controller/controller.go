//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml openapi.yaml
package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/service"
	"zadanie-6105/internal/util"
)

type Controller struct {
	zapLogger     *zap.SugaredLogger
	tenderService *service.TenderService
}

func NewController(l *zap.SugaredLogger, ts *service.TenderService) *Controller {
	return &Controller{
		zapLogger:     l,
		tenderService: ts,
	}
}

// CheckServer (GET /api/ping)
func (c *Controller) CheckServer(ctx echo.Context) error {
	ctx.JSON(http.StatusOK, "ok")
	return nil
}

// GetTenders (GET /api/tenders/)
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

func (c *Controller) CreateTender(ctx echo.Context) error {
	var tender models.Tender

	if err := util.DecodeJSONBody(ctx.Request(), &tender); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
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

func (c *Controller) GetTenderStatus(ctx echo.Context, tenderId TenderId, params GetTenderStatusParams) error {
	tender, err := c.tenderService.GetTenderStatus(ctx.Request(), tenderId, *params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, tender)
	return nil
}

func (c *Controller) UpdateTenderStatus(ctx echo.Context, tenderId TenderId, params UpdateTenderStatusParams) error {
	status, err := c.tenderService.UpdateTenderStatus(ctx.Request(), tenderId, string(params.Status), params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

func (c *Controller) EditTender(ctx echo.Context, tenderId TenderId, params EditTenderParams) error {
	var tender models.Tender
	if err := util.DecodeJSONBody(ctx.Request(), &tender); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newTender, err := c.tenderService.EditTender(ctx.Request(), &tender, tenderId, params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, newTender)
	return nil
}

func (c *Controller) RollbackTender(ctx echo.Context, tenderId TenderId, version int32, params RollbackTenderParams) error {
	newTender, err := c.tenderService.RollbackTender(ctx.Request(), tenderId, version, params.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	ctx.JSON(http.StatusOK, newTender)
	return nil
}

func (c *Controller) GetUserBids(ctx echo.Context, params GetUserBidsParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) CreateBid(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) EditBid(ctx echo.Context, bidId BidId, params EditBidParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) SubmitBidFeedback(ctx echo.Context, bidId BidId, params SubmitBidFeedbackParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) RollbackBid(ctx echo.Context, bidId BidId, version int32, params RollbackBidParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) GetBidStatus(ctx echo.Context, bidId BidId, params GetBidStatusParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) UpdateBidStatus(ctx echo.Context, bidId BidId, params UpdateBidStatusParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) SubmitBidDecision(ctx echo.Context, bidId BidId, params SubmitBidDecisionParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) GetBidsForTender(ctx echo.Context, tenderId TenderId, params GetBidsForTenderParams) error {
	//TODO implement me
	panic("implement me")
}

func (c *Controller) GetBidReviews(ctx echo.Context, tenderId TenderId, params GetBidReviewsParams) error {
	//TODO implement me
	panic("implement me")
}