package equal

import (
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

// MARK: DeepEqualMatcher

type DeepMatcher[T any] struct {
	Expected T
}

func (m DeepMatcher[T]) valueAsString(v any, opts ...any) string {
	if v == nil {
		return "nil"
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		// FUTURE: improve the formatting of structs
		// to be more readable...
		//
		// encoding/json MarshalIndent is nice, but limite to exported fields
		// and quotes field names (yuck)
		//
		// magicjson is better but not in the standard library and still suffers
		// from the problem of quoting field names (rules out JSON in general)
		//
		// For now, just use %#v to get the type and value; the ideal would be
		// to report the type name separately and format the result differently:
		//
		//    test.foo{name:"arthur", age:42}
		//
		// -> test.foo {
		//      name: "arthur",
		//      age : 42
		//    }
		return fmt.Sprintf("%#v", v)
	default:
		return opt.ValueAsString(v, opts...)
	}
}

func (m DeepMatcher[T]) Match(got T, opts ...any) bool {
	if equable, ok := any(m.Expected).(interface{ Equal(T) bool }); ok {
		return equable.Equal(got)
	}

	if cmp, ok := opt.Get[func(T, T) bool](opts); ok {
		return cmp(m.Expected, got)
	}

	if cmp, ok := opt.Get[func(any, any) bool](opts); ok {
		return cmp(m.Expected, got)
	}

	return reflect.DeepEqual(m.Expected, got)
}

func (m DeepMatcher[T]) OnTestFailure(got T, opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{
			"expected to not equal: " + m.valueAsString(got, opts...),
		}
	}

	ef := m.valueAsString(m.Expected, opts...)
	gf := m.valueAsString(got, opts...)

	// if the expected and got values are small, use a one line report
	if len(ef) < 10 && len(gf) < 10 {
		return []string{fmt.Sprintf("expected %s, got %s", ef, gf)}
	}

	// otherwise, use a multi-line report
	return []string{
		fmt.Sprintf("expected: %s", ef),
		fmt.Sprintf("got     : %s", gf),
	}
}
