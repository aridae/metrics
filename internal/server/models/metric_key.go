package models

type MetricKey string

func (k MetricKey) String() string {
	return string(k)
}

func BuildMetricKey(name string, mtype MetricType) MetricKey {
	return MetricKey(mtype.String() + ":" + name)
}
