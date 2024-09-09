package tsstorage

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Key string

type TimeseriesValue interface {
	GetDatetime() time.Time
}

type MemTimeseriesStorage struct {
	mu    sync.RWMutex
	store map[Key][]TimeseriesValue
}

func New() *MemTimeseriesStorage {
	return &MemTimeseriesStorage{
		store: make(map[Key][]TimeseriesValue),
	}
}

func (mem *MemTimeseriesStorage) Save(_ context.Context, key Key, value TimeseriesValue) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.store[key] = append(mem.store[key], value)
}

func (mem *MemTimeseriesStorage) GetLatest(_ context.Context, key Key) TimeseriesValue {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	timeseries, found := mem.store[key]
	if !found || len(timeseries) == 0 {
		return nil
	}

	sort.SliceStable(timeseries, func(i, j int) bool {
		return timeseries[i].GetDatetime().Before(timeseries[j].GetDatetime())
	})

	return timeseries[len(timeseries)-1]
}
