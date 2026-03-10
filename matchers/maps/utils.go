package maps

import (
	"reflect"

	"strings"

	"github.com/blugnu/test/opt"
)

func as[T any](v any) T {
	result, _ := v.(T)
	return result
}

func containsMap[K comparable, V any](m, c map[K]V, opts ...any) bool {
	if len(c) == 0 {
		return len(m) == 0
	}

	eq := compareFuncFor[K, V](opts...)
	for k, v := range c {
		if vGot, exists := m[k]; !(exists && eq(v, vGot)) {
			return false
		}
	}

	return true
}

func compareFuncFor[K comparable, V any](opts ...any) func(a, b any) bool {
	type equatable interface{ Equal(V) bool }

	// check for Equatable or a comparison function in the options
	// to use

	var cmp func(a, b any) bool

	v := *new(V)
	if _, ok := any(v).(equatable); ok {
		cmp = func(a, b any) bool {
			return as[equatable](a).Equal(as[V](b))
		}
	} else if fn, ok := opt.Get[func(V, V) bool](opts); ok {
		cmp = func(a, b any) bool {
			return fn(as[V](a), as[V](b))
		}
	} else if fn, ok := opt.Get[func(any, any) bool](opts); ok {
		cmp = fn
	}

	// FUTURE: options should also be applied when comparing *values* (V) that are also maps.
	//
	// This may not be straightforward as the values in a map value cannot be assumed to be
	// of the same type as the map being tested so must be treated as a map[comparable]any
	// with functions equivalent to containsMap etc that operate on K:comparable, V:any

	switch reflect.ValueOf(v).Kind() { //nolint:exhaustive // exhaustive is not needed here
	case reflect.Slice, reflect.Array:
		if cmp == nil {
			cmp = func(a, b any) bool {
				return slicesEqual(
					reflect.ValueOf(a),
					reflect.ValueOf(b),
					opts...,
				)
			}
		}
	case reflect.String:
		// if the map values are strings and the CaseSensitive option is set to false
		// we wrap the comparison function to first normalize the strings to lower case
		// before comparing them
		if opt.IsSet(opts, opt.CaseInsensitive) {
			eq := cmp
			if eq == nil {
				eq = reflect.DeepEqual
			}
			cmp = func(a, b any) bool {
				a = strings.ToLower(as[string](a))
				b = strings.ToLower(as[string](b))
				return eq(a, b)
			}
		}
	}

	// if still nil then no comparison function was identified or provided so
	// fall-back on reflect.DeepEqual

	if cmp == nil {
		cmp = reflect.DeepEqual
	}

	return cmp
}

func slicesEqual(a, b reflect.Value, opts ...any) bool {
	switch {
	case a.Len() != b.Len():
		return false

	case a.IsNil():
		return b.IsNil()

	case b.IsNil():
		return false

	case a.Len() == 0:
		return true
	}

	cmp := reflect.DeepEqual
	if a.Index(0).Kind() == reflect.String && opt.IsSet(opts, opt.CaseInsensitive) {
		cmp = func(a, b any) bool {
			a = strings.ToLower(as[string](a))
			b = strings.ToLower(as[string](b))
			return a == b
		}
	}

	if !opt.IsSet(opts, opt.AnyOrder()) {
		for i := 0; i < a.Len(); i++ {
			if !cmp(a.Index(i).Interface(), b.Index(i).Interface()) {
				return false
			}
		}
		return true
	}

	return slicesEquivalent(a, b, cmp)
}

func slicesEquivalent(a, b reflect.Value, cmp func(a, b any) bool) bool {
	matched := make(map[int]struct{}, a.Len())
	for i := 0; i < a.Len(); i++ {
		a := a.Index(i).Interface()
		for j := 0; j < b.Len(); j++ {
			// if we already matched this element, skip it
			if _, ok := matched[j]; ok {
				continue
			}

			// if the elements are equal, mark j as matched
			if cmp(a, b.Index(j).Interface()) {
				matched[j] = struct{}{}
				break
			}
		}
	}

	return len(matched) == a.Len()
}
