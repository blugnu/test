package test

import (
	"cmp"

	"github.com/blugnu/test/matchers/ordered"
)

// BeBetween returns a matcher that will fail if the matched value is not
// between the limits of a defined interval. The type T must be ordered
// (i.e. it must support the comparison operators <, <=, >, >=).
//
// BeBetween accepts a single value as an argument providing one limit of
// the interval. The second limit is provided by calling the [And] method
// on the returned initializer:
//
//	Expect(n).To(BeBetween(10).And(20))
//
// By default the matcher compares the matched value to the closed interval
// defined by the specified values.
//
// # Supported Options
//
// The interval can be changed by providing an [opt.IntervalClosure] option.
// The default interval is [opt.IntervalClosed].  Other options are:
//
//   - [opt.IntervalOpen]    for an open interval (min < x < max)
//   - [opt.IntervalOpenMin] for a half-open interval (min < x <= max)
//   - [opt.IntervalOpenMax] for a half-open interval (min <= x < max)
func BeBetween[T cmp.Ordered](v T) ordered.IsBetweenInitializer[T] {
	return ordered.IsBetweenInitializer[T]{Value: v}
}

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
// To compare values using the <= operator, call the [OrEqual] modifier
// method on the returned matcher:
//
//	Expect(n).To(BeGreaterThan(10).OrEqual())
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
func BeGreaterThan[T cmp.Ordered](want T) ordered.RelativeMatcher[T] {
	return ordered.RelativeMatcher[T]{
		Expected:   want,
		Comparison: ordered.GreaterThan,
	}
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
// To compare values using the <= operator, call the [OrEqual] modifier
// method on the returned matcher:
//
//	Expect(n).To(BeLessThan(10).OrEqual())
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
func BeLessThan[T cmp.Ordered](want T) ordered.RelativeMatcher[T] {
	return ordered.RelativeMatcher[T]{
		Expected:   want,
		Comparison: ordered.LessThan,
	}
}
