package pointer

import "fmt"

// ExampleTo показывает базовый пример использования функции To.
func ExampleTo() {
	v := To(42, 0)
	fmt.Println(*v)

	v = To(0, 0)
	fmt.Println(v)

	// Output:
	// 42
	// <nil>
}
