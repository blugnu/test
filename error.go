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
	ErrorDefault      ErrorFormat = "%v"  // errors are printed using the default representation
	ErrorString       ErrorFormat = "%s"  // errors are printed using their Error() string
	ErrorDecl         ErrorFormat = "%#v" // errors are printed using their Go declaration
)

// ErrorIs fails the test if got does not satisfy errors.Is()
// with respect to a wanted error:
//
//   - if wanted is nil, got must be nil
//   - if wanted is not nil, got must satisfy errors.Is()
//     with respect to wanted
//
// An optional ErrorFormat value may be provided to specify
// the format of errors displayed in test failure reports.
// If not provided, ErrorDefault is used (%v).
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
//		test.ErrorIs(t, nil, err)
//	  }
func ErrorIs(t *testing.T, wanted, got error, opt ...ErrorFormat) {
	t.Helper()

	f := Format(ErrorDefault)
	if len(opt) > 0 {
		f = Format(opt[0])
	}

	if wanted == nil && got != nil {
		t.Errorf("\nunexpected error: %s", format(got, f))
	} else if !errors.Is(got, wanted) {
		t.Errorf("\nwanted error: %s\ngot         : %s", format(wanted, f), format(got, f))
	}
}

// UnexpectedError fails the test if got is not nil.  The
// test is equivalent to test.Error(t, nil, got).
//
// An optional ErrorFormat value may be provided to specify
// the format of the error displayed in a test failure report.
// If not provided, ErrorDefault is used (%v).
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
//		test.UnexpectedError(t, err)
//	  }
func UnexpectedError(t *testing.T, got error, opt ...ErrorFormat) {
	t.Helper()
	ErrorIs(t, nil, got, opt...)
}
