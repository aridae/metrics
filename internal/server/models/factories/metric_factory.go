package factories

import (
	"github.com/aridae/go-metrics-store/internal/server/models"
)

// MetricFactory абстрактная фабрика метрик и связанных сущностей.
type MetricFactory interface {
	CreateMetricKey(mname string) models.MetricKey
	ParseMetricValue(val string) (models.MetricValue, error)
	CreateMetricUpsert(mname string, val models.MetricValue) models.MetricUpsert
}
