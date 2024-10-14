package usecases

import (
	"context"
	"time"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type healthcheckable interface {
	Healthcheck(ctx context.Context) error
}

type metricsRepo interface {
	healthcheckable
	Save(ctx context.Context, metric models.ScalarMetric) error
	GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error)
	GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error)
}

type Controller struct {
	metricsRepo  metricsRepo
	postgresConn healthcheckable
	now          func() time.Time
}

func NewController(
	metricsRepo metricsRepo,
	postgresConn healthcheckable,
) *Controller {
	return &Controller{
		metricsRepo:  metricsRepo,
		postgresConn: postgresConn,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}
