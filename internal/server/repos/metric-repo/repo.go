package metricrepo

import (
	"context"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.Metric, error)
	GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error)
	Save(ctx context.Context, metric models.Metric) error
}
