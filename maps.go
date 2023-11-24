package test

import "testing"

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
//		test.Maps(t, want, got)
//	  }
func Maps[K comparable, V comparable](t *testing.T, want, got map[K]V) (keys []K, values []V) {
	t.Helper()

	ok := len(want) == len(got)
	if ok {
		wk := make([]K, 0, len(want))
		for k := range want {
			wk = append(wk, k)
		}
		for _, k := range wk {
			if ok = want[k] == got[k]; !ok {
				break
			}
		}
	}

	if !ok {
		t.Errorf("\nwanted: %#v\ngot    : %#v", want, got)
	}

	return
}
