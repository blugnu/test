package test

import (
	"fmt"
	"reflect"
)

// MARK: EqualMatcher

type EqualMatcher[T comparable] struct {
	Expecting[T]
}

func (m EqualMatcher[T]) Format(v any) string {
	return DeepEqualMatcher[T](m).Format(v)
}

func (m EqualMatcher[T]) Match(got T, _ ...any) bool {
	return m.Expecting.value == got
}

func Equal[T comparable](e T) EqualMatcher[T] {
	return EqualMatcher[T]{Expecting[T]{e}}
}

// MARK: DeepEqualMatcher

type DeepEqualMatcher[T any] struct {
	Expecting[T]
}

func (m DeepEqualMatcher[T]) Format(v any) string {
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
		switch v := any(v).(type) {
		case string:
			return fmt.Sprintf("%q", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
}

func (m DeepEqualMatcher[T]) Match(actual T, _ ...any) bool {
	return reflect.DeepEqual(m.Expecting.value, actual)
}

func DeepEqual[T any](e T) DeepEqualMatcher[T] {
	return DeepEqualMatcher[T]{Expecting[T]{e}}
}
