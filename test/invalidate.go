package test

import (
	"fmt"
	"strings"

	"github.com/blugnu/test/internal/testframe"
)

// runner is an interface that defines the method set required to report
// an invalid test or error, extending Helper with error reporting and
// test failure methods
type runner interface {
	Helper
	Errorf(s string, args ...any)
	FailNow()
}

// Error is used to indicate a test is invalid due to some error having occurred.
// This is intended to be used in test helpers and matchers to report an error
// that invalidates a test.
//
// i.e. the error does not indicate that the test failed, but rather that
// the test is invalid and therefore unreliable, due to an error that occurred
// during execution or evaluation of a test helper or matcher.
//
// It should not be confused with the [github.com/blugnu/test.Error] or
// [github.com/blugnu/test.Errorf] functions.
//
// If a valid test frame is available, it will report the error using the [Errorf]
// method of that test frame.
//
// If no valid test frame is available, this function will panic with the
// provided error and message(s).  Panicking ensures that test execution fails,
// avoiding a false positive outcome.
//
// # Alternatives
//
// To indicate that a test is invalid without any specific error, use the [Invalid]
// function.
//
// To draw attention to a non-fatal issue in a test, use the [Warning] function.
func Error(err error, msg ...string) {
	if s := strings.Join(msg, "\n"); len(s) > 0 {
		err = fmt.Errorf("%w\n%s", err, s)
	}

	if t, ok := testframe.Peek[runner](); ok {
		t.Helper()
		t.Errorf("<== INVALID TEST\nERROR: %s", err.Error())
		t.FailNow()
		return
	}

	panic(fmt.Errorf("INVALID TEST\n%w", err))
}

// Invalid is used to mark a test as invalid.  It should be called by a matcher
// when the test cannot be run due to an invalid condition, such as attempting to
// use a matcher with an unsupported type, or when the test is not properly set up.
//
// Calling this function will fail the current test with an error message that
// includes the provided message(s).  If no message is provided, it will simply
// mark the test as invalid without any additional information.
//
// The Go standard library testing package does not provide a way to mark a test
// as invalid, so this function is used to provide a consistent way to do so.
//
// An invalid test is identified by a "<== INVALID TEST" error following the
// test location, followed by any message provided.
func Invalid(msg ...string) {
	msg = append([]string{"INVALID TEST"}, msg...)
	s := strings.Join(msg, "\n")

	// if we can obtain a TestRunner from the current test frame then we will
	// use it to report the test as invalid, otherwise we must panic, to avoid
	// a test yielding a false positive result
	if t, ok := testframe.Peek[runner](); ok {
		t.Helper()
		t.Errorf("<== %s", s)
		t.FailNow()
		return
	}

	panic(s)
}

// Warning is used to report a warning in a test.  This should be used to
// indicate a condition that is not an error, but may indicate a problem or
// unexpected behavior in the test.
//
// Although a warning does not invalidate a test, it does fail any current
// test execution.
//
// For example, if all test cases in a table-driven test are skipped, a
// warning is produced indicating that the test did not run any test cases.
// This does not represent an error or failure in any individual test case,
// but fails the test, avoiding a false positive result.
//
// If a valid test frame is available, it will report the warning using the
// Errorf, otherwise it will panic with the warning message.
//
// The warning message will be prefixed with "WARNING: " to indicate that it is
// a warning and not an error.  This is useful for indicating that the test
// is not invalid, but there is something noteworthy that should be considered
// by the developer.
func Warning(msg string) {
	msg = "WARNING: " + msg

	// if we can obtain a TestRunner from the current test frame then we will
	// use it to report the test as invalid, otherwise we must panic, to avoid
	// a test yielding a false positive result
	if t, ok := testframe.Peek[runner](); ok {
		t.Helper()
		t.Errorf("<== " + msg)
		return
	}

	panic(msg)
}
