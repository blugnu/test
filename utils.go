package test

import "golang.org/x/exp/constraints"

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func remove[T comparable](s []T, e T) []T {
	out := make([]T, len(s)-1)
	ie := 0
	for _, s := range s {
		if s == e {
			continue
		}
		out[ie] = s
		ie++
	}
	return out
}
