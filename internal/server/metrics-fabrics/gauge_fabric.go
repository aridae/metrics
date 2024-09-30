package metricsfabrics

import (
	"fmt"
	metricsupsertstrategies "github.com/aridae/go-metrics-store/internal/server/metrics-upsert-strategies"
	"github.com/aridae/go-metrics-store/internal/server/models"
	"strconv"
)

type gaugeMetricFactory struct{}

func NewGaugeMetricFactory() ScalarMetricFactory {
	return &gaugeMetricFactory{}
}

func (f *gaugeMetricFactory) CreateMetricKey(metricName string) models.MetricKey {
	return models.MetricKey("gauge:" + metricName)
}

func (f *gaugeMetricFactory) ParseScalarMetricValue(v string) (models.ScalarMetricValue, error) {
	float64Val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse '%s' to int64 value: %w", v, err)
	}

	return models.NewFloat64MetricValue(float64Val), nil
}

func (f *gaugeMetricFactory) CreateScalarMetricToRegister(key models.MetricKey, val models.ScalarMetricValue) models.ScalarMetricToRegister {
	return models.NewScalarMetricToRegister(key, val, models.ScalarMetricTypeGauge)
}

func (f *gaugeMetricFactory) ProvideUpsertStrategy() metricsupsertstrategies.Strategy {
	return metricsupsertstrategies.NewOverrideUpsertStrategy()
}
