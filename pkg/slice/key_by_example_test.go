package slice

import (
	"fmt"
	"strings"
)

// ExampleKeyBy демонстрирует использование функции KeyBy.
func ExampleKeyBy() {
	sliceInt := []int{1, 2, 3}
	fnKeyByInt := func(x int) int { return x * 2 }

	resultInt := KeyBy(sliceInt, fnKeyByInt)
	fmt.Println(resultInt)

	sliceStr := []string{"apple", "banana", "cherry"}
	fnKeyByStr := strings.ToUpper
	resultStr := KeyBy(sliceStr, fnKeyByStr)
	fmt.Println(resultStr)

	// Output:
	// map[2:1 4:2 6:3]
	// map[APPLE:apple BANANA:banana CHERRY:cherry]
}
