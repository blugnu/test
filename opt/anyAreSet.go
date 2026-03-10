package opt

import "slices"

// AnyAreSet checks if any of the provided options are set in the opts slice.
// It returns true if at least one of the options is found, otherwise false.
func AnyAreSet(opts []any, opt ...any) bool {
	for _, o := range opt {
		if slices.Contains(opts, o) {
			return true
		}
	}
	return false
}
