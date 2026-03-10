package require

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// Equal asserts that the actual value is equal to the expected value. If the
// assertion fails, the test will be halted.
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
//	opt.IsRequired(true)     // the expectation is required to pass; no further
//	                         // expectations in the current test will be evaluated if
//	                         // the expectation fails
//	opt.Required()           // equivalent to opt.IsRequired(true)
func Equal[T comparable](expected, actual T, opts ...any) {
	test.T().Helper()
	expect.Equal(expected, actual, opt.Force(opts, opt.Required())...)
}
