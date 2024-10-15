package metricsfabrics

import (
	"fmt"
	"strconv"
	"sync"

	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

var (
	_counterFactory *counterMetricFactory
	_counterOnce    sync.Once
)

type counterMetricFactory struct{}

func ObtainCounterMetricFactory() ScalarMetricFactory {
	_counterOnce.Do(func() {
		_counterFactory = &counterMetricFactory{}
	})

	return _counterFactory
}

func (f *counterMetricFactory) CreateMetricKey(metricName string) models.MetricKey {
	return models.BuildMetricKey(metricName, models.ScalarMetricTypeCounter)
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

func (f *counterMetricFactory) CreateScalarMetricToRegister(name string, val models.ScalarMetricValue) models.ScalarMetricToRegister {
	return models.NewScalarMetricToRegister(name, val, models.ScalarMetricTypeCounter)
}

func (f *counterMetricFactory) ProvideUpsertStrategy() metricsupsertstrategies.Strategy {
	return metricsupsertstrategies.NewIncrementUpsertStrategy()
}
