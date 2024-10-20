package inmem

import (
	"context"
	"os"
	"sort"
	"sync"
	"time"
)

type Key string

func (k Key) String() string {
	return string(k)
}

type TimeseriesValue interface {
	GetDatetime() time.Time
}

type MemTimeseriesStorage struct {
	storeMu sync.RWMutex
	store   map[Key][]TimeseriesValue

	fileMu         sync.RWMutex
	backupFile     *os.File
	backupInterval time.Duration
}

func New() *MemTimeseriesStorage {
	return &MemTimeseriesStorage{
		store: make(map[Key][]TimeseriesValue),
	}
}

func (mem *MemTimeseriesStorage) Save(_ context.Context, key Key, value TimeseriesValue) {
	mem.storeMu.Lock()
	defer mem.storeMu.Unlock()

	timeseries := mem.store[key]

	timeseries = append(timeseries, value)

	sort.SliceStable(timeseries, func(i, j int) bool {
		return timeseries[i].GetDatetime().Before(timeseries[j].GetDatetime())
	})

	mem.store[key] = timeseries
}

func (mem *MemTimeseriesStorage) GetLatest(_ context.Context, key Key) TimeseriesValue {
	mem.storeMu.RLock()
	defer mem.storeMu.RUnlock()

	timeseries := mem.store[key]
	if len(timeseries) == 0 {
		return nil
	}

	return timeseries[len(timeseries)-1]
}

func (mem *MemTimeseriesStorage) GetAllLatest(_ context.Context) []TimeseriesValue {
	mem.storeMu.RLock()
	defer mem.storeMu.RUnlock()

	res := make([]TimeseriesValue, 0, len(mem.store))
	for _, timeseries := range mem.store {
		if len(timeseries) == 0 {
			continue
		}
		res = append(res, timeseries[len(timeseries)-1])
	}

	return res
}
