package slice

import (
	"fmt"
	"testing"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		target   interface{}
		name     string
		slice    []interface{}
		expected bool
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

func BenchmarkContains(b *testing.B) {
	targets := []int{0, 1, 100}

	var testSlices = [][]int{
		{1, 2, 3},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}

	for _, s := range testSlices {
		for _, t := range targets {
			b.Run(fmt.Sprintf("SliceSize:%d_Target:%d", len(s), t), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					Contains(s, t)
				}
			})
		}
	}
}
