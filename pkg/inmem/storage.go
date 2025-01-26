package inmem

import (
	"context"
	"encoding/gob"
	"sync"
	"time"
)

type encoder interface {
	Encode(e any) error
}

type decoder interface {
	Decode(e any) error
}

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

	provideFileEncoder  func(backupFile file) encoder
	providerFileDecoder func(backupFile file) decoder
}

func New[Key comparable, Value any]() *Storage[Key, Value] {
	return &Storage[Key, Value]{
		store: make(map[Key]Value),
		provideFileEncoder: func(backupFile file) encoder {
			return gob.NewEncoder(backupFile)
		},
		providerFileDecoder: func(backupFile file) decoder {
			return gob.NewDecoder(backupFile)
		},
	}
}

// Save сохраняет пару ключ-значение в хранилище.
//
// Аргументы:
// key (Key): Ключ для сохранения.
// value (Value): Значение для сохранения.
func (s *Storage[Key, Value]) Save(_ context.Context, key Key, value Value) {
	s.storeMu.Lock()
	defer s.storeMu.Unlock()

	s.store[key] = value
}

// Get возвращает значение по указанному ключу.
//
// Аргументы:
// _ (context.Context): Контекст выполнения запроса (не используется).
// key (Key): Ключ для поиска значения.
//
// Возвращает:
// Value: Найденное значение.
// bool: Признак наличия ключа в хранилище.
func (s *Storage[Key, Value]) Get(_ context.Context, key Key) (Value, bool) {
	s.storeMu.RLock()
	defer s.storeMu.RUnlock()

	val, ok := s.store[key]

	return val, ok
}

// GetAll возвращает все значения из хранилища.
//
// Аргументы:
// _ (context.Context): Контекст выполнения запроса (не используется).
//
// Возвращает:
// []Value: Все сохраненные значения.
func (s *Storage[Key, Value]) GetAll(_ context.Context) []Value {
	s.storeMu.RLock()
	defer s.storeMu.RUnlock()

	vals := make([]Value, 0, len(s.store))
	for _, v := range s.store {
		vals = append(vals, v)
	}

	return vals
}
