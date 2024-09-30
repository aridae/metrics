package metricsupsertstrategies

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"time"
)

type metricsRepo interface {
	GetLatestState(ctx context.Context, metricKey models.MetricKey) (*models.ScalarMetric, error)
	Save(ctx context.Context, metric models.ScalarMetric) error
}

type Strategy interface {
	Upsert(ctx context.Context, metricsRepo metricsRepo, metricToRegister models.ScalarMetricToRegister, now time.Time) error
}
