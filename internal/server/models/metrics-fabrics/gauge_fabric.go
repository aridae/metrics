package metricsfabrics

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

var (
	_gaugeFactory *gaugeMetricFactory
	_gaugeOnce    sync.Once
)

type gaugeMetricFactory struct{}

func ObtainGaugeMetricFactory() MetricFactory {
	_gaugeOnce.Do(func() {
		_gaugeFactory = &gaugeMetricFactory{}
	})

	return _gaugeFactory
}

func (f *gaugeMetricFactory) CreateMetricKey(name string) models.MetricKey {
	return models.BuildMetricKey(name, models.MetricTypeGauge)
}

func (f *gaugeMetricFactory) ParseMetricValue(v string) (models.MetricValue, error) {
	float64Val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse '%s' to int64 value: %w", v, err)
	}

	return models.NewFloat64MetricValue(float64Val), nil
}

func (f *gaugeMetricFactory) CreateMetricUpsert(name string, val models.MetricValue) models.MetricUpsert {
	return models.NewMetricUpsert(name, val, models.MetricTypeGauge)
}
