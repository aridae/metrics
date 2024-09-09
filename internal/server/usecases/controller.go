package usecases

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"github.com/aridae/go-metrics-store/internal/server/usecases/counter"
	"github.com/aridae/go-metrics-store/internal/server/usecases/gauge"
)

type metricsRepo interface {
	GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error)
	GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error)
}

type Controller struct {
	metricsRepo            metricsRepo
	counterUseCasesHandler *counter.Handler
	gaugeUseCasesHandler   *gauge.Handler
}

func NewController(
	metricsRepo metricsRepo,
	counterUseCasesHandler *counter.Handler,
	gaugeUseCasesHandler *gauge.Handler,
) *Controller {
	return &Controller{
		metricsRepo:            metricsRepo,
		counterUseCasesHandler: counterUseCasesHandler,
		gaugeUseCasesHandler:   gaugeUseCasesHandler,
	}
}
