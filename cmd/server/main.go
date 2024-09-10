package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/config"
	"github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	"github.com/aridae/go-metrics-store/internal/server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/server/usecases/gauge"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
	"log"
)

func main() {
	ctx := context.Background()
	cnf := config.ObtainFromFlags()

	memStore := tsstorage.New()
	metricsRepo := scalarmetric.NewRepository(memStore)

	counterUseCases := counter.NewHandler(metricsRepo)
	gaugeUseCases := gauge.NewHandler(metricsRepo)
	useCaseController := usecases.NewController(metricsRepo, counterUseCases, gaugeUseCases)

	httpRouter := handlers.NewRouter(useCaseController)
	httpServer := http.NewServer(cnf.GetAddress(), httpRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
