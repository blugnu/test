package expect

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

// Type tests that a value is of an expected type T.  If the test passes,
// the value is returned as that type, with true. If the test fails the zero
// value of the specified type is returned, with false.
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test failure report
//	opt.Name(string)         // }-- equivalents --
//	opt.Namef(s, args)       // }
//
//	opt.FailReporter(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
//
//	opt.OnFailure(string)    // a simple string to output as the
//	                         // failure report if the test fails.
//
//	opt.IsRequired(bool)     // if true, the expectation is required to pass;
//	                         // no further expectations in the current test will be
//	                         // evaluated if the expectation fails.
//
//	opt.Required()           // equivalent to opt.IsRequired(true)
//
// # Alternatives
//
// If a test has no use for the returned T value, consider using the
// [test.BeOfType] matcher with [test.Expect], to avoid lint warnings about
// unused return values:
//
//	Expect(got).Should(BeOfType[T]())
//
// If the test can only continue when the value is of the expected type,
// consider using the [require.Type] function to halt test execution if
// when the value is not of the expected type:
//
//	v := require.Type[T](got)
func Type[T any](got any, opts ...any) (T, bool) {
	test.T().Helper()

	// the test passes if got asserts as type T
	if result, ok := got.(T); ok {
		return result, true
	}

	// otherwise, the test fails...

	// prepare default failure report identifying expected and actual types
	expectedType := report.TypeName[T]()
	gotType := report.TypeName(got)
	if got == nil {
		gotType = "nil"
	}

	// by separating and recombining options we ensure we have a canonical
	// set of expectation and matcher options
	opts, mopts := internal.SeparateOptions(opts)
	opts = append(opts, mopts...)

	// if no custom failure report is provided, add a default one that identifies
	// the expected and actual types
	if _, ok := opt.Get[opt.FailReporter](opts); !ok {
		opts = append(opts,
			opt.FailReporter(func(...any) []string {
				return []string{
					"expected type: " + expectedType,
					"got          : " + gotType,
				}
			}),
		)
	}

	test.Fail(opts...)

	// we return the zero value of T and false, in case the caller did not
	// specify opt.IsRequired(true)

	var zero T
	return zero, false
}
