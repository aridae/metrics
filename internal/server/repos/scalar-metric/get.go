package scalarmetric

import (
	"context"
	"fmt"
	"github.com/aridae/go-metrics-store/internal/server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

func (r *Repository) GetLatestState(ctx context.Context, key models.MetricKey) (*models.ScalarMetric, error) {
	val := r.storage.GetLatest(ctx, tsstorage.Key(key.String()))
	if val == nil {
		return nil, nil
	}

	scalar, ok := val.(*models.ScalarMetric)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T, expecting *models.ScalarMetric type", val)
	}

	return scalar, nil
}
