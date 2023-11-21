package test

import (
	"reflect"
	"testing"
)

// Type fails the test if got is of type T.  The function returns
// the value of got as type T and a boolean indicating whether
// the test failed.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ACT
//		got := doSomething()
//
//		// ASSERT
//		if got, ok := test.Type[Foo](t, got); ok {
//			// apply further tests to got (of type Foo)
//		}
//	  }
func Type[T any](t *testing.T, got any) (T, bool) {
	z := *new(T)
	if reflect.TypeOf(got) != reflect.TypeOf(z) {
		t.Errorf("\nwanted: %T\ngot   : %T", z, got)
		return z, false
	}
	return got.(T), true
}
