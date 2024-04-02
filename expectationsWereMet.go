package test

import "testing"

// Mock is an interface describing the methods required to be supported by a
// mock that can be tested using test.MockExpectations().
type Mock interface {
	ExpectationsWereMet() error
	Resetter
}

// ExpectationsWereMet is a convenience function for testing and resetting mocks.
// It verifies that all expectations set on the mocks have been met and resets
// the mocks to their initial state.  If any expectations were not met an error
// is reported.  All mocks are tested, even if some mocks have failed to meet
// their expectations.
//
// It is recommended that the test function should call this function using the
// defer statement to ensure that the mocks are tested and reset after the test
// function completes.
//
// This function support mocks that implement the test.Mock interface:
//
//	type Mock interface {
//		ExpectationsWereMet() error
//		Resetter
//	}
//
//	type Resetter interface {
//		Reset()
//	}
//
// The Resetter interface is used to reset the mocks to their initial state and
// is also used to support resetting fakes (see: test.Fake[R]).  A mock is always
// reset after expectations have been tested, regardless of whether the test
// passed or failed.
//
// Where a mock is used to provide a more complex fake implementation of an
// interface but details (i.e. expectations) of how the interface is used are not
// of interest to this test, the Resetter interface can be used to reset the
// mock without testing the expectations.
//
// # Example
//
//	func TestSomething(t *testing.T) {
//		// ARRANGE
//		mock1 := NewMock()
//		mock2 := NewMock()
//		defer test.ExpectationsWereMet(t, mock1, mock2)
//
//		// .. set expectations on mock1 and mock2
//
//		// ACT
//		// .. call the code under test
//
//		// ASSERT
//		// .. assertions additional to testing mock expectations (if any)
//	}
func ExpectationsWereMet(t *testing.T, m ...Mock) {
	t.Helper()

	// ensure that all mocks are reset after all expectations have been tested
	defer func() {
		for _, mock := range m {
			mock.Reset()
		}
	}()

	for _, mock := range m {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}
}
