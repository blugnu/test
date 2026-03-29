package test

import (
	"fmt"
	"runtime"

	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/matchers/typecheck"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
	"github.com/blugnu/test/test"
)

// MARK: error

// BeError returns a matcher that checks if the value is an error of the expected type.
//
// The matcher may be used in two modes:
//
//   - When called with an expected error value and type E of `error`, the test will
//     pass if the subject is a non-nil error that satisfies [errors.Is] for the
//     expected error;
//
//   - When called with no expected error value and a specific type E (that implements
//     `error`), the test will pass if the subject is a non-nil error that satisfies
//     [errors.As] for the type E.
//
// Calling BeError with 2 or more error values is invalid and will cause the test to
// fail with an appropriate error message.
//
// # Alternatives
//
// To perform further tests on a value that passes this check using a strongly
// typed value of type E, use:
//
//	expect.Error[E](err)    // returns an E and an indicator whether the test passed
//	                        // or failed, without halting test execution on failure
//
//	require.Error[E](err)   // returns an E when the test passes; halts current test
//	                        // execution when the test fails, avoiding the need
//	                        // for a returned indicator
func BeError[E error](err ...E) Matcher[error] {
	switch len(err) {
	case 0:
		// drop out to return the zero value which indicates no specific
		// target error to use errors.As to check for type E
		//
		// NOTE: explicitly returning here would create an unreachable code path
		// because the default case causes the current test to halt execution
		// of the current test

	case 1:
		return typecheck.MatchIsError{Target: err[0]}

	default:
		T().Helper()
		test.Invalid(fmt.Sprintf("BeError: at most one error argument is supported, got %d", len(err)))
	}

	return typecheck.MatchAsError[E]{}
}

// MARK: occurs

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
// the test will pass if any error occurred which may or may not be an
// expected error.  This may be acceptable in very simple cases but
// it is usually better to test for a specific error using one of the
// other methods described below.:
//
//	Expect(err).Is(expectedError)
//	Expect(err).To(BeError(expectedErr))
//	expect.Error[E](err)
func (e expectation[T]) DidOccur(opts ...any) {
	e.t.Helper()

	switch v := any(e.subject).(type) {
	case panics.Expected:
		match := &panics.MatchRecovered{R: recover()}

		if match.R != nil {
			const bufsize = 65536
			stk := make([]byte, bufsize)
			n := runtime.Stack(stk, false)
			match.Stack = stk[:n-1]
		}

		if !match.Match(v, opts...) {
			e.fail(match, opts...)
		}

	case error:
		if v != nil {
			return
		}
		e.err("expected error, got nil")

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
// a panic having occurred.
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
		// for a "DidNotOccur" test, things are more complicated:

		// first let's grab any recoverable value and create a
		// matcher which we'll use later...
		matcher := &panics.MatchRecovered{R: recover()}

		// first, using DidNotOccur with Panic(nil) is invalid since it
		// is likely to cause confusion
		if expected.R == opt.NoPanicExpected(true) {
			test.Invalid("DidNotOccur: may not be used with Panic(nil); did you mean NilPanic()?")
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
			// we add the ToNotMatch(true) option to indicate that the
			// expectation was that the panic should not have occurred
			e.fail(matcher, append(opts, opt.ToNotMatch(true))...)
			return
		}

		// but we're not done yet...
		//
		// the recovered value did not match the expected value, and if that
		// recovered value is not nil, then we have an unexpected panic to report...
		if matcher.R != nil {
			// the existing matcher has already been used to test the recovered
			// value against an expected value where-as we now need to report an
			// unexpected panic (i.e. expected nil)
			//
			// so we create a new panic matcher, matching against an expected R:nil
			// and use THAT to report the failure
			matcher := &panics.MatchRecovered{R: matcher.R}
			matcher.Match(panics.Expected{R: nil})
			e.fail(matcher, opts...)
		}

	case error:
		opts = append(opts, opt.FailReporter(func(opts ...any) []string {
			return []string{
				"expected: <no error>",
				"got     : " + report.Value(expected, opts...),
			}
		}))
		Expect(expected).IsNil(opts...)

	case nil:
		return

	default:
		test.Invalid("test.DidNotOccur: may only be used with Panic() or error values")
	}
}

// MARK: panics

// NilPanic returns an expectation that a panic will occur that recovers
// a *runtime.PanicNilError.
//
// see: https://go.dev/blog/compat#expanded-godebug-support-in-go-121
//
// For more information, refer to the # The Panic(nil) Special Case
// described in the [Panic] function documentation.
func NilPanic() panics.Expected {
	return panics.Expected{R: &runtime.PanicNilError{}}
}

// Panic returns an expectation subject that can be used to test whether a
// panic has occurred, optionally identifying a value that should match the
// value recovered from the expected panic.
//
// NOTE: At most ONE panic test should be expected per function.  In addition,
// extreme care should be exercised when combining panic tests with other
// deferred recover() calls as these will also interfere with a panic test
// (or vice versa).
//
// # Usage
//
//   - If called with no arguments, any panic will satisfy the expectation,
//     regardless of the value recovered.
//
//   - If called with a single argument, it will expect to recover a panic that
//     recovers that value (unless the argument is nil; see The Panic(nil)
//     Special Case, below)
//
//   - If called with > 1 argument, the test will be failed as invalid.
//
// # The Panic(nil) Special Case
//
// Panic(nil) is a special case that is equivalent to "no panic expected".
// This is motivated by table-driven tests to avoid having to write conditional
// code to handle test cases where a panic is expected vs those where not.
//
// Treating Panic(nil) as "no panic expected" allows you to write:
//
//	defer Expect(Panic(testcase.expectedPanic)).DidOccur()
//
// When testcase.expectedPanic is nil, this is equivalent to:
//
//	defer Expect(Panic()).DidNotOccur()
//
// Should you need to test for an actual panic(nil), use:
//
//	defer Expect(NilPanic()).DidOccur()
//
// Or, in a table-driven test, specify an expected recovery value of
// &runtime.PanicNilError{}.
func Panic(r ...any) panics.Expected {
	switch len(r) {
	case 0:
		return panics.Expected{}
	case 1:
		if r[0] == nil {
			return panics.Expected{R: opt.NoPanicExpected(true)}
		}
		return panics.Expected{R: r[0]}
	}

	T().Helper()
	test.Invalid(fmt.Sprintf("Panic: expected at most one argument, got %d", len(r)))

	return panics.Expected{}
}
