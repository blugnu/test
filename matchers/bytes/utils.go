package bytes

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
