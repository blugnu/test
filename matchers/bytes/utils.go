package bytes

import "cmp"

func max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func remove[T comparable](s []T, e T) []T {
	result := make([]T, 0, len(s))
	for _, s := range s {
		if s == e {
			continue
		}
		result = append(result, s)
	}
	return result
}
