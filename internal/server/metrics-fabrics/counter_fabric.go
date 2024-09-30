package metricsfabrics

import (
	"fmt"
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"strconv"
)

type counterMetricFactory struct{}

func NewCounterMetricFactory() ScalarMetricFactory {
	return &counterMetricFactory{}
}

func (f *counterMetricFactory) CreateMetricKey(metricName string) models.MetricKey {
	return models.MetricKey(models.ScalarMetricTypeCounter.String() + ":" + metricName)
}

func (f *counterMetricFactory) CastScalarMetricValue(v any) (models.ScalarMetricValue, error) {
	switch value := v.(type) {
	case int64:
		return models.NewInt64MetricValue(value), nil
	case int32:
		return models.NewInt64MetricValue(int64(value)), nil
	case int:
		return models.NewInt64MetricValue(int64(value)), nil
	case int16:
		return models.NewInt64MetricValue(int64(value)), nil
	case int8:
		return models.NewInt64MetricValue(int64(value)), nil
	default:
		return nil, fmt.Errorf("can't cast to int numeric safely, unsupported scalar value type: %T", value)
	}
}

func (f *counterMetricFactory) ParseScalarMetricValue(v string) (models.ScalarMetricValue, error) {
	int64Val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse '%s' to int64 value: %w", v, err)
	}

	return models.NewInt64MetricValue(int64Val), nil
}

func (f *counterMetricFactory) CreateScalarMetricToRegister(key models.MetricKey, val models.ScalarMetricValue) models.ScalarMetricToRegister {
	return models.NewScalarMetricToRegister(key, val, models.ScalarMetricTypeCounter)
}

func (f *counterMetricFactory) ProvideUpsertStrategy() metricsupsertstrategies.Strategy {
	return metricsupsertstrategies.NewIncrementUpsertStrategy()
}
