package main

import (
	"context"
	"flag"
	"github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	"github.com/aridae/go-metrics-store/internal/server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/server/usecases/gauge"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
	"log"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "a", ":8080", "Address of server, default: localhost:8080")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	memStore := tsstorage.New()
	metricsRepo := scalarmetric.NewRepository(memStore)

	counterUseCases := counter.NewHandler(metricsRepo)
	gaugeUseCases := gauge.NewHandler(metricsRepo)
	useCaseController := usecases.NewController(metricsRepo, counterUseCases, gaugeUseCases)

	httpRouter := handlers.NewRouter(useCaseController)
	httpServer := http.NewServer(address, httpRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
