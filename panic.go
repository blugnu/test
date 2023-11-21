package test

import (
	"errors"
	"testing"
)

// ExpectedPanic is used to test for panics.
type ExpectedPanic struct {
	error
}

// Assert fails the test if the expected panic does not occur.
//
// The test will pass if:
//
//   - there is no expected panic and none has occurred
//   - there is a panic and the recovered value is an error
//     that satisfies errors.Is() with respect to the expected error
//
// Panics are tested by arranging an ExpectedPanic and then deferring
// a call to Assert() in the test function.
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
// Assert may be called on a nil receiver and is equivalent to
// calling Assert() on an ExpectedPanic with a nil error. This
// simplifies panic tests in data-driven tests where the expected
// panic may be nil for some test cases (indicating no panic is
// expected).
func (e *ExpectedPanic) Assert(t *testing.T) {
	t.Helper()

	r := recover()

	switch {
	case e == nil && r == nil:
		return
	case e == nil && r != nil:
		t.Errorf("\nunexpected panic: %v", r)
	case e != nil && r == nil:
		t.Errorf("\nwanted (panic): %v\ngot           : (did not panic)", e.error)
	case e != nil && (e.error) != nil && r != nil:
		if got, ok := r.(error); !ok || !errors.Is(got, e.error) {
			t.Errorf("\nwanted (panic): %v\ngot    (panic): %s", e.error, r)
		}
	}
}

// ExpectPanic returns an ExpectedPanic that can be used to test for panics.
func ExpectPanic(err error) *ExpectedPanic {
	return &ExpectedPanic{err}
}
