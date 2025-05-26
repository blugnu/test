package test

import "github.com/blugnu/test/matchers/bytes"

// EqualBytes returns a matcher that checks if a byte slice is equal to an
// expected byte slice.
//
// The type T must be byte or a type that is assignable to byte. The matcher
// uses reflect.DeepEqual to compare the slices; this enables the matcher to
// be used with slices of custom byte types.
//
// The matcher reports differences in human-readable format.
//
// The failure report identifies the offsets of all differences and shows a
// portion of the expected and actual byte slices to highlight the first such
// difference.  See the example for a demonstration.
//
// # Supported Options
//
// This is a highly specialised matcher; the only supported option is a
// custom failure report function:
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the slices are not equal
func EqualBytes[T ~byte](want []T) *bytes.EqualMatcher[T] {
	return &bytes.EqualMatcher[T]{Expected: want}
}
