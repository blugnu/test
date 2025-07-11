package test

import (
	"github.com/blugnu/test/internal/testcase"
	"github.com/blugnu/test/test"
)

// TestExecutor is an interface that defines a function to execute a test case
type TestExecutor[T any] interface {
	Execute(string, T)
}

// Case adds a test case to the runner. The name is used to identify the test
// case in the test output and any Before/Each scaffolding functions.  The tc
// provides the data for the test case.
//
// If the name is an empty or whitespace string, a name will be determined
// as follows:
//
//   - if the test case data is a struct (or pointer to a struct) and has a
//     string field named "Scenario", "scenario", "Name", or "name", with a
//     non-empty or whitespace value, that will be used as the name. Otherwise
//     a default name is derived in the format "testcase-NNN" where NNN is
//     the 1-based index of the test case in the list of test cases.
//
// If the test case data has a bool debug/Debug field that is set true, the
// test case will be marked as a debug test case.
//
// If the test case data has a bool skip/Skip field that is set true, the
// test case will be marked as a skip test case.
func Case[T any](name string, tc T) testcase.Registration[T] {
	return func(r *testcase.Runner[T], flags testcase.Flags) {
		r.AddCase(name, tc, flags)
	}
}

// Cases creates a Runner to run a set of test cases.
func Cases[T any](cases []T) testcase.Registration[T] {
	return func(r *testcase.Runner[T], flags testcase.Flags) {
		for _, c := range cases {
			r.AddCase("", c, flags)
		}
	}
}

// Debug adds a test case to the runner and marks it as a debug target.
//
// When running test cases, if any cases are marked as debug only those
// cases will be run. This is useful for debugging specific test cases
// without running the entire suite.
//
// Adding a test case using Debug() overrides any debug/skip fields that
// may be present in the test case data itself.
//
// When any debug test cases are run, the test case runner itself will fail
// the test with a warning that only debug cases were evaluated.  This is
// intended to ensure that Debug() cases are not accidentally left in the
// test suite.
func Debug[T any](name string, tc T) testcase.Registration[T] {
	return func(r *testcase.Runner[T], flags testcase.Flags) {
		r.AddCase(name, tc, flags|testcase.Debug)
	}
}

// ParallelCase adds a test case to the runner for parallel execution.
//
// Aside from the parallel execution of the test case subtests, this function
// is otherwise identical to Case().
func ParallelCase[T any](name string, tc T) testcase.Registration[T] {
	return func(r *testcase.Runner[T], flags testcase.Flags) {
		r.AddCase(name, tc, flags|testcase.Parallel)
	}
}

// Skip adds a test case to the runner and marks it to be skipped.
//
// Adding a test case using Skip() overrides any skip/debug fields that
// may be present in the test case data itself.
func Skip[T any](name string, tc T) testcase.Registration[T] {
	return func(r *testcase.Runner[T], flags testcase.Flags) {
		r.AddCase(name, tc, flags|testcase.Skip)
	}
}

// For creates a TestExecutor that uses the provided function to execute
// each test case. The function is called with the name of the test case and
// the test case data. This allows for variations in test execution based on
// the test case name, which can be useful for more complex test scenarios.
func For[T any](exec func(string, T)) TestExecutor[T] {
	if exec == nil {
		GetT().Helper()
		test.Invalid("For() function cannot be nil")
	}
	return testcase.NewExecutor[T](exec)
}

// ForEach creates a TestExecutor that uses the provided function to execute
// each test case. The function is called with the test case data, and the
// test case name is not provided.
func ForEach[T any](exec func(T)) TestExecutor[T] {
	if exec == nil {
		GetT().Helper()
		test.Invalid("ForEach() function cannot be nil")
	}
	return testcase.NewExecutor[T](exec)
}

// ParallelCases creates a Runner to run a set of test cases in parallel.
//
// Aside from the parallel execution of the test cases, this function
// is otherwise identical to Testcases().
func ParallelCases[T any](exec TestExecutor[T], cases ...testcase.Registration[T]) testcase.Runner[T] {
	if IsParallel() {
		GetT().Helper()
		test.Invalid("ParallelCases() cannot be run from a parallel test")
	}

	runner := testcase.NewRunner(exec)
	for _, reg := range cases {
		reg(&runner, testcase.Parallel)
	}

	return runner
}

// Testcases creates a Runner that can be used to run a set of test cases; each
// test case is performed in its own subtest using a test executor function
// that is provided as the first argument.
//
// # Test Executors
//
// For tests where the test executor is consistent for all test cases, the
// ForEach() function provides a single test executor function that requires
// only the test case itself.
//
// If test execution varies for some test cases, the For() function instead
// provides a test executor that will be called with the name of each test case
// in addition to the test case itself.  The test executor function can then
// use the test case name to perform any test case specifica variations as
// required.
//
// # Test Cases
//
// Test cases may be specified in two different ways:
//
//   - by calling the Case(), ParallelCase, Debug() or Skip() fluent methods
//     on the returned Runner; these functions require a name for each case;
//
//   - as a variadic list of test cases following the test executor;
//     these test cases are expected to be simple values or structs.  If a
//     struct type is used that has a name/Name/scenario/Scenario string field,
//     that field will be used as the name for the test case.  Otherwise a
//     default name will be generated in the format "testcase-NNN" where NNN
//     is the 1-based index of the test case in the list.
//
// The first form is recommended for cases where the name/scenario of each
// test case is significant, whether for varying test execution or simply
// in identifying a failed test case.
//
// The second form is useful for simple test cases where the test executor
// is the same for all test cases and the index-based identity for each test
// case is sufficient.
func Testcases[T any](exec TestExecutor[T], cases ...testcase.Registration[T]) testcase.Runner[T] {
	runner := testcase.NewRunner(exec)

	for _, reg := range cases {
		reg(&runner, testcase.NoFlags)
	}

	return runner
}
