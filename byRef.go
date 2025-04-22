package test

// ByRef returns the address of the value passed as argument.  If the
// argument is a value-type variable, the address will be that of a copy
// of the value, not the value itself.
//
// This function is provided for situations where a test requires that the
// address of some value is provided; it as a convenience for the fact that
// Golang does not allow the & operator to be used with literals or constants.
func ByRef[T any](v T) *T {
	return &v
}
