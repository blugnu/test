package test

import (
	"testing"

	"github.com/blugnu/test/test"
)

// testRunner is a simple wrapper around a test function that allows it to
// be run as a subtest in the current test frame. It implements the Runnable
// interface, allowing it to be used with the Run() function.
type testRunner struct {
	name     string
	fn       func()
	parallel bool
}

// Run runs the named test function as a subtest in the current test frame
func (tr testRunner) Run() {
	t := T()
	t.Helper()

	if tr.parallel && IsParallel() {
		test.Invalid("ParallelTest() cannot be run from a parallel test")
	}

	t.Run(tr.name, func(t *testing.T) {
		With(t)
		t.Helper()

		if tr.parallel {
			t.Parallel()
		}

		tr.fn()
	})
}

// ParallelTest creates a test runner to run a function as a subtest
// with the provided name, running it in parallel.
//
// If the current test is already parallel, this function will
// fail the test as invalid since it is not allowed to nest parallel
// tests.
func ParallelTest(name string, fn func()) testRunner {
	T().Helper()

	if IsParallel() {
		test.Invalid("ParallelTest() cannot be run from a parallel test")
	}

	return testRunner{
		name:     name,
		fn:       fn,
		parallel: true,
	}
}

// Test creates a named test runner that can be used to run a test function as a
// subtest in the current test frame. The name is used as the subtest name,
// and the function is the test code to be executed.
func Test(name string, fn func()) testRunner {
	return testRunner{
		name: name,
		fn:   fn,
	}
}
