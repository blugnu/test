package test

import (
	"fmt"
	"reflect"
	"testing"
)

// provides methods for testing a map.
type MapTest[K comparable, V any] struct {
	testable[map[K]V]
}

// returns a MapTest for testing a specified map.
func Map[K comparable, V any](t *testing.T, got map[K]V, opts ...any) MapTest[K, V] {
	t.Helper()

	n := "map"
	f := FormatDefault
	ffn := *new(func(map[K]V) string)
	checkOptTypes(t, optTypes(n, f, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	getOpt(&ffn, append(opts, func(v map[K]V) string { return fmt.Sprintf(string(f), v) })...)

	return MapTest[K, V]{newTestable(t, got, n, f, ffn)}
}

// compares two maps and fails the test if they are not equal.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		want := map[string]int{"a": 1, "b": 2}
//		got := map[string]int{"a": 1, "b": 2}
//
//		// ASSERT
//		test.Map(t, "got", got).Equals(want)
//	  }
func (m MapTest[K, V]) Equals(wanted map[K]V) {
	m.Helper()

	ok := len(m.got) == len(wanted)
	if ok {
		for k, got := range m.got {
			if ok = reflect.DeepEqual(got, wanted[k]); !ok {
				break
			}
		}
	}

	if !ok {
		m.fail(m.T, wanted)
	}
}
