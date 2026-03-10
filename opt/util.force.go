package opt

type forcedPosition int

const (
	// ForceAtStart indicates that the forced option should be placed at the
	// start of the options slice.
	ForceAtStart forcedPosition = iota

	// ForceAtEnd indicates that the forced option should be placed at the
	// end of the options slice.
	ForceAtEnd
)

// Force returns a slice of options guaranteed to contain a specified instance of
// an option type T. The returned array is based on a provided opts slice.
//
// Any instances of type T in the opts slice are removed and replaced by the
// specified option instance.
//
// By default, the specified option is placed at the end of the returned options
// slice.  This may be changed by specifying ForceAtStart in the optional pos
// argument.
func Force[T any](opts []any, opt T, pos ...forcedPosition) []any {
	if len(pos) == 0 {
		pos = []forcedPosition{ForceAtEnd}
	}

	if pos[0] == ForceAtStart {
		return append([]any{opt}, removeAll[T](opts)...)
	}

	return append(removeAll[T](opts), opt)
}

// removeAll removes all instances of type T from the provided opts slice. If
// there are no instances of type T, a copy of the original opts slice is returned.
func removeAll[T any](opts []any) []any {
	result := make([]any, 0, len(opts))
	for _, opt := range opts {
		if _, ok := opt.(T); ok {
			continue
		}
		result = append(result, opt)
	}

	return result
}
