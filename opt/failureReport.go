package opt

// FailureReport is a function type that implements an OnTestFailure function.
//
// This option allows a test to override the default error report when a test
// fails, replacing it with a custom message.  All options supplied to the
// matcher are also passed to the function.  The custom function must accept
// the options as a variadic number of arguments, but may choose to ignore them.
//
//	Expect(got).To(BeTrue(), opt.FailureReport(func(...any) []string {
//		return []string{"custom failure message"}
//	}))
type FailureReport func(...any) []string

// TestFailure implements the TestFailure(...any) []string interface,
// calling the function with the provided options.
func (f FailureReport) OnTestFailure(opts ...any) []string {
	return f(opts...)
}

func OnFailure(msg string) FailureReport {
	return func(...any) []string {
		return []string{msg}
	}
}
