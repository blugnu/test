package test

// AddressOf returns the address of the value passed as argument.  If the
// argument is a value-type variable, the address will be that of a copy
// of the value, not the value itself.
//
// This function is provided for situations where a test requires that the
// address of some value is provided; it as a convenience for the fact that
// Golang does not allow the & operator to be used with literals or constants.
//
// # Example
//
//	func TestSomething(t *testing.T) {
//		// given that SomeFunction returns a struct:
//		type resultModel struct {
//			v *string
//		}
//
//		// and that SomeFunction is expected to return a resultModel with
//		// v set to the address of a supplied string:
//		result := SomeFunction("value")
//
//		// then the test could be written as:
//		test.That(t, result).Equals(resultModel{
//			v: test.AddressOf("value"),
//		})
//	}
func AddressOf[T any](v T) *T {
	return &v
}
