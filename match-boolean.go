package test

import (
	"github.com/blugnu/test/matchers/bools"
)

// BeFalse returns a matcher that will fail if the matched value is not false.
//
// # Supported Options
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
func BeFalse() bools.BooleanMatcher {
	return bools.BooleanMatcher{Expected: false}
}

// BeTrue returns a matcher that will fail if the matched value is not true.
//
// # Supported Options
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
func BeTrue() bools.BooleanMatcher {
	return bools.BooleanMatcher{Expected: true}
}
