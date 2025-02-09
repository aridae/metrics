package slice

func KeyBy[T any, K comparable](slice []T, fn func(T) K) map[K]T {
	result := make(map[K]T)

	for _, elem := range slice {
		key := fn(elem)
		result[key] = elem
	}

	return result
}
