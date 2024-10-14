package inmemory

import (
	"context"

	"github.com/aridae/go-metrics-store/internal/server/models"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

func (r *repo) Save(ctx context.Context, metric models.ScalarMetric) error {
	key := metric.Key().String()

	r.storage.Save(ctx, tsstorage.Key(key), metric)

	return nil
}
