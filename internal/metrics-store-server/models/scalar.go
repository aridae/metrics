package models

import "time"

type ScalarMetricType string

const (
	ScalarMetricTypeCounter ScalarMetricType = "counter"
	ScalarMetricTypeGauge   ScalarMetricType = "gauge"
)

type CounterValue int64
type GaugeValue float64

type ScalarMetricUpdater struct {
	Type  ScalarMetricType
	Name  string
	Value any
}

func (su ScalarMetricUpdater) AsCounterValue() CounterValue {
	return su.Value.(CounterValue)
}

func (su ScalarMetricUpdater) AsGaugeValue() GaugeValue {
	return su.Value.(GaugeValue)
}

type ScalarMetric struct {
	ScalarMetricUpdater
	Datetime time.Time
}

func (s ScalarMetric) Time() time.Time {
	return s.Datetime
}
