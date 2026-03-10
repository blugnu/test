package require

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// False fails and halts execution of the current test if a specified bool
// is not false.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(false), opt.Required())
//	Expect(got).To(BeFalse(), opt.Required())
//	expect.False(got, opt.Required())
//	Require(got).To(Equal(false))
//	Require(got).To(BeFalse())
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
func False[T ~bool](got T, opts ...any) {
	test.T().Helper()
	expect.False(got, append(opts, opt.Required())...)
}

// True fails and halts execution of the current test if a specified bool
// is not true.
//
// This test is a convenience for these equivalent alternatives:
//
//	Expect(got).To(Equal(true), opt.Required())
//	Expect(got).To(BeTrue(), opt.Required())
//	expect.True(got, opt.Required())
//	Require(got).To(Equal(true))
//	Require(got).To(BeTrue())
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
func True[T ~bool](got T, opts ...any) {
	test.T().Helper()
	expect.True(got, append(opts, opt.Required())...)
}
