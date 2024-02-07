package test

import (
	"errors"
	"testing"
)

// ErrorFormat values are used to specify the format of
// errors displayed in test failure reports.
//
// Supported formats are:
//
//   - ErrorDefault		(%v)
//   - ErrorString		(%s)
//   - ErrorDecl		(%#v)
type ErrorFormat string

const (
	errorFormatNotSet ErrorFormat = ""    // zero value for ErrorFormat, indicates not set
	ErrorDefault      ErrorFormat = "%v"  // errors are printed using %v representation
	ErrorString       ErrorFormat = "%s"  // errors are printed using %s (their Error() string)
	ErrorDecl         ErrorFormat = "%#v" // errors are printed using %#v (Go declaration)
)

// provides methods for testing an error value.
type ErrorTest struct {
	testable[error]
}

// returns a value that may be used to apply tests to a specified error.  A name for
// the error must be provided.  An optional ErrorFormat value may be provided to specify
// the format of errors displayed in test failure reports.  If an ErrorFormat is not
// provided, ErrorDefault is used (%v).
func Error(t *testing.T, got error, opts ...any) ErrorTest {
	n := "error"
	f := ErrorDefault
	checkOptTypes(t, optTypes(n, f), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	return ErrorTest{newTestable(t, got, n, Format(f))}
}

// fails the test if the error being tested does not satisfy errors.Is()
// with respect to the specified error:
//
//   - if wanted is nil, got must be nil
//   - if wanted is not nil, got must satisfy errors.Is()
//     with respect to wanted
//
// Example:
//
//	test.Error(t, "returned", err).Is(ErrExpected)
func (et ErrorTest) Is(wanted error) {
	et.Helper()

	if wanted == nil && et.got == nil {
		return
	}

	if wanted == nil {
		et.IsNil()
		return
	}

	et.run(func(t *testing.T) {
		t.Helper()
		if !errors.Is(et.got, wanted) {
			et.errorf(t, "wanted error: %s\ngot         : %s", et.format(wanted), et.format(et.got))
		}
	})
}

// fails the test if got is not nil.
//
// If got is not nil the test fails and a test failure report will be output:
//
//	unexpected error: <%T of got>: <%v of got>
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		var err error
//
//		// ACT
//		err = doSomething()
//
//		// ASSERT
//		test.Error(t, "returned", err).IsNil()
//	  }
func (et ErrorTest) IsNil() {
	et.Helper()
	if et.got != nil {
		et.errorf(et.T, "unexpected error: %s", et.format(et.got))
	}
}
