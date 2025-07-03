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
	t, ok := testframe.Peek[runner]()
	if !ok {
		panic(s)
	}

	t.Helper()
	t.Errorf("<== %s", s)
	t.FailNow()
}

// Error is used to indicate an error in a test.  This should be used to report
// errors that occur during the execution of a test rendering the test outcome
// invalid.
//
// i.e. the error does not indicate that the test failed, but rather that
// the test is invalid and therefore unreliable, due to an error that occurred
// during its execution.
//
// If a valid test frame is available, it will report the error using the Error
// method.
//
// If no valid test frame is available, this function will panic with the
// provided error and message(s), avoiding a false positive test outcome.
func Error(err error, msg ...string) {
	t, tok := testframe.Peek[runner]()
	s := strings.Join(msg, "\n")

	switch {
	// we have a t we can use to report the error
	case tok && len(s) > 0:
		t.Helper()
		t.Errorf("<== INVALID TEST\nERROR: %s\n%s", err.Error(), s)
	case tok:
		t.Helper()
		t.Errorf("<== INVALID TEST\nERROR: %s", err.Error())

	// no t, we must panic
	case len(s) > 0:
		panic(fmt.Errorf("INVALID TEST\n%w\n%s", err, s))
	default:
		panic(fmt.Errorf("INVALID TEST\n%w", err))
	}
}

// Warning is used to report a warning in a test.  This should be used to
// indicate a condition that is not an error, but may indicate a problem or
// unexpected behavior in the test.
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
	t, ok := testframe.Peek[runner]()
	if !ok {
		panic(msg)
	}

	t.Helper()
	t.Errorf("<== " + msg)
}
