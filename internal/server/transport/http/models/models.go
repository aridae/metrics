package models

type MetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
