package test

import (
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

// testRunner implements [Runnable] to execute a function with an associated
// test name.
type testRunner struct {
	name     string
	fn       func()
	parallel bool
}

func newTestRunner(name string, fn func(), parallel bool) testRunner {
	return testRunner{name: name, fn: fn, parallel: parallel}
}

type Helper interface {
	Helper()
}

// Run runs the named test function as a subtest in the current test frame
func (tr testRunner) Run() {
	type T interface {
		Helper()
		Run(string, func(*testing.T)) bool
	}
	t := testframe.MustPeek[T]()
	t.Helper()

	if tr.parallel && testframe.IsParallel() {
		test.Invalid("ParallelTest() cannot be run from a parallel test")
	}

	t.Run(tr.name, func(t *testing.T) {
		t.Helper()
		testframe.PushWithCleanup(t)

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
	if testframe.IsParallel() {
		T().Helper()
		test.Invalid("ParallelTest() cannot be run from a parallel test")
	}

	return newTestRunner(name, fn, true)
}

// Test creates a named test runner that can be used to run a test function as a
// subtest in the current test frame. The name is used as the subtest name,
// and the function is the test code to be executed.
func Test(name string, fn func()) testRunner {
	return newTestRunner(name, fn, false)
}
