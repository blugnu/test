package test

import (
	"errors"
	"testing"
)

// Panic is used to test for panics.
type Panic struct {
	error
}

// IsRecovered fails the test if the expected panic does not occur.
//
// The test will pass if:
//
//   - there is no expected panic and none has occurred
//   - there is a panic and the recovered value is an error
//     that satisfies errors.Is() with respect to the expected error
//
// Panics are tested by arranging an ExpectedPanic and then deferring
// a call to IsRecovered() in the test function.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer test.ExpectPanic(err).IsRecovered(t)
//
//		// ACT
//		doSomething()
//	  }
//
// IsRecovered may be called on a nil receiver and is equivalent to
// calling IsRecovered() on a *Panic with a nil error. This
// simplifies panic tests in data-driven tests where the expected
// panic may be nil for some test cases (indicating no panic is
// expected).
func (e *Panic) IsRecovered(t *testing.T) {
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

// ExpectPanic returns a *Panic that can be used to test that an expected
// panic occured and recovered a specified error.
//
// The same test can be used to verify that no panic occured, by calling
// the IsRecovered(t) method on a nil *Panic or passing nil as the error
// argument to ExpectPanic (which will then return a nil *Panic).
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		defer test.ExpectPanic(err).IsRecovered(t)
//
//		// ACT
//		doSomething()
//	  }
func ExpectPanic(err error) *Panic {
	if err == nil {
		return nil
	}
	return &Panic{err}
}
