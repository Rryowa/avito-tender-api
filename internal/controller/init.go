//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=openapi/cfg.yaml openapi/yaml

package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
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

// InternalError to return Internal Server Error.
func InternalError(ctx echo.Context, err error) error {
	var customErr util.MyResponseError
	if errors.As(err, &customErr) {
		ctx.JSON(customErr.Status, ErrorResponse{Reason: customErr.Msg})
		return err
	}

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
	return err
}