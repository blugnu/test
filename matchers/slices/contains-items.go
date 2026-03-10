package slices

import (
	"reflect"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
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
			at, _ := a.(T)
			bt, _ := b.(T)
			return fn(at, bt)
		}
	} else if fn, ok := opt.Get[func(any, any) bool](opts); ok {
		cmp = fn
	}

	if opt.IsSet(opts, opt.AnyItem) {
		return slice[T](got).containsAny(m.Expected, cmp)
	}

	return slice[T](got).containsItems(m.Expected, cmp)
}

func (m ContainsItemsMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	set := "items"
	if opt.IsSet(opts, opt.AnyItem) {
		set = "any of"
	}

	result := make([]string, 0, 2+len(got)+len(m.Expected))
	cond := "containing " + set
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not containing " + set
	}

	result = report.AppendSlice(result, m.Expected, opt.WithNamef(opts, "expected %T %s:", got, cond)...)
	return report.AppendSlice(result, got, opt.Force(opts, opt.Name("got:"))...)
}
