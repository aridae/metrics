package metricinmemrepo

import (
	"context"
	metricrepo "github.com/aridae/go-metrics-store/internal/server/repos/metric"
	"github.com/aridae/go-metrics-store/pkg/inmem"
)

type inmemDatabase interface {
	Save(ctx context.Context, key inmem.Key, value inmem.TimeseriesValue)
	GetLatest(ctx context.Context, key inmem.Key) inmem.TimeseriesValue
	GetAllLatest(ctx context.Context) []inmem.TimeseriesValue
}

type repo struct {
	db inmemDatabase
}

func NewRepositoryImplementation(db inmemDatabase) metricrepo.Repository {
	return &repo{db: db}
}
