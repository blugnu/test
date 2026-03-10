package expect

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
)

// Equal fails a test if two values are not equal.
//
// Supported Options:
//
//	string                   // a name for the value, for use in any test failure report
//	opt.Name(string)         // }-- equivalents --
//	opt.Namef(s, args)       // }
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
//
//	opt.OnFailure(string)    // a simple string to output as the
//	                         // failure report if the test fails.
//
//	opt.Required()           // the expectation is required to pass; no further
//	                         // expectations in the current test will be evaluated if
//	                         // the expectation fails
//
//	opt.IsRequired(true)      // equivalent to opt.Required()
func Equal[T comparable](actual, expected T, opts ...any) {
	test.T().Helper()
	eopts, mopts := internal.SeparateOptions(opts)
	test.Expect(actual, eopts...).To(test.Equal(expected), mopts...)
}
