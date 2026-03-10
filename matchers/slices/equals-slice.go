package slices

import (
	"reflect"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

// EqualMatcher is a matcher for []T that will match the []T
// if the two slices are the same length and contain the same items.
//
// By default, the items in the slice must be in the same order as
// the items in the expected slice.  This may be overridden by
// specifying the ExactOrder(false) option in which case the order
// of the items in the slice is not significant.
type EqualMatcher[T any] struct {
	Expected []T
}

// Match satisfies the matcher interface for []T.
func (m EqualMatcher[T]) Match(got []T, opts ...any) bool {
	switch {
	case len(got) == 0 && len(m.Expected) == 0:
		return true
	case len(got) != len(m.Expected):
		return false
	}

	cmp := reflect.DeepEqual
	if fn, ok := opt.Get[func(T, T) bool](opts); ok {
		cmp = func(a, b any) bool {
			at, _ := a.(T)
			bt, _ := b.(T)
			return fn(at, bt)
		}
	} else if fn, ok := opt.Get[func(any, any) bool](opts); ok {
		cmp = fn
	}

	if opt.IsSet(opts, opt.ExactOrder(false)) {
		return slice[T](got).containsItems(m.Expected, cmp)
	}
	return slice[T](got).containsSlice(m.Expected, cmp)
}

// OnTestFailure returns a report of the failure for the matcher.
func (m EqualMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	result := make([]string, 0, 2+len(got)+len(m.Expected))
	cond := "equal to"
	inv := opt.IsSet(opts, opt.ToNotMatch(true))
	if inv {
		cond = "not equal to"
	}

	result = report.AppendSlice(result, slice[T](m.Expected), opt.WithNamef(opts, "expected %T %s:", got, cond)...)
	if inv {
		return result
	}

	return report.AppendSlice(result, slice[T](got), opt.Force(opts, opt.Name("got:"))...)
}
