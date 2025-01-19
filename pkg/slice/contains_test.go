package slice

import (
	"testing"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string        // Название теста
		slice    []interface{} // Входной срез
		target   interface{}   // Целевой элемент
		expected bool          // Ожидаемый результат
	}{
		{
			name:     "IntSliceFound",
			slice:    []interface{}{1, 2, 3},
			target:   2,
			expected: true,
		},
		{
			name:     "IntSliceNotFound",
			slice:    []interface{}{4, 5, 6},
			target:   7,
			expected: false,
		},
		{
			name:     "StringSliceFound",
			slice:    []interface{}{"apple", "banana", "cherry"},
			target:   "banana",
			expected: true,
		},
		{
			name:     "StringSliceNotFound",
			slice:    []interface{}{"mango", "orange", "pear"},
			target:   "grapes",
			expected: false,
		},
		{
			name:     "Float64SliceFound",
			slice:    []interface{}{1.0, 2.0, 3.0},
			target:   2.0,
			expected: true,
		},
		{
			name:     "Float64SliceNotFound",
			slice:    []interface{}{4.0, 5.0, 6.0},
			target:   7.0,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual := Contains(
				tc.slice,
				tc.target,
			)

			if actual != tc.expected {
				t.Errorf("Expected %v, got %v for slice %v and target %v", tc.expected, actual, tc.slice, tc.target)
			}
		})
	}
}
