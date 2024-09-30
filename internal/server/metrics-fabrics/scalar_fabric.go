package metricsfabrics

import (
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

type ScalarMetricFactory interface {
	CreateMetricKey(metricName string) models.MetricKey
	CastScalarMetricValue(v any) (models.ScalarMetricValue, error)
	ParseScalarMetricValue(v string) (models.ScalarMetricValue, error)
	CreateScalarMetricToRegister(key models.MetricKey, val models.ScalarMetricValue) models.ScalarMetricToRegister
	ProvideUpsertStrategy() metricsupsertstrategies.Strategy
}
