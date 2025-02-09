package slice

import (
	"reflect"
	"testing"
)

func TestKeyBy(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 31},
		{Name: "Dave", Age: 25},
	}

	fn := func(person Person) int {
		return person.Age
	}

	expected := map[int]Person{
		30: {Name: "Alice", Age: 30},
		31: {Name: "Charlie", Age: 31},
		25: {Name: "Dave", Age: 25},
	}

	actual := KeyBy(people, fn)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, actual)
	}
}
