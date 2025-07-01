package test

import (
	"github.com/blugnu/test/matchers/slices"
)

// ContainItem returns a matcher for slices of comparable types that is
// satisfied if the slice contains at least one item that is equal to the
// expected item.
//
// The Matcher supports the following options:
//
//   - func(T, T) bool
//   - func(any, any) bool
//
// The options allow a test to supply a custom comparison function to be
// used to compare the expected item with the items in the slice. The
// comparison function must return true if the two arguments are equal.
//
// If both types of functions are provided the matcher will panic with
// ErrInvalidOption.
//
// The default comparison function is reflect.DeepEqual.
func ContainItem[T any](item T) slices.ContainsItemMatcher[T] {
	return slices.ContainsItemMatcher[T]{Expected: item}
}

// ContainItems returns a matcher for slices of comparable types that is
// satisfied if the slice contains all of the items in the expected
// slice.  The items in the expected slice may occur in any order in the
// got slice and need not be contiguous.
//
// The the following options are supported:
//
//   - func(T, T) bool
//   - func(any, any) bool
//
// The options allow a test to supply a custom comparison function to be
// used to compare the expected item with the items in the slice. The
// comparison function must return true if the two arguments are equal.
//
// If both types of functions are provided the matcher will panic with
// ErrInvalidOption.
//
// The default comparison function is reflect.DeepEqual.
func ContainItems[T any](items []T) slices.ContainsItemsMatcher[T] {
	return slices.ContainsItemsMatcher[T]{Expected: items}
}

// ContainSlice returns a matcher for a slice that is satisfied if the
// slice contains items that correspond to the items in a given other
// slice.
//
// Comparison of items is performed using reflect.DeepEqual by default
// but may be overridden by supplying one of:
//
//	func(T, T) bool
//	func(any, any) bool
func ContainSlice[T comparable](e []T) slices.ContainsSliceMatcher[T] {
	return slices.ContainsSliceMatcher[T]{Expected: e}
}

// EqualSlice compares the actual slice with an expected slice and fails
// the test if they are not equal.
//
// By default, the order of elements in each slice is significant.  That
// is, the nth each slice must be equal. If the order of elements is not
// significant, use the ExactOrder option to specify that the order of
// elements is not significant, e.g.:
//
//	got := []int{1, 2, 3}
//	expected := []int{3, 2, 1}
//	Expect(got).To(EqualSlice(expected))                         // will fail
//	Expect(got).To(EqualSlice(expected), opt.ExactOrder(false))  // will pass
func EqualSlice[T comparable](e []T) slices.EqualMatcher[T] {
	return slices.EqualMatcher[T]{Expected: e}
}
