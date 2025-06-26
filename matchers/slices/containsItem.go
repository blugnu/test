package slices

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/blugnu/test/opt"
)

// ContainsItemsMatcher is a matcher for []T that will match the []T
// if it contains at least one element that is equal to the expected
// element.
type ContainsItemMatcher[T any] struct {
	Expected T
}

// Match satisfies the matcher interface for []T.
func (m ContainsItemMatcher[T]) Match(got []T, opts ...any) bool {
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

	// if the values being compared are strings and opt.CaseSensitive(false) is set
	// then we normalise strings to lowercase before comparing, making sure to
	// preserve the use of any custom comparison function
	if reflect.ValueOf(m.Expected).Kind() == reflect.String &&
		opt.IsSet(opts, opt.CaseSensitive(false)) {
		og := cmp
		cmp = func(a, b any) bool {
			a = strings.ToLower(fmt.Sprintf("%v", a))
			b = strings.ToLower(fmt.Sprintf("%v", b))
			return og(a, b)
		}
	}

	return slice[T](got).next(m.Expected, 0, cmp) != -1
}

func (m ContainsItemMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	const minlen = 1
	var maxlen = 2 + len(got)

	result := make([]string, minlen, maxlen)
	switch opt.IsSet(opts, opt.ToNotMatch(true)) {
	case true:
		result[0] = fmt.Sprintf("expected: %T not containing: "+opt.ValueAsString(m.Expected, opts...), got)
	default:
		result[0] = fmt.Sprintf("expected: %T containing: "+opt.ValueAsString(m.Expected, opts...), got)
	}

	result = slice[T](got).appendToTestReport(result, "got:", opts...)

	if reflect.ValueOf(m.Expected).Kind() == reflect.String &&
		opt.IsSet(opts, opt.CaseSensitive(false)) {
		result = append(result, "(case insensitive comparison)")
	}
	return result
}
