package slice

import "fmt"

// ExampleContains демонстрирует использование функции Contains.
func ExampleContains() {
	slice := []int{1, 2, 3}
	target := 2

	found := Contains(slice, target)
	fmt.Println(found)

	target = 4
	found = Contains(slice, target)
	fmt.Println(found)

	// Output:
	// true
	// false
}
