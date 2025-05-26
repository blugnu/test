package opt

// AnyAreSet checks if any of the provided options are set in the opts slice.
// It returns true if at least one of the options is found, otherwise false.
func AnyAreSet(opts []any, opt ...any) bool {
	for _, o := range opts {
		for _, oo := range opt {
			if o == oo {
				return true
			}
		}
	}
	return false
}
