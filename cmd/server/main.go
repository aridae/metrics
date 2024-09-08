package main

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/repos/metric"
	"github.com/aridae/go-metrics-store/internal/server/transport/http"
	"github.com/aridae/go-metrics-store/internal/server/transport/http/handlers"
	"github.com/aridae/go-metrics-store/internal/server/usecases"
	"github.com/aridae/go-metrics-store/internal/server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/server/usecases/gauge"
	tsstorage "github.com/aridae/go-metrics-store/pkg/ts-storage"
	"log"
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
	ctx := context.Background()

	memStore := tsstorage.New()

	metricsRepo := metric.NewRepository(memStore)

	counterUseCases := counter.NewHandler(metricsRepo)
	gaugeUseCases := gauge.NewHandler(metricsRepo)

	useCaseController := usecases.NewController(counterUseCases, gaugeUseCases)

	httpRouter := handlers.NewRouter(useCaseController)
	httpServer := http.NewServer(address, httpRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
