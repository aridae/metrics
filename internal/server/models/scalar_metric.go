package models

import "time"

type ScalarMetricType string

func (t ScalarMetricType) String() string {
	return string(t)
}

const (
	ScalarMetricTypeCounter ScalarMetricType = "counter"
	ScalarMetricTypeGauge   ScalarMetricType = "gauge"
)

type ScalarMetricToRegister struct {
	name  string
	val   ScalarMetricValue
	mtype ScalarMetricType
}

func NewScalarMetricToRegister(name string, val ScalarMetricValue, mtype ScalarMetricType) ScalarMetricToRegister {
	return ScalarMetricToRegister{
		name:  name,
		val:   val,
		mtype: mtype,
	}
}

func (s ScalarMetricToRegister) Key() MetricKey {
	return BuildMetricKey(s.name, s.mtype)
}

func (s ScalarMetricToRegister) Name() string {
	return s.name
}

func (s ScalarMetricToRegister) Value() ScalarMetricValue {
	return s.val
}

func (s ScalarMetricToRegister) Type() ScalarMetricType {
	return s.mtype
}

func (s ScalarMetricToRegister) WithValue(v ScalarMetricValue) ScalarMetricToRegister {
	s.val = v // local copy only
	return s
}

func (s ScalarMetricToRegister) AtDatetime(now time.Time) ScalarMetric {
	return ScalarMetric{
		ScalarMetricToRegister: s,
		dt:                     now,
	}
}

type ScalarMetric struct {
	ScalarMetricToRegister
	dt time.Time
}

func (s ScalarMetric) Datetime() time.Time {
	return s.dt
}
