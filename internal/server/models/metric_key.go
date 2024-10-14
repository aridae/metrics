package models

type MetricKey string

func (k MetricKey) String() string {
	return string(k)
}

func BuildMetricKey(name string, mtype ScalarMetricType) MetricKey {
	return MetricKey(mtype.String() + ":" + name)
}
