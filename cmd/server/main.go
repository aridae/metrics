package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/repos/metric"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/transport/http"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/metrics-store-server/usecases/gauge"
	tsstorage "github.com/aridae/go-metrics-store/pkg/ts-storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Сервер должен быть доступен по адресу http://localhost:8080
// Принимать метрики по протоколу HTTP методом POST.
// Принимать данные в формате
//    http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>,
//    Content-Type: text/plain.
//

const (
	address = ":8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	memStore := tsstorage.New()

	metricsRepo := metric.NewRepository(memStore)

	counterUseCases := counter.NewHandler(metricsRepo)
	gaugeUseCases := gauge.NewHandler(metricsRepo)

	useCaseController := usecases.NewController(counterUseCases, gaugeUseCases)

	httpRouter := handlers.NewRouter(useCaseController)
	httpServer := http.NewServer(address, httpRouter)

	go func() {
		defer cancel()
		_ = httpServer.Run(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-quit:
			if err := httpServer.Shutdown(ctx); err != nil {
				log.Fatalf("Server shutdown failed: %v", err)
			}
		case <-ctx.Done():
		}
	}

}
