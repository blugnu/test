package opt

import "slices"

// HasAnyOf checks if any of the provided options are set in the opts slice.
// It returns true if AT LEAST ONE of the options is found, otherwise false.
//
// No information is returned about WHICH options were found.
func HasAnyOf(opts []any, opt ...any) bool {
	for _, o := range opt {
		if slices.Contains(opts, o) {
			return true
		}
	}

	return false
}
