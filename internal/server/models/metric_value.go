package models

type MetricValue interface {
	String() string
	Inc(v MetricValue) (MetricValue, error)
}
