package typecheck

import (
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type Matcher[T any] struct{}

// Match checks that the value is of the expected type.  If the
// value is not of the expected type, the test fails.  If the
// value is of the expected type, the test passes.
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
func (m Matcher[T]) Match(got any, opts ...any) bool {
	_, ok := got.(T)
	return ok
}

func (m Matcher[T]) OnTestFailure(got any, opts ...any) []string {
	expectedType := report.TypeName[T]()

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{"should not be of type: " + expectedType}
	}

	return []string{
		"expected type: " + expectedType,
		"got          : " + report.TypeName(got),
	}
}
