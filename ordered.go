package test

import (
	"cmp"

	"github.com/blugnu/test/matchers/ordered"
)

// BeGreaterThan returns a matcher that will fail if the matched value is not
// greater than the expected value. The type T must be ordered (i.e. it must
// support the comparison operators <, <=, >, >=).
//
// By default the matcher uses the > operator to compare the values. This
// can be overridden by providing a custom comparison function as an
// option. The function must take two arguments of type T and return a
// boolean indicating whether the first argument is greater than the second
// argument.
//
// # Supported Options
//
//	func(T, T) bool          // a custom comparison function to compare values
//	                         // (overriding the use of the > operator)
//
//	opt.QuotedStrings(bool)  // determines whether string values are quoted
//	                         // in test failure report (quoted by default);
//	                         // the option has no effect if the value is not
//	                         // a string type
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
//
//	opt.OnFailure(string)    // a string to output as the failure
//	                         // report if the test fails
func BeGreaterThan[T cmp.Ordered](want T) ordered.IsGreaterThan[T] {
	return ordered.IsGreaterThan[T]{Expected: want}
}

// BeLessThan returns a matcher that will fail if the matched value is not
// less than the expected value. The type T must be ordered (i.e. it must
// support the comparison operators <, <=, >, >=).
//
// By default the matcher uses the < operator to compare the values. This
// can be overridden by providing a custom comparison function as an
// option. The function must take two arguments of type T and return a
// boolean indicating whether the first argument is less than the second
// argument.
//
// # Supported Options
//
//	func(T, T) bool          // a custom comparison function to compare values
//	                         // (overriding the use of the < operator)
//
//	opt.QuotedStrings(bool)  // determines whether string values are quoted
//	                         // in test failure report (quoted by default);
//	                         // the option has no effect if the value is not
//	                         // a string type
//
//	opt.FailureReport(...)   // a function returning a custom failure report
//	                         // in the event that the test fails
//
//	opt.OnFailure(string)    // a string to output as the failure
//	                         // report if the test fails
func BeLessThan[T cmp.Ordered](want T) ordered.IsLessThan[T] {
	return ordered.IsLessThan[T]{Expected: want}
}
