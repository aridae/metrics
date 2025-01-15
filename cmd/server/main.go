package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/repos/metric"
	"github.com/aridae/go-metrics-store/internal/server/repos/metric/metric-inmem-repo"
	"github.com/aridae/go-metrics-store/internal/server/repos/metric/metric-pg-repo"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw/gzip-mw"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw/logging-mw"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw/sha256-mw"
	"github.com/aridae/go-metrics-store/pkg/logger"
	nooptrm "github.com/aridae/go-metrics-store/pkg/noop-trm"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/config"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	"github.com/aridae/go-metrics-store/pkg/inmem"
	"github.com/aridae/go-metrics-store/pkg/postgres"
	trmman "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

		<-signalCh

		logger.Infof("Got signal, shutting down...")

		// If you fail to cancel the context, the goroutine that WithCancel or WithTimeout created
		// will be retained in memory indefinitely (until the program shuts down), causing a memory leak.
		cancel()
	}()

	cnf := config.Obtain()

	var txManager trm.Manager
	var metricRepo metric.Repository
	var routerOptions []handlers.RouterOption

	if cnf.DatabaseDsn != "" {
		pgClient := mustInitPostgresClient(ctx, cnf)
		txManager = trmman.Must(trmpgx.NewDefaultFactory(pgClient))

		var err error
		metricRepo, err = metricpgrepo.NewRepositoryImplementation(ctx, pgClient, trmpgx.DefaultCtxGetter)
		if err != nil {
			logger.Fatalf("failed to init metricRepo: %v", err)
		}

		routerOptions = append(routerOptions, handlers.CheckAvailableOnPing(pgClient))
	}

	if metricRepo == nil {
		memStore := mustInitMemStore(ctx, cnf)
		metricRepo = metricinmemrepo.NewRepositoryImplementation(memStore)
		txManager = nooptrm.NewNoopTransactionManager()
	}

	useCaseController := usecases.NewController(metricRepo, txManager)

	httpRouter := handlers.NewRouter(useCaseController, routerOptions...)

	httpServer := http.NewServer(cnf.Address, httpRouter,
		gzipmw.GzipDecompressRequestMiddleware,
		sha256mw.ValidateRequestServerMiddleware(cnf.Key),
		sha256mw.SignResponseServerMiddleware(cnf.Key),
		gzipmw.GzipCompressResponseMiddleware,
		loggingmw.LoggingMiddleware,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}

func mustInitMemStore(ctx context.Context, cnf *config.Config) *inmem.MemTimeseriesStorage {
	memStore := inmem.New()

	err := memStore.InitBackup(ctx, cnf.FileStoragePath, cnf.StoreInterval, map[string]any{
		"Metric":             models.Metric{},
		"Int64MetricValue":   models.NewInt64MetricValue(0),
		"Float64MetricValue": models.NewFloat64MetricValue(0),
	})
	if err != nil {
		logger.Fatalf("failed to init mem store backup: %v", err)
	}

	if cnf.Restore {
		err = memStore.LoadFromBackup()
		if err != nil {
			logger.Fatalf("failed to load mem store from backup: %v", err)
		}
	}

	return memStore
}

func mustInitPostgresClient(ctx context.Context, cnf *config.Config) *postgres.Client {
	client, err := postgres.NewClient(ctx, cnf.DatabaseDsn,
		postgres.WithInitialReconnectBackoffOnFail(time.Second),
	)
	if err != nil {
		logger.Fatalf("failed to init postgres client: %v", err)
	}

	return client
}
