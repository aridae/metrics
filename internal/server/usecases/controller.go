package usecases

import (
	"context"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type metricsRepo interface {
	Save(ctx context.Context, metric models.ScalarMetric) error
	GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error)
	GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error)
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
