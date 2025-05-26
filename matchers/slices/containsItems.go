package slices

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

// ContainsItemsMatcher is a matcher for []T that will match the []T
// if it contains all of the elements in another slice.
//
// The items in the expected slice may occur in any order in the got
// slice and need not be contiguous.
type ContainsItemsMatcher[T any] struct {
	Expected []T
}

// Match satisfies the matcher interface for []T.
func (m ContainsItemsMatcher[T]) Match(got []T, opts ...any) bool {
	cmp := reflect.DeepEqual
	if fn, ok := opt.Get[func(T, T) bool](opts); ok {
		cmp = func(a, b any) bool {
			// this is a type-safe matcher; we can safely cast without checking
			return fn(a.(T), b.(T))
		}
	} else if fn, ok := opt.Get[func(any, any) bool](opts); ok {
		cmp = fn
	}

	return slice[T](got).containsItems(m.Expected, cmp)
}

func (m ContainsItemsMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	result := make([]string, 0, 2+len(got)+len(m.Expected))
	cond := "containing items"
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not containing items"
	}

	result = slice[T](m.Expected).appendToTestReport(result, fmt.Sprintf("expected: %T %s:", got, cond), opts...)
	return slice[T](got).appendToTestReport(result, "got:", opts...)
}
