package inmem

import (
	"context"
	"fmt"
	"sort"
)

func ExampleStorage_Save() {
	storage := New[string, string]()
	key := "my-key"
	value := "my-value"

	storage.Save(context.Background(), key, value)
	gotVal, found := storage.Get(context.Background(), key)
	if !found || gotVal != value {
		panic("value not saved correctly")
	}

	fmt.Printf("Saved and retrieved value for key %q: %q\n", key, gotVal)
	// Output: Saved and retrieved value for key "my-key": "my-value"
}

func ExampleStorage_Get() {
	storage := New[int, string]()
	key := 42
	value := "the answer"

	storage.Save(context.Background(), key, value)
	gotVal, found := storage.Get(context.Background(), key)
	if !found || gotVal != value {
		panic("value not found or incorrect")
	}

	fmt.Printf("Retrieved value for key %d: %q\n", key, gotVal)
	// Output: Retrieved value for key 42: "the answer"
}

func ExampleStorage_GetAll() {
	storage := New[string, string]()
	keysAndValues := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}

	for k, v := range keysAndValues {
		storage.Save(context.Background(), k, v)
	}

	allVals := storage.GetAll(context.Background())

	sort.Strings(allVals)

	fmt.Printf("All values: %+v\n", allVals)
	// Output: All values: [val1 val2 val3]
}
