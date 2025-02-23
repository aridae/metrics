package pointer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTo(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		nillif   []int
		expected *int
	}{
		{
			name:     "Nil value",
			value:    0,
			nillif:   []int{0},
			expected: nil,
		},
		{
			name:     "Non-nil value",
			value:    42,
			nillif:   []int{0},
			expected: &[]int{42}[0],
		},
		{
			name:     "Empty nillif list",
			value:    0,
			nillif:   []int{},
			expected: &[]int{0}[0],
		},
		{
			name:     "Multiple nillif values",
			value:    7,
			nillif:   []int{0, 7},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := To(tt.value, tt.nillif...)
			require.Equal(t, result, tt.expected)
		})
	}
}

// BenchmarkTo измеряет производительность функции To.
func BenchmarkTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		To(i, 0)
	}
}
