package expect

import (
	"errors"
	"fmt"
	"slices"

	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"

	"github.com/blugnu/test"
)

// Error asserts than an error is non-nil or satisfies `errors.Is` for a
// specified target error.
//
// # Supported Options
//
//	error                    // a target error
//	                         //
//	                         // When specified, the test will pass if the got error is
//	                         // non-nil and satisfies errors.Is for the target error
//	                         //
//	                         // When not specified, the test will pass if the got error
//	                         // is non-nil
//
//	string                   // a name for the value, for use in any test failure report
//	opt.SubjectName(string)  // }
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
// If a test can only continue when the error is non-nil or matches a target
// error, consider using the [require] package function to halt test execution
// when the error does not match:
//
//	require.Error(err, opts...) // opts may include a target error
func Error(got error, opts ...any) {
	test.T().Helper()

	if target, isSet := opt.Get[error](opts); isSet {
		test.Expect(got).Is(target, opts...)
		return
	}

	opts = slices.Clone(opts)
	opts = append(opts, opt.OnFailure("expected error"))

	test.Expect(got).IsNotNil(opts...)
}

// ErrorAs tests that an error is of an expected type E, using [errors.As].
// An ok indicator is returned indicating whether the test passed (true)
// or failed (false).  If the test passes, the error is returned as a
// value of type E, otherwise the zero value of E is returned and should
// be ignored.
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test failure report
//	opt.SubjectName(string)  // }
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
// If a test has no use for the returned E value, consider using the
// [test.BeError] matcher with [test.Expect], to avoid lint warnings about
// unused return values:
//
//	Expect(err).Should(BeError[E]())
//
// If the test can only continue when the error is of the expected type,
// consider using the require package ErrorAs function to halt test execution
// when the error is not of the expected type:
//
//	e := require.ErrorAs[E](err)
func ErrorAs[E error](got error, opts ...any) (E, bool) {
	test.T().Helper()

	// the test passes if the error is non-nil and satisfies errors.As(E)
	var target E
	if got != nil && errors.As(got, &target) {
		return target, true
	}

	// otherwise, the test fails...

	// by separating and recombining options we ensure we have a canonical
	// set of expectation and matcher options
	opts, mopts := internal.SeparateOptions(opts)
	opts = append(opts, mopts...)

	// if no custom failure report is provided, add a default one
	// that identifies the expected type
	if _, ok := opt.Get[opt.FailReporter](opts); !ok {
		opts = append(opts,
			opt.FailReporter(func(...any) []string {
				expectedType := internal.TypeName[E]()
				if internal.IsNil(got) {
					return []string{
						"expected error: " + expectedType,
						"got           : nil",
					}
				}
				return []string{
					"expected error to match or wrap: " + expectedType,
					"got: " + report.Value(got, opts...),
				}
			}),
		)
	}

	test.Fail(opts...)

	return target, false
}

// NoError asserts that an error is nil.
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test failure report
//	opt.SubjectName(string)  // }
//	opt.Name(string)         // }-- equivalents --
//	opt.Namef(s, args)       // }
//
//	opt.FailReporter(func)   // a function returning a custom failure report
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
func NoError(got error, opts ...any) {
	test.T().Helper()

	opts = slices.Clone(opts)
	opts = append(opts, opt.FailReporter(func(a ...any) []string {
		return []string{
			"unexpected error",
			fmt.Sprintf("got: %[1]q [%[1]T]", got),
		}
	}))

	test.Expect(got).IsNil(opts...)
}
