package opt

// Extract removes all instances of type T from a set of options, returning
// the first such value, if present, and the remaining options (in their
// original order).
//
// An indicator value is also returned to indicate whether any value of
// type T was found.  If the indicator is false, the returned T value will
// be the zero value of T and should be ignored.
func Extract[T any](opts []any) (T, []any, bool) {
	var (
		extracted T
		found     bool
		remaining = make([]any, 0, len(opts))
	)

	for _, opt := range opts {
		if v, ok := opt.(T); ok {
			if !found {
				extracted = v
				found = true
			}
			continue
		}

		remaining = append(remaining, opt)
	}

	return extracted, remaining, found
}
