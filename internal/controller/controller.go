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
	bidService    *service.BidService
}

func NewController(l *zap.SugaredLogger, ts *service.TenderService, bs *service.BidService) *Controller {
	return &Controller{
		zapLogger:     l,
		tenderService: ts,
		bidService:    bs,
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

func (c *Controller) CreateBid(ctx echo.Context) error {
	var bid models.Bid

	if err := util.DecodeJSONBody(ctx.Request(), &bid); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newBid, err := c.bidService.CreateBid(ctx.Request(), &bid)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, newBid)
	return nil
}

func (c *Controller) GetUserBids(ctx echo.Context, params GetUserBidsParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	bids, err := c.bidService.GetUserBids(ctx.Request(), offset, limit, *params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bids)
	return nil
}

func (c *Controller) GetBidsForTender(ctx echo.Context, tenderId TenderId, params GetBidsForTenderParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	bids, err := c.bidService.GetBidsForTender(ctx.Request(), tenderId, offset, limit, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bids)
	return nil
}

func (c *Controller) GetBidStatus(ctx echo.Context, bidId BidId, params GetBidStatusParams) error {
	bids, err := c.bidService.GetBidStatus(ctx.Request(), bidId, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bids)
	return nil
}

func (c *Controller) UpdateBidStatus(ctx echo.Context, bidId BidId, params UpdateBidStatusParams) error {
	status, err := c.bidService.UpdateBidStatus(ctx.Request(), bidId, string(params.Status), params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

func (c *Controller) EditBid(ctx echo.Context, bidId BidId, params EditBidParams) error {
	var bid models.Bid
	if err := util.DecodeJSONBody(ctx.Request(), &bid); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequest
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newBid, err := c.bidService.EditBid(ctx.Request(), &bid, bidId, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, newBid)
	return nil
}

func (c *Controller) SubmitBidDecision(ctx echo.Context, bidId BidId, params SubmitBidDecisionParams) error {
	status, err := c.bidService.SubmitBidDecision(ctx.Request(), bidId, string(params.Decision), params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

func (c *Controller) SubmitBidFeedback(ctx echo.Context, bidId BidId, params SubmitBidFeedbackParams) error {
	status, err := c.bidService.SubmitBidFeedback(ctx.Request(), bidId, params.BidFeedback, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

func (c *Controller) GetBidReviews(ctx echo.Context, tenderId TenderId, params GetBidReviewsParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	reviews, err := c.bidService.GetBidReviews(ctx.Request(), tenderId, params.AuthorUsername, params.RequesterUsername, offset, limit)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, reviews)
	return nil
}

func (c *Controller) RollbackBid(ctx echo.Context, bidId BidId, version int32, params RollbackBidParams) error {
	reviews, err := c.bidService.RollbackBid(ctx.Request(), bidId, version, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, reviews)
	return nil
}

func HandleError(ctx echo.Context, err error) error {
	var customErr util.MyErrorResponse
	if errors.As(err, &customErr) {
		ctx.JSON(customErr.Status, ErrorResponse{Reason: customErr.Msg})
		return err
	}

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
	return err
}