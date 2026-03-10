package test

import (
	"github.com/blugnu/test/matchers/length"
)

// HaveLen returns a matcher that checks if the value has len() equal to n.
//
// The returned matcher is an [AnyMatcher] that must be used with
// [expectation.Should] or [expectation.ShouldNot] expectations.  The matcher
// supports values of any type that is compatible with the built-in len()
// function.  That is:
//
//   - array
//   - channel
//   - map
//   - slice
//   - string
//
// A nil value of any of these types is considered to have a length of 0.
//
// If the value is of any other type, the test fails as invalid, with a
// message similar to:
//
//	`length.Matcher: requires a value that is a string, slice, channel, or map: got <type>`
func HaveLen(n int) *length.Matcher {
	return &length.Matcher{Length: n}
}
