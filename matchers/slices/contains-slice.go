package slices

import (
	"reflect"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

// ContainsSliceMatcher is a matcher for []T that will match the []T
// if it contains all of the elements in another slice.
//
// The items in the slice must be in the same order as the items in
// the expected slice.
type ContainsSliceMatcher[T any] struct {
	Expected []T
}

// Match satisfies the matcher interface for []T.
func (m ContainsSliceMatcher[T]) Match(got []T, opts ...any) bool {
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

	return slice[T](got).containsSlice(m.Expected, cmp)
}

// OnTestFailure returns a report of the failure for the matcher.
func (m ContainsSliceMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	result := make([]string, 0, 2+len(got)+len(m.Expected))
	cond := "containing slice"
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		cond = "not containing slice"
	}

	result = report.AppendSlice(result, m.Expected, opt.WithNamef(opts, "expected %T %s:", got, cond)...)
	return report.AppendSlice(result, got, opt.Force(opts, opt.Name("got:"))...)
}
