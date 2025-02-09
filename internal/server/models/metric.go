package models

import "time"

type MetricType string

func (t MetricType) String() string {
	return string(t)
}

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

type MetricUpsert struct {
	Val   MetricValue
	MName string
	Mtype MetricType
}

func NewMetricUpsert(name string, val MetricValue, mtype MetricType) MetricUpsert {
	return MetricUpsert{
		MName: name,
		Val:   val,
		Mtype: mtype,
	}
}

func (s MetricUpsert) GetKey() MetricKey {
	return BuildMetricKey(s.MName, s.Mtype)
}

func (s MetricUpsert) GetName() string {
	return s.MName
}

func (s MetricUpsert) GetValue() MetricValue {
	return s.Val
}

func (s MetricUpsert) GetType() MetricType {
	return s.Mtype
}

func (s MetricUpsert) WithValue(v MetricValue) MetricUpsert {
	s.Val = v // local copy only
	return s
}

func (s MetricUpsert) WithDatetime(now time.Time) Metric {
	return Metric{
		MetricUpsert: s,
		Datetime:     now,
	}
}

type Metric struct {
	MetricUpsert
	Datetime time.Time
}

func (s Metric) GetDatetime() time.Time {
	return s.Datetime
}
