package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/config"
	"github.com/aridae/go-metrics-store/internal/server/logger"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/mw"
	scalarmetricinmem "github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric/inmemory"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	"github.com/aridae/go-metrics-store/pkg/postgres"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

		<-signalCh

		logger.Obtain().Info("Got signal, shutting down...")

		// If you fail to cancel the context, the goroutine that WithCancel or WithTimeout created
		// will be retained in memory indefinitely (until the program shuts down), causing a memory leak.
		cancel()
	}()

	cnf := config.Obtain()

	var pgClient *postgres.Client
	if cnf.DatabaseDsn != "" {
		pgClient = mustInitPostgresClient(ctx, cnf)
	}

	memStore := mustInitMemStore(ctx, cnf)

	metricsRepo := scalarmetricinmem.NewRepositoryImplementation(memStore)

	useCaseController := usecases.NewController(metricsRepo, pgClient)

	httpRouter := handlers.NewRouter(useCaseController)

	httpServer := http.NewServer(cnf.Address, httpRouter,
		mw.LoggingMiddleware,
		mw.GzipDecompressRequestMiddleware,
		mw.GzipCompressResponseMiddleware,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Obtain().Fatalf("failed to start server: %v", err)
	}
}

func mustInitMemStore(ctx context.Context, cnf *config.Config) *tsstorage.MemTimeseriesStorage {
	memStore := tsstorage.New()

	// NOTE: tsstorage.MemTimeseriesStorage работает с интерфейсом TimeseriesValue
	// и не знает о том, какие модельки передаются под капотом. Но из-за этого,
	// при бэкапе в файл, и последующем чтении из файла MemTimeseriesStorage не может знать,
	// в какую структуру/структурки десереализовать содержимое файла.
	// Чтобы не писать свои маршаллеры/анмаршраллеры на рефлексии,
	// я регистрирую типы для использования в gob.Encoder/Decoder.
	// Но это делает стор зависимым от гошных моделек, и мне от этого грустно.
	err := memStore.InitBackup(ctx, cnf.FileStoragePath, cnf.StoreInterval, map[string]any{
		"ScalarMetric":       models.ScalarMetric{},
		"Int64MetricValue":   models.NewInt64MetricValue(0),
		"Float64MetricValue": models.NewFloat64MetricValue(0),
	})
	if err != nil {
		logger.Obtain().Fatalf("failed to init mem store backup: %v", err)
	}

	if cnf.Restore {
		err = memStore.LoadFromBackup()
		if err != nil {
			logger.Obtain().Fatalf("failed to load mem store from backup: %v", err)
		}
	}

	return memStore
}

func mustInitPostgresClient(ctx context.Context, cnf *config.Config) *postgres.Client {
	client, err := postgres.NewClient(ctx, cnf.DatabaseDsn,
		postgres.WithInitialReconnectBackoffOnFail(time.Second),
	)
	if err != nil {
		logger.Obtain().Fatalf("failed to init postgres client: %v", err)
	}

	return client
}
