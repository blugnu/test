package test

import "github.com/blugnu/test/matchers/typecheck"

// BeOfType returns a matcher that checks if the value is of the expected type.
//
// To perform further tests on a value that passes this check using a strongly
// typed value of the expected type, use expect.Type or require.Type instead.
func BeOfType[T any]() typecheck.Matcher[T] {
	return typecheck.Matcher[T]{}
}
