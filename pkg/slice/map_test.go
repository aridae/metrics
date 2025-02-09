package slice

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapBatch(t *testing.T) {
	type mapperFunc[In any, Out any] func(In) (Out, error)

	var testCases = []struct {
		expectedErr error
		mapperFunc  mapperFunc[int, int]
		name        string
		input       []int
		expected    []int
	}{
		{
			name:        "SuccessAllElementsMapped",
			input:       []int{1, 2, 3},
			mapperFunc:  func(i int) (int, error) { return i + 10, nil },
			expected:    []int{11, 12, 13},
			expectedErr: nil,
		},
		{
			name:  "ErrorOnFirstElement",
			input: []int{1, 2, 3},
			mapperFunc: func(i int) (int, error) {
				if i == 1 {
					return 0, errors.New("error on first element")
				} else {
					return i + 20, nil
				}
			},
			expected:    nil,
			expectedErr: errors.New("error on first element"),
		},
		{
			name:        "NoInput",
			input:       []int{},
			mapperFunc:  func(i int) (int, error) { return i + 30, nil },
			expected:    []int{},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapBatch(tc.input, tc.mapperFunc)

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.EqualValues(t, tc.expected, result)
			}
		})
	}
}
