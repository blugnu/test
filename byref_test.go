package test

// byref is a test utility returning a typed pointer to some value.  It is used in
// some test cases where pointers to literals would otherwise be helpful.
func byref[T any](v T) *T {
	return &v
}
