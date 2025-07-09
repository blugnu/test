package test

import (
	"fmt"

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
	Expect(false).To(BeTrue(), opt.OnFailure(msg))
}

// Errorf explicitly and unconditionally fails the current test
// with the formatted message.
//
// This should not be confused with the `test.Error` function
// used to report an error condition in a test helper (from the
// blugnu/test/test package).
func Errorf(s string, args ...any) {
	T().Helper()
	Expect(false).To(BeTrue(), opt.OnFailure(fmt.Sprintf(s, args...)))
}
