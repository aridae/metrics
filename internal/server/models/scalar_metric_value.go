package models

type ScalarMetricValue interface {
	String() string
	Inc(v ScalarMetricValue) (ScalarMetricValue, error)
}
