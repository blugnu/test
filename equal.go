package test

import "github.com/blugnu/test/matchers/equal"

// Equal returns a matcher that checks if a value of type T is equal to some
// expected value of type T. The type T must be comparable.
//
// Equality is determined using ONE of the following, in order of preference:
//
//  1. any comparison function provided as an option
//  2. a T.Equal(T) bool method, if it exists
//  3. the == operator
//
// If specified, a custom comparison function must take two arguments of type
// T, returning a boolean indicating whether the values are considered equal.
//
// # Supported Options
//
//	func(T, T) bool            // a function to compare the values
//	                           // (overriding the use of the == operator or
//	                           // T.Equal(T) method, if it exists)
//
//	opt.QuotedString(bool)     // determines whether string values in failure reports
//	                           // are quoted (default is true); the option has no
//	                           // effect on values that are not strings
//	                           //
//	                           // this option applies only to values being directly
//	                           // compared; it is not applied to string fields of
//	                           // struct types, for example
//
//	opt.FailureReport(func)    // a function returning a custom failure report
//	                           // when the values are not equal
func Equal[T comparable](want T) equal.Matcher[T] {
	return equal.Matcher[T]{Expected: want}
}

// DeepEqual returns a matcher that checks if a value of type T is equal
// to some expected value of type T.
//
// Equality is always evaluated using reflect.DeepEqual.  i.e. the matcher does
// not support the use of a custom comparison function and will not use any
// Equal(T) method implemented by type T.
//
// To use a custom comparison function or a T.Equal(T) method, use the Equal()
// matcher.
//
// # Supported Options
//
//	opt.QuotedString(bool)     // determines whether string values in failure reports
//	                           // are quoted (default is true); the option has no
//	                           // effect on values that are not strings
//	                           //
//	                           // this option applies only to values being directly
//	                           // compared; it is not applied to string fields of
//	                           // struct types, for example
//
//	opt.FailureReport(func)    // a function returning a custom failure report
//	                           // when the values are not equal
func DeepEqual[T any](want T) equal.DeepMatcher[T] {
	return equal.DeepMatcher[T]{Expected: want}
}
