package test

import "testing"

// Equal fails the test if got is not equal to wanted.  Wanted and got
// must be of the same type which may be any type that is comparable:
//
// An optional Format value may be provided to specify the format of
// values displayed in test failure reports. If not specified,
// FormatDefault is used (%v).
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		a := 1
//		b := 2
//
//		// ACT
//		c := a + b
//
//		// ASSERT
//		test.Equal(t, 3, c)
//	  }
func Equal[T comparable](t *testing.T, wanted, got T, opt ...Format) {
	t.Helper()
	if wanted != got {
		f := FormatDefault
		if len(opt) > 0 {
			f = opt[0]
		}
		t.Errorf("\nwanted: %s\ngot   : %s", format(wanted, f), format(got, f))
	}
}
