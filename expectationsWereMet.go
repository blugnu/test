package test

// Mock is an interface describing the methods required to be supported by a
// mock that can be tested using ExpectationsWereMet().
type Mock interface {
	ExpectationsWereMet() error
	Resetter
}

// ExpectationsWereMet is a convenience function for testing and resetting one or
// more mocks implementing the test.Mock interface:
//
//	type Mock interface {
//		ExpectationsWereMet() error
//		Reset()
//	}
//
// All mocks are guaranteed to be reset after the test has completed, regardless of
// whether any mock passed or failed the test.
//
// If a test is not concerned with expectations being met and is using a mock
// simply to provide complex mocked responses, the mock can be reset without
// checking expectations by calling the mock's Reset() method directly.
//
// # Example
//
//	func TestSomething(t *testing.T) {
//		test.With(t)
//
//		// ARRANGE
//		mock1 := NewMock()
//		mock2 := NewMock()
//		mock3 := NewMock()
//		defer test.ExpectationsWereMet(mock1, mock2)
//		defer mock3.Reset()
//
//		// .. configure expectations on mocks
//
//		// ACT
//		// .. call the code under test
//
//		// ASSERT
//		// .. additional assertions prior to the deferred testing
//		//    of mock expectations (if any)...
//	}
func ExpectationsWereMet(m Mock, opts ...any) {
	T().Helper()

	defer m.Reset()

	Expect(m.ExpectationsWereMet(), opts...).IsNil()
}
