package metricsfabrics

import (
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

// ScalarMetricFactory абстрактная фабрика метрик и связанных сущностей.
// NOTE: в рамках PR по прошлому спринту была рекомендация использовать
// фабрику для создания метрик разных типов: https://github.com/aridae/go-metrics-store/pull/6#discussion_r1754304094
type ScalarMetricFactory interface {
	CreateMetricKey(metricName string) models.MetricKey
	CastScalarMetricValue(v any) (models.ScalarMetricValue, error)
	ParseScalarMetricValue(v string) (models.ScalarMetricValue, error)
	CreateScalarMetricToRegister(name string, val models.ScalarMetricValue) models.ScalarMetricToRegister
	ProvideUpsertStrategy() metricsupsertstrategies.Strategy
}
