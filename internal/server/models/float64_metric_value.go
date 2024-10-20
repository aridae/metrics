package models

import (
	"fmt"
	"strconv"
)

type float64MetricValue struct {
	Val float64
}

func NewFloat64MetricValue(val float64) MetricValue {
	return float64MetricValue{Val: val}
}

func (mv float64MetricValue) String() string {
	return strconv.FormatFloat(mv.Val, 'f', -1, 64)
}

func (mv float64MetricValue) Inc(v MetricValue) (MetricValue, error) {
	float64Val, ok := v.(float64MetricValue)
	if !ok {
		return nil, fmt.Errorf("expected float64 metric value, got %T", v)
	}

	newVal := mv.Val + float64Val.Val

	return float64MetricValue{Val: newVal}, nil
}
