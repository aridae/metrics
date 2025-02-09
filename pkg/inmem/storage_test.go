package inmem

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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

// BenchmarkSave измеряет производительность метода Save.
func BenchmarkSave(b *testing.B) {
	ctx := context.Background()

	// Создаем хранилище
	s := New[string, string]()

	keys := make([]string, b.N)
	values := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = strconv.Itoa(i)
		values[i] = "test value"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Save(ctx, keys[i], values[i])
	}
}

// BenchmarkGet измеряет производительность метода Get.
func BenchmarkGet(b *testing.B) {
	ctx := context.Background()

	// Создаем хранилище и заполняем его данными
	s := New[string, string]()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		value := "test value"
		s.Save(ctx, key, value)
	}

	keys := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = strconv.Itoa(rand.Intn(b.N))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(ctx, keys[i])
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

// BenchmarkGetAll измеряет производительность метода GetAll.
func BenchmarkGetAll(b *testing.B) {
	ctx := context.Background()

	// Создаем хранилище и заполняем его данными
	s := New[string, string]()
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		value := "test string"
		s.Save(ctx, key, value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.GetAll(ctx)
	}
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
