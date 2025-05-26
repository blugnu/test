package test

import (
	"fmt"
	"reflect"
	"testing"
)

// isParallel checks if the given TestingT is running in parallel.
//
// The TestingG must be a *testing.T; any other implementation will return
// false.
//
// For *testing.T values, reflection is used to examine the internal state
// to determine if it is running in parallel.
func isParallel(t TestingT) bool {
	if _, ok := T().(*testing.T); !ok {
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

// Parallel marks a test for possible parallel execution.
//
// Parallel must not be called from a test that is already parallel; doing so
// will result in the test failing as invalid.
//
// Parallel must be called with 0 or 1 test runner arguments.  If called with
// 0 arguments, it marks the current test frame for parallel execution.
//
// If called with 1 argument, it marks the test runner for parallel execution,
// as an alternative to calling With(t) and then t.Parallel().
//
// If called with more than 1 argument, the function will fail any current
// test frame as invalid.  If there is no current test frame, the function
// will panic.
func Parallel(t ...TestingT) {
	parallel := func(t TestingT) {
		if isParallel(t) {
			t.Helper()
			invalidTest("Parallel() must not be called from a parallel test")
			return
		}

		t.Parallel()
	}

	switch len(t) {
	case 0:
		T().Helper()
		parallel(T())

	case 1:
		t[0].Helper()
		With(t[0])
		parallel(t[0])

	default:
		if t := hasT(); t != nil {
			t.Helper()
			invalidTest("Parallel() must be called with 0 or 1 test runner arguments")
			return
		}
		panic(fmt.Errorf("%w: Parallel() must be called with 0 or 1 test runner arguments", ErrInvalidArgument))
	}
}

// RunParallel runs a named sub-test in parallel.
//
// RunParallel must not be called from a test that is already parallel; doing so
// will result in the test failing as invalid.
func RunParallel(n string, fn func()) {
	t := GetT()
	t.Helper()

	if IsParallel() {
		invalidTest("RunParallel() must not be called from a parallel test")
		return
	}

	t.Run(n, func(st *testing.T) {
		Parallel(st)
		fn()
	})
}

// RunParallelScenarios runs a set of test cases in parallel.
//
// RunParallelScenarios must not be called from a test that is already parallel; doing so
// will result in the test failing as invalid.
func RunParallelScenarios[T any](f func(tc *T, i int), scns []T) {
	GetT().Helper()

	if IsParallel() {
		invalidTest("RunParallelScenarios() must not be called from a parallel test")
		return
	}

	RunScenarios(
		func(tc *T, i int) {
			t := GetT()
			t.Helper()
			t.Parallel()

			f(tc, i)
		},
		scns,
	)
}
