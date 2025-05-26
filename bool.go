package test

import "github.com/blugnu/test/matchers/bools"

// BeFalse returns a matcher that will fail if the matched value is not false.
//
// # Supported Options
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
func BeFalse() bools.BooleanMatcher {
	return bools.BooleanMatcher{Expected: false}
}

// BeTrue returns a matcher that will fail if the matched value is not true.
//
// # Supported Options
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
func BeTrue() bools.BooleanMatcher {
	return bools.BooleanMatcher{Expected: true}
}

// ExpectFalse fails a test if a specified bool is not false.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(false))
//	Expect(got).To(BeFalse())
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test
//	                         // failure report
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
func ExpectFalse[T ~bool](got T, opts ...any) {
	GetT().Helper()
	Expect(bool(got), opts...).To(BeFalse(), opts...)
}

// ExpectTrue fails a test if a specified bool is not true.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(true))
//	Expect(got).To(BeTrue())
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test
//	                         // failure report
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
func ExpectTrue[T ~bool](got T, opts ...any) {
	GetT().Helper()
	Expect(bool(got), opts...).To(BeTrue(), opts...)
}
