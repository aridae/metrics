package inmem

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

// TestSaveAndGet проверяет корректность сохранения и получения значения по ключу.
func TestSaveAndGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	storage := New[string, string]()

	storage.Save(ctx, key, value)

	gotVal, found := storage.Get(ctx, key)
	if !found || gotVal != value {
		t.Errorf("Expected to find %q with value %q, but found %v", key, value, gotVal)
	}
}

// TestGetNonExistent проверяет получение несуществующего ключа.
func TestGetNonExistent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	key := "non_existent_key"

	storage := New[string, string]()

	_, found := storage.Get(ctx, key)
	if found {
		t.Error("Expected not to find non-existent key")
	}
}

// TestGetAll проверяет получение всех значений из хранилища.
func TestGetAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := New[int, string]()

	keys := []int{1, 2, 3}
	expectedValues := []string{"one", "two", "three"}

	for i, k := range keys {
		storage.Save(ctx, k, expectedValues[i])
	}

	actualValues := storage.GetAll(ctx)

	require.ElementsMatch(t, expectedValues, actualValues)
}

// TestConcurrentAccess проверяет безопасность конкурентного доступа к хранилищу.
func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	storage := New[int, string]()

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			storage.Save(ctx, i, "value")
		}(i)
	}

	wg.Wait()

	allValues := storage.GetAll(ctx)

	if len(allValues) != numGoroutines {
		t.Fatalf("Expected %d values, got %d", numGoroutines, len(allValues))
	}
}
