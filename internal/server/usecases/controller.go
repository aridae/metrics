package usecases

import (
	"context"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type transactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type metricsRepo interface {
	Save(ctx context.Context, metric models.Metric) error
	GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error)
	GetAll(ctx context.Context) ([]models.Metric, error)
}

type Controller struct {
	metricsRepo        metricsRepo
	transactionManager transactionManager
	now                func() time.Time
}

func NewController(
	metricsRepo metricsRepo,
	transactionManager transactionManager,
) *Controller {
	return &Controller{
		metricsRepo:        metricsRepo,
		transactionManager: transactionManager,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}
