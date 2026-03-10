package test

import (
	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

// IsParallel returns true if the current test is running in parallel or is a
// sub-test of a parallel test.
func IsParallel() bool {
	return testframe.IsParallel()
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
	// mark t as a helper if possible
	if t, ok := testframe.Peek[TestingT](); ok {
		t.Helper()
	}

	switch {
	case t == nil:
		test.Invalid("Parallel() cannot be called with nil")

	case testframe.IsParallel():
		test.Invalid("Parallel() must not be called from a parallel test")
	}

	With(t)
	t.Parallel()
}
