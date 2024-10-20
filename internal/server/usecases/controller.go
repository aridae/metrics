package usecases

import (
	"context"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type metricsRepo interface {
	Save(ctx context.Context, metric models.Metric) error
	GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error)
	GetAll(ctx context.Context) ([]models.Metric, error)
}

type Controller struct {
	metricsRepo metricsRepo
	now         func() time.Time
}

func NewController(
	metricsRepo metricsRepo,
) *Controller {
	return &Controller{
		metricsRepo: metricsRepo,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}
