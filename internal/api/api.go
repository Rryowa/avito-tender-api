package api

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"zadanie-6105/internal/handlers"
	"zadanie-6105/internal/models/config"
)

const (
	shutdownTimeout = 5 * time.Second
)

type API struct {
	server        *http.Server
	handler       *handlers.MyHandler
	zapLogger     *zap.SugaredLogger
	telemetryAddr string
}

func NewAPI(h *handlers.MyHandler, l *zap.SugaredLogger, sc *config.ServerConfig) *API {
	return &API{
		server: &http.Server{
			Addr:         sc.ServerAddr,
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		handler:       h,
		zapLogger:     l,
		telemetryAddr: sc.TelemetryAddr,
	}
}

func (a *API) Run(ctxBackground context.Context) {
	ctx, stop := signal.NotifyContext(ctxBackground, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := mux.NewRouter()
	swagger, err := handlers.GetSwagger()
	if err != nil {
		a.zapLogger.Fatalf("Failed to load OpenAPI specification: %v", err)
	}
	swagger.Servers = nil

	// Use the middleware to validate the incoming requests
	router.Use(nethttpmiddleware.OapiRequestValidator(swagger))

	//the generated code sets up the routing to match the OpenAPI spec and
	//delegates request handling to my GetApiPing method.
	//The generated ServerInterfaceWrapper wraps my handler and calls my GetApiPing
	//method when the corresponding route (/api/ping) is accessed.
	handlers.HandlerFromMux(a.handler, router)

	a.server.Handler = router

	a.ListenGracefulShutdown(ctx)
}

func (a *API) ListenGracefulShutdown(ctx context.Context) {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.zapLogger.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()
	a.zapLogger.Infof("Listening on: %v\n", a.server.Addr)

	<-ctx.Done()
	a.zapLogger.Info("Shutting down server...\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.zapLogger.Errorf("shutdown: %v", err)
	}

	longShutdown := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- struct{}{}
	}()

	select {
	case <-shutdownCtx.Done():
		a.zapLogger.Errorf("server shutdown: %v", ctx.Err())
	case <-longShutdown:
		a.zapLogger.Infof("finished")
	}
}