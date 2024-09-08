//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config cfg.yaml ../../openapi/api.yaml
package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"zadanie-6105/internal/util"
)

type MyHandler struct {
	zapLogger *zap.SugaredLogger
}

func NewMyHandler(l *zap.SugaredLogger) *MyHandler {
	return &MyHandler{
		zapLogger: l,
	}
}

// GetPing (GET /api/ping)
func (h *MyHandler) GetPing(w http.ResponseWriter, r *http.Request) {
	util.WriteJSON(w, map[string]string{"message": "pong"})
}