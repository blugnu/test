package test

import (
	"fmt"

	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

// DidOccur is used to check whether an expected panic or error occurred.
//
// # Testing for Panics
//
// Use the Panic(r) function to create an expectation that a value r will
// be recovered from a panic.  The call to DidOccur() must be deferred:
//
//	defer Expect(Panic(r)).DidOccur(opts...)
//
// If the value r is an error the test will pass only if a panic occurs
// and an error is recovered from the panic that satisfies errors.Is(r).
//
// If the expected recovered value is not an error, the test passes if
// the recovered value is equal to the expected value, based on comparison
// using reflect.DeepEqual or a comparison function.
//
// # Supported Options
//
//	func(a, b any) bool     // a function to compare the values, overriding
//	                        // the use of reflect.DeepEqual.
//
// # Testing for Errors
//
// To test for an error, use the error value as the expected value.
// The test will pass if the error is not nil:
//
//	Expect(err).DidOccur()
//
// This is equivalent to:
//
//	Expect(err).IsNotNil()
//
// NOTE: this approach to testing for errors is not recommended since
// the test will pass if any error occured which may not be the error
// that was expected.  This may be acceptable in very simple cases but
// it is usually better to test for a specific error using:
//
//	Expect(err).Is(expectedError)
func (e expectation[T]) DidOccur(opts ...any) {
	e.t.Helper()

	switch v := any(e.subject).(type) {
	case panics.Expected:
		match := &panics.MatchRecoveredValue{R: recover()}
		if !match.Match(v, opts...) {
			e.err(match.OnTestFailure())
		}
	case error:
		if v != nil {
			return
		}
	case nil:
		e.err("expected error, got nil")
	default:
		test.Invalid("test.DidOccur: may only be used with Panic() or error values")
	}
}

// DidNotOccur is used to ensure that a panic or error did not occur.
//
// # Testing for Panics
//
// Use the Panic() function to create an expectation for a Panic with
// an unspecified recovered value.  The call to DidNotOccur() must be
// deferred:
//
//	defer Expect(Panic()).DidNotOccur(opts...)
//
// The test will pass only if the function scope terminates without
// a panic having occured.
//
// # Testing for Errors
//
// To test for an error, use the error value as the expected value.
// The test will pass if the error is nil:
//
//	Expect(err).DidNotOccur()
//
// This is equivalent to:
//
//	Expect(err).IsNil()
func (e expectation[T]) DidNotOccur(opts ...any) {
	e.t.Helper()

	switch expected := any(e.subject).(type) {
	case panics.Expected:
		// for a "DidNotOccur" test, things are more complicated.

		// first let's grab any recoverable value and create a
		// matcher which we'll use later...
		matcher := &panics.MatchRecoveredValue{R: recover()}

		// first, using DidNotOccur with Panic(nil) is invalid since it
		// is likely to cause confusion
		if expected.R == opt.NoPanicExpected(true) {
			test.Invalid("DidNotOccur: may not be used with Panic(nil); did you mean NilPanic()?")

			// if we did not panic, then we can return early, otherwise
			// we will continue to check the recovered value even though
			// the result may be meaningless
			if matcher.R == nil {
				return
			}
		}

		// if we expect Panic(x) did NOT occur, but Panic(y) DID occur,
		// then although the expectation was met, the UNexpected panic
		// should still be reported as a test failure.
		//
		// so we use the MatchesPanic matcher to determine whether
		// the recovered value matches the expected value...
		recoveredExpectedValue := matcher.Match(expected, opts...)

		// if the recovered value matches the expected value, then
		// the test has failed since this panic should not have occurred...
		if recoveredExpectedValue && expected.R != nil {
			// for the error report, we add the ToNotMatch(true) option
			// to indicate that the expectation was that the panic should
			// not have occurred but did
			e.err(matcher.OnTestFailure(append(opts, opt.ToNotMatch(true))...))
			return
		}

		// but we're not done yet...
		//
		// the recovered value did not match the expected value, and if that
		// recovered value is not nil, then we have an unexpected panic to report...
		if matcher.R != nil {
			// the existing matcher has already been used to test the recovered
			// value against an expected value; we cannot use this to produce
			// the failure report since it will incorrectly report "expecting X, got Y"
			// when we want to report "unexpected panic: got Y"
			//
			// so we create a new matcher with the recovered value, match it
			// against an expected R:nil and use THAT to report the failure
			//
			// we could just emit a test report, but that would duplicate the report
			// handling logic in the matcher and require us to also , so we use the matcher to
			matcher := &panics.MatchRecoveredValue{R: matcher.R}
			matcher.Match(panics.Expected{R: nil})
			e.err(matcher.OnTestFailure(opts...))
		}

	case error:
		opts = append(opts, opt.FailureReport(func(opts ...any) []string {
			return []string{
				"expected: <no error>",
				fmt.Sprintf("got     : %T(%v)", expected, opt.ValueAsString(expected, opts...)),
			}
		}))
		Expect(expected).IsNil(opts...)

	case nil:
		return

	default:
		test.Invalid("test.DidNotOccur: may only be used with Panic() or error values")
	}
}
