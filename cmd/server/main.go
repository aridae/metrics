package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/config"
	"github.com/aridae/go-metrics-store/internal/server/logger"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/mw"
	"github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

		select {
		case <-signalCh:
			logger.Obtain().Infof("Got signal %s, shutting down...", "FUCK YOU")
			cancel()
		}
	}()

	cnf := config.Obtain()

	memStore := mustInitMemStore(ctx, cnf)

	metricsRepo := scalarmetric.NewRepository(memStore)

	useCaseController := usecases.NewController(metricsRepo)

	httpRouter := handlers.NewRouter(useCaseController)

	httpServer := http.NewServer(cnf.Address, httpRouter,
		mw.LoggingMiddleware,
		mw.GzipDecompressRequestMiddleware,
		mw.GzipCompressResponseMiddleware,
	)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func mustInitMemStore(ctx context.Context, cnf *config.Config) *tsstorage.MemTimeseriesStorage {
	memStore := tsstorage.New()

	err := memStore.InitBackup(ctx, cnf.FileStoragePath, cnf.StoreInterval, map[string]any{
		"ScalarMetric":       models.ScalarMetric{},
		"Int64MetricValue":   models.NewInt64MetricValue(0),
		"Float64MetricValue": models.NewFloat64MetricValue(0),
	})
	if err != nil {
		log.Fatalf("failed to init mem store backup: %v", err)
	}

	if cnf.Restore {
		err = memStore.LoadFromBackup()
		if err != nil {
			log.Fatalf("failed to load mem store from backup: %v", err)
		}
	}

	return memStore
}
