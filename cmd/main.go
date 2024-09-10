package main

import (
	"context"
	"zadanie-6105/internal/api"
	"zadanie-6105/internal/controller"
	"zadanie-6105/internal/service"
	"zadanie-6105/internal/storage/postgres"
	"zadanie-6105/internal/util"
)

// eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MjU5MDA3NzAsImV4cCI6MTc1NzQzNjc3MCwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJSb2xlIjpbIk1hbmFnZXIiLCJQcm9qZWN0IEFkbWluaXN0cmF0b3IiXX0.KWq8eGJpcbZt0l5k2VGXvG59VseYzIRVqkg5I6RR5Uc
func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()
	storage := postgres.NewPostgresRepository(ctx, util.NewDbConfig(), zapLogger)
	tenderService := service.NewTenderService(storage)
	ctrl := controller.NewController(zapLogger, tenderService)

	app := api.NewAPI(ctrl, zapLogger, util.NewServerConfig())

	app.Run(ctx)
}