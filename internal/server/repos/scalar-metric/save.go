package scalarmetric

import (
	"context"
	"github.com/aridae/go-metrics-store/internal/server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

func (r *Repository) Save(ctx context.Context, metric models.ScalarMetric) error {
	key := metric.Key().String()

	r.storage.Save(ctx, tsstorage.Key(key), metric)

	return nil
}
