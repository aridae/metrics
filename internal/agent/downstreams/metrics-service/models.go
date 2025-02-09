package metricsservice

// Metric локальная модель метрики, с которой работает Client.
type Metric struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}
