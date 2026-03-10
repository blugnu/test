package require

import (
	"github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

// Nil fails and halts the current test if a specified value is not nil.
// An optional name (string) may be specified to be included in the test
// report in the event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Require(got).IsNil()
//	Require(got).Should(BeNil())
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test
//	                         // failure report
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
func Nil(got any, opts ...any) {
	test.T().Helper()
	expect.Nil(got, append(opts, opt.Required())...)
}

// NotNil fails and halts the current test if a specified value is nil.
// An optional name (string) may be specified to be included in the test
// report in the event of failure.
//
// This test is a convenience for these equivalent alternatives:
//
//	Require(got).IsNot(nil)
//	Require(got).ShouldNot(BeNil())
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test
//	                         // failure report
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
func NotNil(got any, opts ...any) {
	test.T().Helper()
	expect.NotNil(got, append(opts, opt.Required())...)
}
