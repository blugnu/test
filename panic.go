package test

import (
	"fmt"
	"runtime"

	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/opt"
)

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
	invalidTest(fmt.Sprintf("Panic: expected at most one argument, got %d", len(r)))
	return panics.Expected{}
}

// NilPanic returns an expectation that a panic will occur that recovers
// a *runtime.PanicNilError.
//
// Panic(nil) is syntactic sugar for "no panic expected", to simplify
// table-drive tests where each test case may or may not expect a panic,
// enabling the use of a single Expect() call.
//
// i.e. instead of writing:
//
//	if testcase.expectedPanic == nil {
//		defer Expect(Panic()).DidNotOccur()
//	} else {
//		defer Expect(Panic(testcase.expectedPanic)).DidOccur()
//	}
//
// you can write:
//
//	defer Expect(Panic(testcase.expectedPanic)).DidOccur()
//
// When testcase.expectedPanic is nil, this is equivalent to:
//
//	defer Expect(Panic()).DidNotOccur()
//
// Without having to write conditional code to handle the different
// expectations.
//
// # Testing for a nil panic
//
// In the unlikely event that you specifically need to test for a
// panic(nil), you can use the NilPanic() function, which will
// create an expectation for a panic that recovers a *runtime.PanicNilError.
//
// see: https://go.dev/blog/compat#expanded-godebug-support-in-go-121
func NilPanic() panics.Expected {
	// This is a convenience function to create an expectation for a nil panic.
	// It is equivalent to Panic(nil).
	return panics.Expected{R: &runtime.PanicNilError{}}
}
