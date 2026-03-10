package expect

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
)

// Nil fails a test if a specified value is not nil.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).Is(nil)
//	Expect(got).Should(BeNil())
//
// # Supported Options
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
//	opt.IsRequired(bool)     // if true, the expectation is required to pass;
//	                         // no further expectations in the current test will be
//	                         // evaluated if the expectation fails.
//
//	opt.Required()           // equivalent to opt.IsRequired(true)
func Nil(got any, opts ...any) {
	test.T().Helper()
	eopts, mopts := internal.SeparateOptions(opts)
	test.Expect(got, eopts...).Should(test.BeNil(), mopts...)
}

// NotNil fails a test if a specified value is nil.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).IsNot(nil)
//	Expect(got).ShouldNot(BeNil())
//
// # Supported Options
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
//	opt.IsRequired(bool)     // if true, the expectation is required to pass;
//	                         // no further expectations in the current test will be
//	                         // evaluated if the expectation fails.
//
//	opt.Required()           // equivalent to opt.IsRequired(true)
func NotNil(got any, opts ...any) {
	test.T().Helper()
	eopts, mopts := internal.SeparateOptions(opts)
	test.Expect(got, eopts...).ShouldNot(test.BeNil(), mopts...)
}
