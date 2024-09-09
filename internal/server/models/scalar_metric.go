package models

import "time"

type MetricKey string

func (s MetricKey) String() string {
	return string(s)
}

type ScalarMetricType string

const (
	ScalarMetricTypeCounter ScalarMetricType = "counter"
	ScalarMetricTypeGauge   ScalarMetricType = "gauge"
)

func (s ScalarMetricType) String() string {
	return string(s)
}

type CounterValue int64
type GaugeValue float64

type ScalarMetricUpdater struct {
	Type  ScalarMetricType
	Name  string
	Value any
}

func (su ScalarMetricUpdater) Key() MetricKey {
	return MetricKey(su.Type.String() + ":" + su.Name)
}

func (su ScalarMetricUpdater) AsCounterValue() CounterValue {
	if val, ok := su.Value.(CounterValue); ok {
		return val
	}

	if rawVal, ok := su.Value.(int64); ok {
		return CounterValue(rawVal)
	}

	return CounterValue(0)
}

func (su ScalarMetricUpdater) AsGaugeValue() GaugeValue {
	if val, ok := su.Value.(GaugeValue); ok {
		return val
	}

	if rawVal, ok := su.Value.(float64); ok {
		return GaugeValue(rawVal)
	}

	return GaugeValue(0)
}

type ScalarMetric struct {
	ScalarMetricUpdater
	Datetime time.Time
}

func (s ScalarMetric) GetDatetime() time.Time {
	return s.Datetime
}
