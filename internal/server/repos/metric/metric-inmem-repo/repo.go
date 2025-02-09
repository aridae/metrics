package metricinmemrepo

import (
	"context"
	"sort"

	"github.com/aridae/go-metrics-store/internal/server/models"
	metricrepo "github.com/aridae/go-metrics-store/internal/server/repos/metric"
)

type inmemoryStorage[Key comparable, Value any] interface {
	Save(ctx context.Context, key Key, val Value)
	Get(ctx context.Context, key Key) (Value, bool)
	GetAll(ctx context.Context) []Value
}

type repo struct {
	store inmemoryStorage[models.MetricKey, models.Metric]
}

func NewRepositoryImplementation(
	db inmemoryStorage[models.MetricKey, models.Metric],
) metricrepo.Repository {
	return &repo{store: db}
}

func (r *repo) Save(ctx context.Context, metric models.Metric) error {
	key := metric.GetKey()

	r.store.Save(ctx, key, metric)

	return nil
}

func (r *repo) GetByKey(ctx context.Context, key models.MetricKey) (*models.Metric, error) {
	val, isFound := r.store.Get(ctx, key)
	if !isFound {
		return nil, nil
	}

	return &val, nil
}

func (r *repo) GetAll(ctx context.Context) ([]models.Metric, error) {
	metrics := r.store.GetAll(ctx)

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].GetKey() < metrics[j].GetKey()
	})

	return metrics, nil
}
