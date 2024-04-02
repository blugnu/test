package test

// Resetter is an interface that describes the Reset method. A Resetter
// is any type that can be reset to some initial state.
type Resetter interface {
	Reset()
}

// Reset calls the Reset method on each Resetter in the list. A Resetter
// is any type that implements the Resetter interface:
//
//	type Resetter interface {
//		Reset()
//	}
//
// Reset is a convenience function for resetting multiple Resetters,
// such as fakes and mocks.
//
// # Example
//
//	func TestSomething(t *testing.T) {
//		// ARRANGE
//		fake := NewFake()
//		mock := NewMock()
//		defer test.Reset(fake, mock)
//
//		// .. set expectations on mock
//
//		// ACT
//		// .. call the code under test
//
//		// ASSERT
//		// .. assertions additional to testing mock expectations (if any)
//		test.MockExpectations(t, mock)
//	}
func Reset(r ...Resetter) {
	for _, resetter := range r {
		resetter.Reset()
	}
}
