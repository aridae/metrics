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
	MName string
	Mtype ScalarMetricType
	Val   ScalarMetricValue
}

func NewScalarMetricToRegister(name string, val ScalarMetricValue, mtype ScalarMetricType) ScalarMetricToRegister {
	return ScalarMetricToRegister{
		MName: name,
		Val:   val,
		Mtype: mtype,
	}
}

func (s ScalarMetricToRegister) Key() MetricKey {
	return BuildMetricKey(s.MName, s.Mtype)
}

func (s ScalarMetricToRegister) Name() string {
	return s.MName
}

func (s ScalarMetricToRegister) Value() ScalarMetricValue {
	return s.Val
}

func (s ScalarMetricToRegister) Type() ScalarMetricType {
	return s.Mtype
}

func (s ScalarMetricToRegister) WithValue(v ScalarMetricValue) ScalarMetricToRegister {
	s.Val = v // local copy only
	return s
}

func (s ScalarMetricToRegister) AtDatetime(now time.Time) ScalarMetric {
	return ScalarMetric{
		ScalarMetricToRegister: s,
		Datetime:               now,
	}
}

type ScalarMetric struct {
	ScalarMetricToRegister
	Datetime time.Time
}

func (s ScalarMetric) GetDatetime() time.Time {
	return s.Datetime
}
