package models

import "errors"

type MetricRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

func (m MetricRequest) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}

	if m.MType == "" {
		return errors.New("type is required")
	}

	return nil
}

type Metric struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
}

func (m Metric) Validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}

	if m.MType == "" {
		return errors.New("type is required")
	}

	return nil
}

type Metrics []Metric

func (ms Metrics) Validate() error {
	for _, m := range ms {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	return nil
}
