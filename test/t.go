package test

import (
	"github.com/blugnu/test/internal/testframe"
)

// Helper is an interface that defines the method set required to report
// make a frame as part of a test helper call stack
type Helper interface {
	Helper()
}

// T returns the current test frame's helper interface, which is used to
// mark the test as a helper.
//
// The function will panic if the current test frame does not contain a
// valid test frame or if the test frame does not implement the Helper
// interface.
//
// It is recommended to use this to mark functions as helpers before
// reporting an invalid test or error:
//
//	func MyHelperFunction() {
//	   // ... evaluate pre-conditions for helper ...
//	   if !preCondition {
//	      test.T().Helper()
//	      test.Invalid("pre-conditions not met")
//	   }
//	   // ... continue with helper logic ...
//	}
func T() Helper {
	if t, ok := testframe.Peek[Helper](); ok {
		return t
	}
	return noopHelper{}
}
