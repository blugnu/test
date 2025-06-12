package test

import "github.com/blugnu/test/test"

type Runner[T any] struct {
	BeforeEach func(*T, int)
	AfterEach  func(*T, int)
}

// RunParallel is used for running table-driven tests using a function and
// a slice of scenarios where each scenario is run in parallel.
func (r Runner[T]) RunParallel(fn func(*T, int), scns []T) {
	t := GetT()
	t.Helper()

	if IsParallel() {
		test.Invalid("Runner.RunParallel() must not be called from a parallel test")
		return
	}

	RunParallelScenarios(
		func(tc *T, i int) {
			t := GetT()
			t.Helper()

			exec := func(fn func(tc *T, i int)) {
				if fn == nil {
					return
				}
				t.Helper()
				fn(tc, i)
			}

			exec(r.BeforeEach)
			exec(fn)
			exec(r.AfterEach)
		},
		scns,
	)
}

// RunScenarios is used for running table-driven tests using a function and
// a slice of scenarios. Each scenario is run in its own test runner.
//
// The test name for each scenario is determined as follows:
//  1. If the scenario is a struct and has a string field named "Scenario",
//     "scenario", "Name" or "name", that field is used as the name for the test;
//  2. If the scenario is not a struct or does not have a supported name field,
//     the test is named "testcase/000" where 000 is the index of the scenario
//
// If multiple candidate test name fields are found, the first one is used in
// the order of preference as follows: Scenario > scenario > Name > name.
func (r Runner[T]) RunScenarios(fn func(*T, int), scns []T) {
	GetT().Helper()

	RunScenarios(
		func(tc *T, i int) {
			t := GetT()
			t.Helper()

			exec := func(fn func(tc *T, i int), n ...string) {
				if fn == nil {
					return
				}
				if len(n) == 0 {
					fn(tc, i)
					return
				}
				Run(n[0], func() {
					GetT().Helper()
					fn(tc, i)
				})
			}

			exec(r.BeforeEach, "BeforeEach")
			exec(fn)
			exec(r.AfterEach, "AfterEach")
		},
		scns,
	)
}
