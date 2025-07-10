package test

import (
	"github.com/blugnu/test/matchers/mocks"
	"github.com/blugnu/test/opt"
)

// Mock is an interface describing the methods required to be supported by a
// mock that can be tested using ExpectationsWereMet().
type Mock interface {
	ExpectationsWereMet() error
	Resetter
}

// ExpectationsWereMet is a convenience function for testing and resetting expectations
// of a mock that implements the test.Mock interface:
//
//	type Mock interface {
//		ExpectationsWereMet() error
//		Reset()
//	}
//
// The mock is guaranteed to be reset whether expecations were met or not.
//
// If a test is not concerned with expectations being met and is using a mock
// simply to provide complex mocked responses, the mock can be reset without
// checking expectations by passing the mock to the Reset() helper.
//
// # Supported Options
//
//	string    // a name for the expectation; the name is used in
//	          // the failure message if the expectation fails.
//
// // # Example
//
//	func TestSomething(t *testing.T) {
//		With(t)
//
//		// ARRANGE
//		mock1 := NewMock()
//		mock2 := NewMock()
//		mock3 := NewMock()
//		defer Reset(mock1, mock2, mock3)
//
//		// .. configure expectations on mocks
//
//		// ACT
//		// .. call the code under test
//
//		// ASSERT
//		// .. additional assertions prior to the deferred testing
//		test.ExpectationsWereMet(mock1, "mock 1")
//		test.ExpectationsWereMet(mock2, "mock 2")
//		test.ExpectationsWereMet(mock3, "mock 3")
//	}
func ExpectationsWereMet(m Mock, opts ...any) {
	T().Helper()

	defer m.Reset()

	// if err is non-nil we use an IsNil test to report the
	// failure with any options applied
	if err := m.ExpectationsWereMet(); err != nil {
		Expect(err).IsNil(append(opts, opt.OnFailure(err.Error()))...)
	}
}

// MeetExpectations is a matcher that checks whether the expectations of a mock
// were met.  It is used in conjunction with the Expect() function to assert
func MeetExpectations() *mocks.Matcher {
	return &mocks.Matcher{}
}
