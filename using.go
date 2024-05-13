package test

// Using is a helper function that allows you to temporarily change the value of a variable.
// The original value is restored when the returned function is called.
//
// # example
//
//	func TestSomething(t *testing.T) {
//	  // ARRANGE
//	  defer test.Using(&now, func() time.Time { return time.Date(2020,1,1,0,0,0,0,time.UTC) })()
//
//	  // ACT + ASSERT ...
//
//	  // execute some code which uses the now variable function to get the current time
//	  // and assert that the code behaves as expected given the time is 2020-01-01
//	}
func Using[T any](v *T, r T) func() {
	og := *v
	*v = r
	return func() { *v = og }
}
