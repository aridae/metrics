package metricsfabrics

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/aridae/go-metrics-store/internal/server/models"
)

var (
	_counterFactory *counterMetricFactory
	_counterOnce    sync.Once
)

type counterMetricFactory struct{}

func ObtainCounterMetricFactory() MetricFactory {
	_counterOnce.Do(func() {
		_counterFactory = &counterMetricFactory{}
	})

	return _counterFactory
}

func (f *counterMetricFactory) CreateMetricKey(metricName string) models.MetricKey {
	return models.BuildMetricKey(metricName, models.MetricTypeCounter)
}

func (f *counterMetricFactory) ParseMetricValue(v string) (models.MetricValue, error) {
	int64Val, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse '%s' to int64 value: %w", v, err)
	}

	return models.NewInt64MetricValue(int64Val), nil
}

func (f *counterMetricFactory) CreateMetricUpsert(name string, val models.MetricValue) models.MetricUpsert {
	return models.NewMetricUpsert(name, val, models.MetricTypeCounter)
}
