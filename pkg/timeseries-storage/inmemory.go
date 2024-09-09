package timeseriesstorage

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

	timeseries := mem.store[key]

	timeseries = append(timeseries, value)

	sort.SliceStable(timeseries, func(i, j int) bool {
		return timeseries[i].GetDatetime().Before(timeseries[j].GetDatetime())
	})

	mem.store[key] = timeseries
}

func (mem *MemTimeseriesStorage) GetLatest(_ context.Context, key Key) TimeseriesValue {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	timeseries := mem.store[key]
	if len(timeseries) == 0 {
		return nil
	}

	return timeseries[len(timeseries)-1]
}
