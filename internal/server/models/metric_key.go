package models

type MetricKey string

func (k MetricKey) String() string {
	return string(k)
}
