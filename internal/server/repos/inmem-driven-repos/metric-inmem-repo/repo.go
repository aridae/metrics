package metricinmemrepo

import (
	"context"
	metricrepo "github.com/aridae/go-metrics-store/internal/server/repos"
	"github.com/aridae/go-metrics-store/pkg/inmem"
)

type db interface {
	Save(ctx context.Context, key inmem.Key, value inmem.TimeseriesValue)
	GetLatest(ctx context.Context, key inmem.Key) inmem.TimeseriesValue
	GetAllLatest(ctx context.Context) []inmem.TimeseriesValue
}

type repo struct {
	db db
}

func NewRepositoryImplementation(db db) metricrepo.MetricRepository {
	return &repo{db: db}
}
