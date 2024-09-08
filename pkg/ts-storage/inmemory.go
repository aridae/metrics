package tsstorage

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Key string

type Value interface {
	Time() time.Time
}

type MemStorage struct {
	mu    sync.RWMutex
	store map[Key][]Value
}

func New() *MemStorage {
	return &MemStorage{
		store: make(map[Key][]Value),
	}
}

func (inmem *MemStorage) Save(_ context.Context, key Key, value Value) {
	inmem.mu.Lock()
	defer inmem.mu.Unlock()

	inmem.store[key] = append(inmem.store[key], value)
}

func (inmem *MemStorage) GetLatest(_ context.Context, key Key) Value {
	inmem.mu.RLock()
	defer inmem.mu.RUnlock()

	timeseries, found := inmem.store[key]
	if !found || len(timeseries) == 0 {
		return nil
	}

	sort.SliceStable(timeseries, func(i, j int) bool {
		return timeseries[i].Time().Before(timeseries[j].Time())
	})

	return timeseries[len(timeseries)-1]
}
