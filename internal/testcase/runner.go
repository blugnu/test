package testcase

import (
	"fmt"
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/test"
)

type Flags int

const (
	// NoFlags indicates that no flags are set for the test case
	NoFlags Flags = 0

	// Parallel indicates that the test case should be run in parallel
	Parallel Flags = 1

	// Debug indicates that the test case is a debug target
	Debug Flags = 2

	// Skip indicates that the test case should be skipped
	Skip Flags = 4
)

// TestExecutor[T] is an interface that defines a method to execute an individual
// test case of type T.
type TestExecutor[T any] interface {
	Execute(name string, tc T)
}

// NamedCaseExecution[T] is a function type that takes a test case name and a pointer
// to a test case of type T.
type NamedCaseExecution[T any] func(name string, tc T)

// AnonCaseExecution[T] is a function type that takes a pointer to a test case of type T
// without a name.
type AnonCaseExecution[T any] func(tc T)

// Runner[T] is a struct that holds a testing.T interface, a test executor,
// a setup controller, and a list of case controllers, used for executing a number
// of test cases of type T.
//
// It provides methods to configure setup functions and add and run test cases.
type Runner[T any] struct {
	TestingT
	TestExecutor[T]

	CaseControllers []Controller[T]
}

type Registration[T any] func(*Runner[T], Flags)

// NewRunner creates a new Runner instance with the provided test executor
func NewRunner[T any](exec TestExecutor[T]) Runner[T] {
	t := testframe.MustPeek[TestingT]()

	if exec == nil {
		t.Helper()
		test.Invalid("test executor cannot be nil")
	}

	return Runner[T]{
		TestingT:     t,
		TestExecutor: exec,
	}
}

// AddCase adds a test case to the runner. The name is used to identify the test
// case in the test output.  The tc provides the data for the test case.
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
func (tcr *Runner[T]) AddCase(name string, tc T, flags ...Flags) {
	idx := len(tcr.CaseControllers)
	ctrl := NewController(tc, idx+1, name)
	tcr.CaseControllers = append(tcr.CaseControllers, ctrl)

	f := Flags(0)
	for _, flag := range flags {
		f |= flag
	}

	ctrl.debug = ctrl.debug || (f&Debug != 0)
	ctrl.skip = (ctrl.skip && (f&Debug == 0)) || (f&Skip != 0)
	ctrl.parallel = ctrl.parallel || (f&Parallel != 0)

	tcr.CaseControllers[idx] = ctrl // update the controller in the list
}

// Run runs the test cases in the runner. Each test case is run as a subtest
// in the current test frame. The test case name is used to identify the test
// case in the test output, and any Before/After scaffolding functions are
// called with the test case data.
func (tcr Runner[T]) Run() {
	t := tcr.TestingT
	t.Helper()

	if len(tcr.CaseControllers) == 0 {
		test.Invalid("no test cases provided")
		return
	}

	// initially assume that all test cases are runnable
	runnable := tcr.CaseControllers

	// if there are any debug test cases these are the only ones that will be run
	debugging := tcr.getDebugCases()
	if len(debugging) > 0 {
		runnable = debugging
	}

	nSkipped := 0
	for _, tc := range runnable {
		name := tc.name
		tcr.TestingT.Run(name, func(t *testing.T) {
			testframe.Push(t)
			t.Helper()

			if tc.skip {
				nSkipped++
				t.SkipNow()
			}

			if tc.parallel {
				t.Parallel()
			}

			tc := tc.data // copy the test case data
			tcr.TestExecutor.Execute(name, tc)
		})
	}

	tcr.doWarnings(len(debugging), nSkipped)
}

// doWarnings reports any warnings that should be issued after running the
// test cases.
func (tcr Runner[T]) doWarnings(nDebugged, nSkipped int) {
	tcr.Helper()

	switch nDebugged {
	case 0:
		// no debug cases, nothing to report
	case len(tcr.CaseControllers):
		// a bit odd, but this means that all cases were evaluated
		// so nothing to report
	default:
		test.Warning(fmt.Sprintf("only %d of %d cases were evaluated (debug mode)", nDebugged, len(tcr.CaseControllers)))
	}

	switch nSkipped {
	case 0:
		// no cases were skipped, nothing to report
	case len(tcr.CaseControllers):
		test.Warning("all cases were skipped")
	default:
		test.Warning(fmt.Sprintf("%d of %d cases were skipped", nSkipped, len(tcr.CaseControllers)))
	}
}

func (tcr Runner[T]) getDebugCases() []Controller[T] {
	result := make([]Controller[T], 0, len(tcr.CaseControllers))

	for _, tc := range tcr.CaseControllers {
		if tc.debug {
			result = append(result, tc)
		}
	}

	return result
}
