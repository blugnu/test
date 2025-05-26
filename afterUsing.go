package test

// AfterUsing is a helper function that allows you to temporarily change the value of a
// variable for the duration of a test.  AfterUsing returns a function which will restore
// the original value.
//
// The intended use is to pass the returned function to Cleanup() or to defer it.
func AfterUsing[T any](v *T, r T) func() {
	og := *v
	*v = r
	return func() { *v = og }
}
