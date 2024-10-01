package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/config"
	"github.com/aridae/go-metrics-store/internal/server/mw"
	"github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
	"log"
)

func main() {
	ctx := context.Background()
	cnf := config.Obtain()

	memStore := tsstorage.New()

	metricsRepo := scalarmetric.NewRepository(memStore)

	useCaseController := usecases.NewController(metricsRepo)

	httpRouter := handlers.NewRouter(useCaseController)

	httpServer := http.NewServer(cnf.GetAddress(), httpRouter,
		mw.LoggingMiddleware,
		mw.GzipDecompressRequestMiddleware,
		mw.GzipCompressResponseMiddleware,
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
