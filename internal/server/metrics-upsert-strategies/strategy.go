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

// Strategy стратегия обновления метрики [для каунтера - суммирование, для gauge - оверрайд прошлого значения]
// NOTE: в рамках PR по прошлому спринту была рекомендация использовать паттерн 'стратегия'
// для выполнения upsert-а метрик разных типов: https://github.com/aridae/go-metrics-store/pull/6#discussion_r1754304115
type Strategy interface {
	Upsert(ctx context.Context, metricsRepo metricsRepo, metricToRegister models.ScalarMetricToRegister, now time.Time) (models.ScalarMetric, error)
}
