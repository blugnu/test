package typecheck

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
	"github.com/blugnu/test/test"
)

type MatchIsError struct {
	Target error
}

// Match checks that the value satisfies [errors.Is] for the specified
// Target error.
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
func (m MatchIsError) Match(got error, _ ...any) bool {
	if m.Target == nil {
		return got == nil
	}

	return errors.Is(got, m.Target)
}

// OnTestFailure returns a failure report indicating that the actual value
// was not of the expected type.
//
// If the ToNotMatch option is set, the report indicates that the value
// should not have been of the expected type.
func (m MatchIsError) OnTestFailure(got error, opts ...any) []string {
	if m.Target == nil {
		test.T().Helper()
		test.Invalid("target error must be non-nil")
	}

	result := func(expected, got string) []string {
		return []string{
			"expected error: " + expected,
			"got           : " + got,
		}
	}

	var (
		expectedString = report.Value(m.Target.Error(), opts...)
		expectedType   = fmt.Sprintf("%T", m.Target)
		expectedErr    = expectedString + " [" + expectedType + "]"
	)

	switch {
	case got == nil:
		return result(expectedErr, "nil")

	case opt.IsSet(opts, opt.ToNotMatch(true)):
		return []string{"should not be: " + expectedErr}

	default:
		var (
			gotString = report.Value(got.Error(), opts...)
			gotType   = fmt.Sprintf("%T", got)
		)

		if gotType == expectedType {
			return result(expectedErr, gotString)
		}

		// pad the expected and got strings so that they are the same width in the report
		// and report both with type information to make it clear that the error types
		// are different

		maxLen := max(len(gotString), len(expectedString))
		spec := "%-" + strconv.Itoa(maxLen) + "s [%s]"

		return result(
			fmt.Sprintf(spec, expectedString, expectedType),
			fmt.Sprintf(spec, gotString, gotType),
		)
	}
}
