package scalarmetric

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type Repository interface {
	Healthcheck(ctx context.Context) error
	GetAllLatestStates(ctx context.Context) ([]models.ScalarMetric, error)
	GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error)
	Save(ctx context.Context, metric models.ScalarMetric) error
}
