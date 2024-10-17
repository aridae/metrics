package inmemory

import (
	"context"
	scalarmetric "github.com/aridae/go-metrics-store/internal/server/repos/scalar-metric"
	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

type timeseriesStorage interface {
	Save(ctx context.Context, key tsstorage.Key, value tsstorage.TimeseriesValue)
	GetLatest(ctx context.Context, key tsstorage.Key) tsstorage.TimeseriesValue
	GetAllLatest(ctx context.Context) []tsstorage.TimeseriesValue
}

type repo struct {
	storage timeseriesStorage
}

func NewRepositoryImplementation(storage timeseriesStorage) scalarmetric.Repository {
	return &repo{storage: storage}
}
