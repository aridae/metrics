package pointer

import (
	"github.com/aridae/go-metrics-store/pkg/slice"
)

func To[T comparable](v T, nillif ...T) *T {
	if slice.Contains(nillif, v) {
		return nil
	}

	return &v
}
