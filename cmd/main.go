package main

import (
	"context"
	"zadanie-6105/internal/api"
	"zadanie-6105/internal/controller"
	"zadanie-6105/internal/service"
	"zadanie-6105/internal/storage/postgres"
	"zadanie-6105/internal/util"
)

// TODO: Решение должно быть загружено в репозиторий в ветку main.
// Необходимо клонировать репозиторий к себе, выполнить задание и
// запушить его обратно.
func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()
	storage := postgres.NewPostgresRepository(ctx, util.NewDBConfig(), zapLogger)
	tenderService := service.NewTenderService(storage)
	bidService := service.NewBidService(storage)
	ctrl := controller.NewController(zapLogger, tenderService, bidService)

	app := api.NewAPI(ctrl, zapLogger, util.NewServerConfig())

	app.Run(ctx)
}