package opt

// Get returns an option of a given type from a slice of options. Only the
// first value of type T in the slice is significant; any additional T values
// are ignored.
//
// The function returns the option value and true if an option of the
// desired type was identified, otherwise the zero value of T is returned
// with false.
func Get[T any](opts []any) (T, bool) {
	for _, opt := range opts {
		if v, ok := opt.(T); ok {
			return v, true
		}
	}

	var zero T
	return zero, false
}
