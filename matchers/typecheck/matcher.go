package typecheck

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
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
	var (
		expected     T
		expectedType = fmt.Sprintf("%T", expected)
	)

	// if we could not determine the type of the expected value using the
	// zero value of the type, we can use a dummy function and reflect
	// the type of the first argument to that function.
	//
	// Q: why not just use the dummy func technique every time?
	// A: because it is more expensive than using the zero value, and
	//    using the zero value provides a more precise (package-qualified)
	//    type name
	if expectedType == "<nil>" {
		fn := func(T) { /* NO-OP */ }
		fn(expected) // ensures that the dummy function is covered by tests

		expectedType = reflect.TypeOf(fn).In(0).Name()
	}

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{"should not be of type: " + expectedType}
	}

	return []string{
		"expected type: " + expectedType,
		fmt.Sprintf("got          : %T", got),
	}
}
