package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/repos"
	"github.com/aridae/go-metrics-store/internal/server/repos/inmem-driven-repos/metric-inmem-repo"
	"github.com/aridae/go-metrics-store/internal/server/repos/pg-driven-repos/metric-pg-repo"
	pgtxman "github.com/aridae/go-metrics-store/internal/server/repos/pg-driven-repos/pg-tx-man"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/mw"
	"github.com/aridae/go-metrics-store/pkg/logger"
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
)

/*
TODO добавить транзакционную модель по тутору https://threedots.tech/post/database-transactions-in-go/
TODO почитать вот эту статью https://threedots.tech/post/repository-pattern-in-go/
*/

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

	var txMan repos.TransactionManager
	var metricRepo repos.MetricRepository
	var routerOptions []handlers.RouterOption

	if cnf.DatabaseDsn != "" {
		pgClient := mustInitPostgresClient(ctx, cnf)

		var err error
		metricRepo, err = metricpgrepo.NewRepositoryImplementation(ctx, pgClient)
		if err != nil {
			logger.Obtain().Fatalf("failed to init metricRepo: %v", err)
		}

		txMan = pgtxman.NewTransactionManagerImplementation(pgClient)
		routerOptions = append(routerOptions, handlers.CheckAvailableOnPing(pgClient))
	}

	if metricRepo == nil {
		memStore := mustInitMemStore(ctx, cnf)
		metricRepo = metricinmemrepo.NewRepositoryImplementation(memStore)
		txMan = repos.NewNoopTransactionManager(&repos.Repositories{MetricRepository: metricRepo})
	}

	useCaseController := usecases.NewController(metricRepo, txMan)

	httpRouter := handlers.NewRouter(useCaseController, routerOptions...)

	httpServer := http.NewServer(cnf.Address, httpRouter,
		mw.LoggingMiddleware,
		mw.GzipDecompressRequestMiddleware,
		mw.GzipCompressResponseMiddleware,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Obtain().Fatalf("failed to start server: %v", err)
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
