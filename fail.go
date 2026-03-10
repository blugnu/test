package test

import (
	"fmt"
	"strings"

	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
)

// Error explicitly and unconditionally fails the current test
// with the given message.
//
// This should not be confused with the `test.Error` function
// used to report an error condition in a test helper (from the
// blugnu/test/test package).
func Error(msg string) {
	T().Helper()
	Fail(opt.OnFailure(msg))
}

// Errorf explicitly and unconditionally fails the current test
// with the formatted message.
//
// This should not be confused with the `test.Error` function
// used to report an error condition in a test helper (from the
// blugnu/test/test package).
func Errorf(s string, args ...any) {
	T().Helper()
	Fail(opt.OnFailure(fmt.Sprintf(s, args...)))
}

// Fail explicitly and unconditionally fails the current test with the given
// options.
//
// If the supplied options do not contain a string, []string or [opt.FailReporter],
// the test is failed with a "test failed" message.
func Fail(opts ...any) {
	getT(opts).Helper()
	FailIf(true, opts...)
}

// FailIf fails the current test with the given options only if a specified
// condition is true.
//
// If the supplied options do not contain a string, []string or [opt.FailReporter],
// the test is failed with a "test failed" message.
func FailIf(cond bool, opts ...any) {
	if !cond {
		return
	}

	T().Helper()

	if _, ok := opt.Get[opt.FailReporter](opts); !ok {
		var msg string = "test failed"

		if s, ok := opt.Get[[]string](opts); ok {
			msg = strings.Join(s, "\n")
		} else if s, ok := opt.Get[string](opts); ok {
			msg = s
		}

		opts = append([]any{opt.OnFailure(msg)}, opts...)
	}

	// options may include both expectation and matcher options
	eopts, mopts := internal.SeparateOptions(opts)

	Expect(cond, eopts...).To(BeFalse(), mopts...)
}

// FailIfNot fails the current test with the given options only if a specified
// condition is not true.
//
// If the supplied options do not contain a string, []string or [opt.FailReporter],
// the test is failed with a "test failed" message.
func FailIfNot(cond bool, opts ...any) {
	T().Helper()
	FailIf(!cond, opts...)
}
