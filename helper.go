package test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type HelperResult bool

const ShouldPass HelperResult = true
const ShouldFail HelperResult = false

// runs a test helper function and compares the result to an
// expected outcome.  The expected outcome (want) must be ONE of
// the following:
//
//   - test.ShouldPass
//   - test.ShouldFail
//   - *test.Panic (see: test.ExpectPanic())
//
// The function returns `stdout` and `stderr` output captured
// while the helper function is executed.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		var a, b int
//
//		// ACT
//		a = 1
//		b = 2
//
//		// ASSERT
//		test.Helper(t, "test case 1", func(st *testing.T) {
//		   test.Compare(st, 3, a+b)
//		}, test.ShouldPass)
//	  }
func Helper(t *testing.T, f func(st *testing.T), want any) (CapturedOutput, CapturedOutput) {
	t.Helper()

	checkresult := true
	result := true
	err := error(nil)

	switch expected := want.(type) {
	case bool:
		result = expected
	case HelperResult:
		result = bool(expected)
	case *Panic:
		err = expected.error
	case nil: // no outcome specified
		checkresult = false
	default:
		panic(fmt.Errorf("%w: accepts only test.ShouldPass (or true), test.ShouldFail (or false) or test.ExpectPanic(error)", ErrInvalidArgument))
	}

	oc := map[bool]string{true: "PASS", false: "FAIL"}

	stdout, stderr, _ := runInternal(t, func(st *testing.T) {
		defer func() {
			r := recover()
			got, ok := r.(error)
			switch {
			case err == nil && r == nil:
				return
			case checkresult && err == nil && r != nil:
				t.Errorf("\nwanted     : %s\ngot (panic): %v", oc[result], r)
				return
			case checkresult && err != nil && r == nil:
				t.Errorf("\nwanted (panic): %v\ngot           : %s", err, oc[!st.Failed()])
				return
			case checkresult && !ok || !errors.Is(got, err):
				t.Errorf("\nwanted (panic): %v\ngot    (panic): %v", err, r)
			}
		}()

		st.Helper()
		f(st)

		if checkresult && result != !st.Failed() {
			t.Errorf("\nwanted: %s\ngot   : %s", oc[result], oc[!st.Failed()])
		}
	})
	return stdout, stderr
}

// matchAll is a test matcher function that always returns true and no error.
// It is used by runInternal to ensure that the internal test is always
// executed by the InternalTest runner, regardless of any match conditions
// specified for the current *testing.T test run.
func matchAll(pat, match string) (bool, error) {
	return true, nil
}

// runInternal is a helper function that runs a specified test function as an
// internal test.
//
// It is used by test.Helper() to run a test function in a separate test runner
// in order to inspect the outcome of the test function without that affecting
// the state of the current test.
func runInternal(t *testing.T, f func(*testing.T)) (CapturedOutput, CapturedOutput, bool) {
	result := false
	stdout, stderr := CaptureOutput(t, func(t *testing.T) {
		t.Helper()
		sut := []testing.InternalTest{{
			Name: "internal",
			F:    f,
		}}
		result = testing.RunTests(matchAll, sut)
	})

	// this is a bit of a hack: when run with -v, the internal test run
	// emits a RUN line and a PASS line when the test passes, which we don't
	// want to see in the captured output.  So we remove them here.
	//
	// (if the internal test fails there is no additional output to worry about)

	for i := len(stdout.s) - 1; i >= 0; i-- {
		s := strings.TrimSpace(stdout.s[i])
		if s == "=== RUN   internal" || strings.HasPrefix(s, "--- PASS: internal") {
			stdout.s = append(stdout.s[:i], stdout.s[i+1:]...)
		}
	}

	return stdout, stderr, result
}
