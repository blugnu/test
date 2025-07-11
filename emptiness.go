package test

import (
	"github.com/blugnu/test/matchers/emptiness"
	"github.com/blugnu/test/matchers/length"
)

// BeEmpty returns a matcher that checks if the value is empty.
//
// The returned matcher is an `AnyMatcher` that may only be used
// with the `Should()` method, or with `To()` where the subject
// is of formal type any.
//
// NOTE: A nil value is not considered empty and will fail this test.
// To test for an empty value that may be nil, use BeEmptyOrNil() instead.
//
// i.e. an empty slice will pass this test, but a nil slice will not.
//
// This test may be used to check for empty strings, arrays, slices,
// channels, maps and any type that implement a Count(), Len() or
// Length() function returning int, int64, uint, or uint64.
//
// If a value of any other type is tested, the test fails with a message
// similar to:
//
//	emptiness.Matcher: requires a value of type string, array, slice, channel or map,
//	                   or a type that implements a Count(), Len(), or Length() function
//	                   returning int, int64, uint, or uint64.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
//
//	func(v T) bool              // a function that returns true if the value
//	                            // is empty; the value argument must be of
//	                            // the same type as the expectation subject.
func BeEmpty() *emptiness.Matcher {
	return &emptiness.Matcher{}
}

// BeEmptyOrNil returns a matcher that checks if a value is empty or nil.
//
// # Compatible Methods and Subjects
//
//	Expect(any(subject)).To(...)       // i.e. where 'subject' is of type `any`
//	Expect(subject).Should(...)        // for any 'subject'
//
// This matcher may be used to check for empty strings and arrays as well as
// empty or nil slices, channels, maps and any type that implements a Count(),
// Len() or Length() method returning int, int64, uint, or uint64.
//
// If the subject is of any other type is tested, the test fails with a message
// similar to:
//
//	emptiness.Matcher: requires a value of type string, array, slice, channel or map,
//	                   or a type that implements a Count(), Len(), or Length() function
//	                   returning int, int64, uint, or uint64.
//
// # Supported Options
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
//
//	func(v T) bool              // a function that returns true if the value
//	                            // is empty or nil; the value argument must be of
//	                            // the same type as the subject.
func BeEmptyOrNil() *emptiness.Matcher {
	return &emptiness.Matcher{TreatNilAsEmpty: true}
}

// HaveLen returns a matcher that checks if the value has len() equal to n.
//
// The returned matcher is an `AnyMatcher` that may only be used with values
// of a type that is compatible with the built-in len() function.  That is:
//
// - string
// - slice
// - array
// - channel
// - map
//
// A nil value of any of these types is considered to have a length of 0.
//
// If the value is of any other type, the test fails as an invalid test,
// with a message similar to:
//
//	length.Matcher: requires a value that is a string, slice, channel, or map: got <type>
func HaveLen(n int) *length.Matcher {
	return &length.Matcher{Length: n}
}
