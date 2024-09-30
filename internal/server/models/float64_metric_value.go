package models

import (
	"fmt"
)

type float64MetricValue struct {
	val float64
}

func NewFloat64MetricValue(val float64) ScalarMetricValue {
	return float64MetricValue{val: val}
}

func (mv float64MetricValue) String() string {
	return fmt.Sprintf("%f", mv.val)
}

func (mv float64MetricValue) Inc(v ScalarMetricValue) (ScalarMetricValue, error) {
	float64Val, ok := v.(float64MetricValue)
	if !ok {
		return nil, fmt.Errorf("expected float64 metric value, got %T", v)
	}

	newVal := mv.val + float64Val.val

	return float64MetricValue{val: newVal}, nil
}
