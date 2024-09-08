package main

import (
	"context"
	"zadanie-6105/internal/api"
	"zadanie-6105/internal/handlers"
	"zadanie-6105/internal/util"
)

func main() {
	ctx := context.Background()
	zapLogger := util.NewZapLogger()
	handler := handlers.NewMyHandler(zapLogger)

	app := api.NewAPI(handler, zapLogger, util.NewServerConfig())

	app.Run(ctx)
}