package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"zadanie-6105/internal/util"
)

func HandleError(ctx echo.Context, err error) error {
	var customErr util.MyResponseError
	if errors.As(err, &customErr) {
		ctx.JSON(customErr.Status, ErrorResponse{Reason: customErr.Msg})
		return err
	}

	ctx.JSON(http.StatusInternalServerError, ErrorResponse{Reason: err.Error()})
	return err
}