package slice

import (
	"fmt"
	"strings"
)

// ExampleMapBatch демонстрирует использование функции MapBatch.
func ExampleMapBatch() {
	inputsInt := []int{1, 2, 3}
	mapperFnInt := func(x int) (int, error) {
		return x * 2, nil
	}

	resultsInt, _ := MapBatch(inputsInt, mapperFnInt)
	fmt.Println(resultsInt)

	inputsStr := []string{"apple", "banana", "cherry"}
	mapperFnStr := func(x string) (string, error) {
		return strings.ToUpper(x), nil
	}

	resultsStr, _ := MapBatch(inputsStr, mapperFnStr)
	fmt.Println(resultsStr)

	// Output:
	// [2 4 6]
	// [APPLE BANANA CHERRY]
}
