package test

// Cleanup adds a function to the list of functions to be called when the test
// finishes.  The function will be called in reverse order of registration.
//
// If a nil function is passed, it is ignored.
//
// This is a wrapper around the *testing.T.Cleanup() method with the addition
// of accepting (but ignoring) a nil function.
func Cleanup(fn func()) {
	if fn == nil {
		return // NO-OP
	}
	T().Cleanup(fn)
}
