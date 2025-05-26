package opt

// Get extracts an option of a given type from a variadic list of options.
// Only the first value of type T is significant; any additional T values
// are ignored.
//
// The function returns the option value and true if an option of the
// desired type was identified, otherwise the zero value of the specified
// option type is returned with false.
func Get[T any](opts []any) (T, bool) {
	for _, opt := range opts {
		switch v := opt.(type) {
		case T:
			return v, true
		}
	}
	return *new(T), false
}
