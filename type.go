package test

import (
	"reflect"
	"testing"
)

// tests that a specified value is of a required type.  The test
// fails if the value is not of type T.
//
// The function returns the value of got as type T (if possible) and
// a boolean indicating whether the test passed (true) or failed (false).
// If the test fails then the returned value of type T is a zero-value,
// unrelated to the value being tested.
//
// If the test fails a test report is written similar to:
//
//	wanted: string
//	got   : int
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ACT
//		got := doSomething()
//
//		// ASSERT
//		if got, ok := test.IsType[Foo](t, got); ok {
//			// further tests performed on shadowed got (of type Foo)
//		}
//	  }
func IsType[T any](t *testing.T, got any) (T, bool) {
	t.Helper()

	z := *new(T)
	if reflect.TypeOf(got) != reflect.TypeOf(z) {
		t.Run("is of type", func(t *testing.T) {
			t.Helper()
			t.Errorf("\nwanted: %T\ngot   : %T", z, got)
		})
		return z, false
	}
	return got.(T), true
}
