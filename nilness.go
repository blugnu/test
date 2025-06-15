package test

import "github.com/blugnu/test/matchers/nilness"

// IsNil checks that the value of the expectation is nil.  If the
// value is not nil, the test fails.  If the value is nil, the test
// passes.
//
// If the value being tested does not support a nil value the test
// will fail and produce a report similar to:
//
//	test.IsNil: values of type '<type>' are not nilable
//
// # Supported Options
//
//	opt.QuotedStrings(bool)     // determines whether any non-nil string
//	                            // values are quoted in any test failure
//	                            // report.  The default is false (string
//	                            // values are quoted).
//	                            //
//	                            // If the value is not a string type this
//	                            // option has no effect.
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsNil(opts ...any) {
	e.t.Helper()
	e.Should(BeNil(), opts...)
}

// IsNotNil checks that a specified value is not nil.  If the value
// is not nil, the test fails.  If the value is nil, the test passes.
//
// NOTE: If the value being tested does not support a nil value the
// test will pass.  This is to allow for testing values that may be
// nilable or non-nilable.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func (e expectation[T]) IsNotNil(opts ...any) {
	e.t.Helper()
	e.ShouldNot(BeNil(), opts...)
}

// BeNil returns a matcher that checks if the value is nil.
//
// The returned matcher is an `AnyMatcher` that may only be used
// with the `Should()` method, or with `To()` where the subject
// is of formal type any.
//
// If used in a To() or Should() test with a subject that is not
// nilable, the test fails with a message similar to:
//
//	test.BeNil: values of type '<type>' are not nilable
//
// If used with ToNo() or ShouldNot() a non-nilable subject will
// pass the test.
//
// # Supported Options
//
//	opt.QuotedStrings(bool)     // determines whether any non-nil string
//	                            // values are quoted in any test failure
//	                            // report.  The default is false (string
//	                            // values are quoted).
//	                            //
//	                            // If the value is not a string type this
//	                            // option has no effect.
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func BeNil() nilness.Matcher {
	return nilness.Matcher{}
}
