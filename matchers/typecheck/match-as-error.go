package typecheck

import (
	"errors"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type MatchAsError[E error] struct {
	Target error
}

// Match checks that the value satisfies [errors.As] for an error of type
// E.  If the value is not of the expected type and does not wrap an error
// of that type, the test fails.
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
func (m MatchAsError[E]) Match(got error, _ ...any) bool {
	var target E
	return errors.As(got, &target)
}

// OnTestFailure returns a failure report indicating that the actual value
// was not of the expected type.
//
// If the ToNotMatch option is set, the report indicates that the value
// should not have been of the expected type.
func (m MatchAsError[E]) OnTestFailure(got any, opts ...any) []string {
	expectedType := report.TypeName[E]()

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{"should not be of type: " + expectedType}
	}

	// if got is nil we include this in the report as no error at all
	// vs an error of the expected type is a significant difference
	if got == nil {
		return []string{
			"expected: " + expectedType,
			"got     : nil",
		}
	}

	// the type of got is not reported when not-nil since the desired
	// error of type E may be in an expected wrapped error; the type of
	// got is not particularly helpful
	return report.ExpectedGot(expectedType, report.TypeName(got))
}
