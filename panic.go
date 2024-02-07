package test

import (
	"errors"
	"testing"
)

// provides methods for testing panics.
type PanicTest struct {
	r any
}

// fails the test if the expected panic does not occur.
//
// The test will pass if:
//
//   - there is no expected panic and none has occurred
//   - there is a panic and the recovered value is an error
//     that satisfies errors.Is() with respect to the expected error
//
// Panics are tested by arranging an ExpectedPanic and then deferring a call to
// Assert() in the test function.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer test.ExpectPanic(err).Assert(t)
//
//		// ACT
//		doSomething()
//	  }
//
// Assert may be called on a nil receiver and is equivalent to calling Assert() on
// a *Panic with a nil error. This simplifies panic tests in data-driven tests where
// the expected panic may be nil for some test cases (indicating no panic is expected).
func (e *PanicTest) Assert(t *testing.T) {
	t.Helper()

	err := error(nil)
	iserr := false
	if e != nil && e.r != nil {
		err, iserr = e.r.(error)
	}

	r := recover()

	switch {
	case e == nil && r == nil:
		return
	case e == nil && r != nil:
		t.Errorf("\nunexpected panic: %[1]T: %[1]v", r)
	case e != nil && r == nil:
		t.Errorf("\nwanted: panic: %[1]T: %[1]v\ngot   : (did not panic)", e.r)
	case iserr && r != nil:
		if got, ok := r.(error); !ok || !errors.Is(got, err) {
			t.Errorf("\nwanted: panic: %[1]T: %[1]v\ngot   : panic: %[2]T: %[2]v", err, r)
		}
	default:
		if e.r != r {
			t.Errorf("\nwanted: panic: %[1]T: %[1]v\ngot   : panic: %[2]T: %[2]v", e.r, r)
		}
	}
}

// returns a *Panic that can be used to test that an expected panic is
// recovered with a specified error.
//
// A *Panic is also be used to verify that no panic occured, by calling
// the Assert(t) method on a nil *Panic or passing nil as the
// expected recovery value (in which case the function will return a
// nil *Panic).
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer test.ExpectPanic(err).Assert(t)
//
//		// ACT
//		doSomething()
//	  }
func ExpectPanic(recover any) *PanicTest {
	if recover == nil {
		return nil
	}
	return &PanicTest{recover}
}
