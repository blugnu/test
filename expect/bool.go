package expect

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
)

// False fails a test if a specified bool is not false.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(false))
//	Expect(got).To(BeFalse())
//
// # Common Options
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
func False[T ~bool](got T, opts ...any) {
	test.T().Helper()
	eopts, mopts := internal.SeparateOptions(opts)
	test.Expect(bool(got), eopts...).To(test.BeFalse(), mopts...)
}

// True fails a test if a specified bool is not true.  An optional
// name (string) may be specified to be included in the test report in the
// event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(true))
//	Expect(got).To(BeTrue())
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test failure report
//	opt.SubjectName(string)  // }
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
func True[T ~bool](got T, opts ...any) {
	test.T().Helper()
	eopts, mopts := internal.SeparateOptions(opts)
	test.Expect(bool(got), eopts...).To(test.BeTrue(), mopts...)
}
