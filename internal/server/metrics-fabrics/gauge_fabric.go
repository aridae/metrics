package metricsfabrics

import (
	"fmt"
	"strconv"
	"sync"

	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
)

var (
	_gaugeFactory *gaugeMetricFactory
	_gaugeOnce    sync.Once
)

type gaugeMetricFactory struct{}

func ObtainGaugeMetricFactory() ScalarMetricFactory {
	_gaugeOnce.Do(func() {
		_gaugeFactory = &gaugeMetricFactory{}
	})

	return _gaugeFactory
}

func (f *gaugeMetricFactory) CreateMetricKey(name string) models.MetricKey {
	return models.BuildMetricKey(name, models.ScalarMetricTypeGauge)
}

func (f *gaugeMetricFactory) ParseScalarMetricValue(v string) (models.ScalarMetricValue, error) {
	float64Val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse '%s' to int64 value: %w", v, err)
	}

	return models.NewFloat64MetricValue(float64Val), nil
}

func (f *gaugeMetricFactory) CastScalarMetricValue(v any) (models.ScalarMetricValue, error) {
	switch value := v.(type) {
	case float64:
		return models.NewFloat64MetricValue(value), nil
	case float32:
		return models.NewFloat64MetricValue(float64(value)), nil
	case int64:
		return models.NewFloat64MetricValue(float64(value)), nil
	case int32:
		return models.NewFloat64MetricValue(float64(value)), nil
	case int:
		return models.NewFloat64MetricValue(float64(value)), nil
	case int16:
		return models.NewFloat64MetricValue(float64(value)), nil
	case int8:
		return models.NewFloat64MetricValue(float64(value)), nil
	default:
		return nil, fmt.Errorf("can't cast to float numeric safely, unsupported scalar value type: %T", value)
	}
}

func (f *gaugeMetricFactory) CreateScalarMetricToRegister(name string, val models.ScalarMetricValue) models.ScalarMetricToRegister {
	return models.NewScalarMetricToRegister(name, val, models.ScalarMetricTypeGauge)
}

func (f *gaugeMetricFactory) ProvideUpsertStrategy() metricsupsertstrategies.Strategy {
	return metricsupsertstrategies.NewOverrideUpsertStrategy()
}
