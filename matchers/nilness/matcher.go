package nilness

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type Matcher struct{}

// IsNil checks that the value of the expectation is nil.  If the
// value is not nil, the test fails.  If the value is nil, the test
// passes.
//
// If the value being tested does not support a nil value the test
// will fail and produce a report similar to:
//
//	nilnessMatcher: values of type '<type>' are not nilable
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
func (m Matcher) Match(subject any, opts ...any) bool {
	if result, isNilable := isNil(subject); isNilable {
		return result
	}

	// if it is expected that the value being tested is NOT nil, the
	// fact that the value is not nilable is not a failure.
	//
	// It may be a pointless test, but we cannot assume that.
	//
	// It may also be a test where the value being tested is an 'any'
	// value which may hold either a nilable or non-nilable value
	// but where in any event it is expected that the value is not nil.
	//
	// otherwise it is expected that the value being tested IS nil, in which
	// case if the value is of non-nilable type, then the test cannot be valid
	if !opt.IsSet(opts, opt.ToNotMatch(true)) {
		test.T().Helper()
		test.Invalid(
			fmt.Sprintf("nilness.Matcher: values of type '%T' are not nilable", subject),
		)
	}

	return false
}

func (m Matcher) OnTestFailure(subject any, opts ...any) string {
	expect := "expected nil"
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		expect = "expected not nil"
	}

	switch got := any(subject).(type) {
	case error:
		return fmt.Sprintf("%s, got error: %v", expect, got)

	default:
		if reflect.ValueOf(got).Kind() == reflect.Pointer {
			v := reflect.Indirect(reflect.ValueOf(got))
			if v.Kind() == reflect.String {
				return fmt.Sprintf("%s, got &(%s)", expect, opt.ValueAsString(v.String(), opts...))
			}
			return fmt.Sprintf("%s, got &(%#v)", expect, v)
		}
		return fmt.Sprintf("%s, got %s", expect, opt.ValueAsString(got, opts...))
	}
}
