package models

type MetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

type Metric struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta *string `json:"delta,omitempty"`
	Value *string `json:"value,omitempty"`
}
