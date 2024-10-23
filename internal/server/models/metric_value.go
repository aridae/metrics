package models

type MetricValue interface {
	String() string
	UnsafeCastInt() int64
	UnsafeCastFloat() float64
	Inc(v MetricValue) (MetricValue, error)
}
