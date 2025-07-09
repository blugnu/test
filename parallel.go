package test

import (
	"reflect"
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

// isParallel checks if the given TestingT is running in parallel.
//
// The underlying type of TestingT must be a *testing.T; any other
// implementation will return false.
//
// For *testing.T values, reflection is used to examine the internal
// state to determine if it is running in parallel.
func isParallel(t TestingT) bool {
	tf := T()
	if _, ok := tf.(*testing.T); !ok {
		// tests cannot be parallel if the TestingT is not a *testing.T
		return false
	}

	c := reflect.Indirect(reflect.ValueOf(t)).FieldByName("common")
	ip := reflect.Indirect(c).FieldByName("isParallel")

	return ip.Bool()
}

// IsParallel returns true if the current test is running in parallel or is a
// sub-test of a parallel test.
func IsParallel() bool {
	return isParallel(T())
}

// Parallel establishes a new test frame scheduled for parallel execution.
// It is intended to be used as an alternative to With(t) for a test that
// is intended to run entirely in parallel.
//
// i.e. use:
//
//	func TestSomething(t *testing.T) {
//	  Parallel(t)
//	  // ... test code here ...
//	}
//
// instead of:
//
//	func TestSomething(t *testing.T) {
//	  With(t)
//
//	  Run(ParallelTest("something", func() {
//	     // ... test code here ...
//	  }))
//	}
//
// Parallel must not be called from a test that is already parallel or with
// a nil argument; in both cases the test will be failed as invalid.
func Parallel(t TestingT) {
	if t == nil {
		if t, ok := testframe.Peek[TestingT](); ok {
			t.Helper()
		}
		test.Invalid("Parallel() cannot be called with nil")
	}

	t.Helper()

	if isParallel(t) {
		test.Invalid("Parallel() cannot be called from a parallel test")
	}

	With(t)
	t.Parallel()
}
