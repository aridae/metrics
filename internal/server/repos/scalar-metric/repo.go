package scalarmetric

import (
	"context"

	tsstorage "github.com/aridae/go-metrics-store/pkg/timeseries-storage"
)

type timeseriesStorage interface {
	Save(ctx context.Context, key tsstorage.Key, value tsstorage.TimeseriesValue)
	GetLatest(ctx context.Context, key tsstorage.Key) tsstorage.TimeseriesValue
	GetAllLatest(ctx context.Context) []tsstorage.TimeseriesValue
}

type Repository struct {
	storage timeseriesStorage
}

func NewRepository(storage timeseriesStorage) *Repository {
	return &Repository{storage: storage}
}
