package inmem

import (
	"context"
	"sync"
	"time"
)

type file interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Truncate(size int64) error
	Seek(offset int64, whence int) (ret int64, err error)
	Close() error
}

type Storage[Key comparable, Value any] struct {
	storeMu sync.RWMutex
	store   map[Key]Value

	backupFileMu   sync.RWMutex
	backupFile     file
	backupInterval time.Duration
}

func New[Key comparable, Value any]() *Storage[Key, Value] {
	return &Storage[Key, Value]{store: make(map[Key]Value)}
}

func (s *Storage[Key, Value]) Save(_ context.Context, key Key, value Value) {
	s.storeMu.Lock()
	defer s.storeMu.Unlock()

	s.store[key] = value
}

func (s *Storage[Key, Value]) Get(_ context.Context, key Key) (Value, bool) {
	s.storeMu.RLock()
	defer s.storeMu.RUnlock()

	val, ok := s.store[key]

	return val, ok
}

func (s *Storage[Key, Value]) GetAll(_ context.Context) []Value {
	s.storeMu.RLock()
	defer s.storeMu.RUnlock()

	vals := make([]Value, 0, len(s.store))
	for _, v := range s.store {
		vals = append(vals, v)
	}

	return vals
}
