package metric

import (
	"context"
	tsstorage "github.com/aridae/go-metrics-store/pkg/ts-storage"
)

type storage interface {
	Save(_ context.Context, key tsstorage.Key, value tsstorage.Value)
	GetLatest(_ context.Context, key tsstorage.Key) tsstorage.Value
}

type Repository struct {
	storage storage
}

func NewRepository(storage storage) *Repository {
	return &Repository{storage: storage}
}
