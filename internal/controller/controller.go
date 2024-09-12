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

// CheckServer (GET /api/ping).
func (c *Controller) CheckServer(ctx echo.Context) error {
	ctx.JSON(http.StatusOK, "ok")
	return nil
}

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

// CreateBid (POST /bids/new).
func (c *Controller) CreateBid(ctx echo.Context) error {
	var bid models.Bid

	if err := util.DecodeJSONBody(ctx.Request(), &bid); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequestError
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

// GetUserBids (GET /bids/my).
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

// GetBidsForTender (GET /bids/{tenderId}/list).
func (c *Controller) GetBidsForTender(ctx echo.Context, tenderID TenderId, params GetBidsForTenderParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	bids, err := c.bidService.GetBidsForTender(ctx.Request(), tenderID, offset, limit, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bids)
	return nil
}

// GetBidStatus (GET /bids/{bidId}/status).
func (c *Controller) GetBidStatus(ctx echo.Context, bidID BidId, params GetBidStatusParams) error {
	bids, err := c.bidService.GetBidStatus(ctx.Request(), bidID, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, bids)
	return nil
}

// UpdateBidStatus (PUT /bids/{bidId}/status).
func (c *Controller) UpdateBidStatus(ctx echo.Context, bidID BidId, params UpdateBidStatusParams) error {
	status, err := c.bidService.UpdateBidStatus(ctx.Request(), bidID, string(params.Status), params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

// EditBid (PATCH /bids/{bidId}/edit).
func (c *Controller) EditBid(ctx echo.Context, bidID BidId, params EditBidParams) error {
	var bid models.Bid
	if err := util.DecodeJSONBody(ctx.Request(), &bid); err != nil {
		c.zapLogger.Error(err)
		var mr *util.MalformedRequestError
		if errors.As(err, &mr) {
			ctx.JSON(mr.Status, ErrorResponse{Reason: mr.Msg})
			return err
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
		return err
	}

	newBid, err := c.bidService.EditBid(ctx.Request(), &bid, bidID, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, newBid)
	return nil
}

// SubmitBidDecision (PUT /bids/{bidId}/submit_decision).
func (c *Controller) SubmitBidDecision(ctx echo.Context, bidID BidId, params SubmitBidDecisionParams) error {
	status, err := c.bidService.SubmitBidDecision(ctx.Request(), bidID, string(params.Decision), params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

// SubmitBidFeedback (PUT /bids/{bidId}/feedback).
func (c *Controller) SubmitBidFeedback(ctx echo.Context, bidID BidId, params SubmitBidFeedbackParams) error {
	status, err := c.bidService.SubmitBidFeedback(ctx.Request(), bidID, params.BidFeedback, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, status)
	return nil
}

// RollbackBid (PUT /bids/{bidId}/rollback/{version}).
func (c *Controller) RollbackBid(ctx echo.Context, bidID BidId, version int32, params RollbackBidParams) error {
	reviews, err := c.bidService.RollbackBid(ctx.Request(), bidID, version, params.Username)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, reviews)
	return nil
}

// GetBidReviews (GET /bids/{tenderId}/reviews).
func (c *Controller) GetBidReviews(ctx echo.Context, tenderID TenderId, params GetBidReviewsParams) error {
	var offset, limit int32 = 0, 5
	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	reviews, err := c.bidService.GetBidReviews(ctx.Request(), tenderID, params.AuthorUsername, params.RequesterUsername, offset, limit)
	if err != nil {
		return HandleError(ctx, err)
	}

	ctx.JSON(http.StatusOK, reviews)
	return nil
}