package test

import (
	"strings"
	"testing"
)

type TestOutcome bool

const ShouldPass TestOutcome = true
const ShouldFail TestOutcome = false
const TestPassed = ShouldPass
const TestFailed = ShouldFail

type R struct {
	t         TestRunner
	recovered any
	Report    []string
	Log       []string
	Outcome   TestOutcome
}

func (ht R) HadOutcome(wanted TestOutcome) {
	ht.t.Helper()
	ht.t.Run("outcome", func(t *testing.T) {
		t.Helper()
		s := map[TestOutcome]string{
			ShouldPass: "PASS",
			ShouldFail: "FAIL",
		}
		if ht.Outcome != wanted {
			t.Errorf("\nwanted: %s\ngot   : %s", s[wanted], s[ht.Outcome])
		}
	})
}

func (ht R) DidFail() {
	ht.t.Helper()
	ht.t.Run("did fail", func(t *testing.T) {
		t.Helper()
		if ht.Outcome != TestFailed {
			t.Errorf("\nwanted: FAIL\ngot   : PASS")
		}
	})
}

func (ht R) DidNotPanic() {
	ht.t.Helper()
	ht.t.Run("did not panic", func(t *testing.T) {
		t.Helper()
		if ht.recovered != nil {
			t.Errorf("\nwanted: (did not panic)\ngot   : panic: %v", ht.recovered)
		}
	})
}

func (ht R) DidPanic(wanted any) {
	ht.t.Helper()
	ht.t.Run("did panic", func(t *testing.T) {
		t.Helper()
		if ht.recovered == nil {
			t.Errorf("\nwanted: panic: %[1]T: %[1]v\ngot   : (did not panic)", wanted)
		} else if ht.recovered != wanted {
			t.Errorf("\nwanted: panic: %[1]T: %[1]v\ngot   : panic: %[2]T: %[2]v", wanted, ht.recovered)
		}
	})
}

func (ht R) DidPass() {
	ht.t.Helper()
	ht.t.Run("did pass", func(t *testing.T) {
		t.Helper()
		if ht.Outcome != TestPassed {
			t.Errorf("\nwanted: PASS\ngot   : FAIL")
		}
	})
}

func (ht R) Assert(r ...any) {
	ht.t.Helper()

	var expected []string
	for _, v := range r {
		switch s := v.(type) {
		case string:
			expected = append(expected, s)
		case []string:
			expected = append(expected, s...)
		}
	}

	if shouldPass := !hasOpt(r, ShouldFail) && (len(r) == 0 || hasOpt(r, ShouldPass) || len(expected) == 0); shouldPass {
		Expect(ht.Outcome, "test should pass").To(Equal(TestPassed))
		Expect(ht.Report).IsEmpty()
		return
	}

	Expect(ht.Outcome, "test should fail").To(Equal(TestFailed))
	Expect(ht.Report).To(ContainStrings(expected))
}

// runs a test helper function and compares the result to an
// expected outcome.  The expected outcome (want) must be ONE of
// the following:
//
//   - test.ShouldPass / test.Passed / true
//   - test.ShouldFail / test.Failed / false
//   - *test.Panic (see: test.ExpectPanic())
//   - error
//   - string / []string
//
// If an error is specified, ExpectPanic(error) is implied.
//
// If a string or []string is specified then test.ShouldFail is
// implied.  If test.ShouldFail, test.Failed or false are specified
// without a string or []string then the test output is not checked.
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
//		stdout, stderr := test.Test(t, "test case 1", func(st *testing.T) {
//		   test.Compare(st, 3, a+b)
//		}, test.ShouldPass)
//	  }
func Test(f func()) R {
	t := T().(*testing.T)
	t.Helper()

	var recovered any
	stdout, stderr, outcome := runInternal(t, func(st *testing.T) {
		st.Helper()
		// defer func() {
		// 	st.Helper()
		// 	if recovered = recover(); recovered != nil {
		// 		st.Errorf("\nunexpected panic: %[1]T: %[1]v", recovered)
		// 	}
		// }()

		With(st)
		f()
	})

	return R{
		t:         t,
		recovered: recovered,
		Log:       stderr,
		Report:    stdout,
		Outcome:   TestOutcome(outcome),
	}
}

// func Helper(t *testing.T, f func(st *testing.T), want ...any) (result HelperTest) {
// 	t.Helper()

// 	result.t = t
// 	result.testOutcome = func(t *testing.T, outcome TestResult) { /* NO-OP */ }
// 	result.testPanic = func(t *testing.T, recovered any) { /* NO-OP */ }
// 	result.testReport = func(t *testing.T, report StringsTest) { /* NO-OP */ }

// 	expect := struct {
// 		result TestResult
// 		panic  *Panic
// 	}{
// 		result: ShouldPass,
// 		panic:  nil,
// 	}

// 	assert := len(want) > 0
// 	for _, want := range want {
// 		switch wanted := want.(type) {
// 		case TestResult:
// 			expect.result = wanted
// 			result.testOutcome = func(t *testing.T, outcome TestResult) { t.Helper(); result.HadOutcome(wanted) }
// 		case bool:
// 			expect.result = TestResult(wanted)
// 			result.testOutcome = func(t *testing.T, outcome TestResult) { t.Helper(); result.HadOutcome(TestResult(wanted)) }
// 		case *Panic:
// 			expect.result = ShouldFail
// 			expect.panic = wanted
// 			result.testPanic = func(t *testing.T, recovered any) { t.Helper(); result.DidPanic(wanted.r) }
// 		case error:
// 			expect.result = ShouldFail
// 			expect.panic = ExpectPanic(wanted)
// 			result.testPanic = func(t *testing.T, recovered any) { t.Helper(); result.DidPanic(wanted) }
// 		case string:
// 			if len(wanted) == 0 {
// 				t.Errorf("invalid expected outcome option: (empty string): did you mean test.ShouldPass?")
// 				return
// 			}

// 			expect.result = ShouldFail
// 			result.testOutcome = func(t *testing.T, outcome TestResult) { t.Helper(); result.DidFail() }
// 			result.testReport = func(t *testing.T, report StringsTest) {
// 				t.Helper()
// 				report.Contains(wanted)
// 			}
// 		case []string:
// 			if len(wanted) == 0 {
// 				t.Errorf("invalid expected outcome option: []string{}: did you mean test.ShouldPass?")
// 				return
// 			}
// 			expect.result = ShouldFail
// 			result.testOutcome = func(t *testing.T, outcome TestResult) { t.Helper(); result.DidFail() }
// 			result.testReport = func(t *testing.T, report StringsTest) {
// 				t.Helper()
// 				report.Contains(wanted)
// 			}

// 		default:
// 			t.Errorf("invalid option: %[1]T: %[1]v: expected test.TestResult/bool, string, []string or *Panic", want)
// 		}
// 	}

// 	stdout, _, outcome := runInternal(t, func(st *testing.T) {
// 		st.Helper()
// 		defer func() {
// 			result.recovered = recover()
// 			result.result = Failed
// 			if assert {
// 				expect.panic.assert(t, result.recovered)
// 			} else if result.recovered != nil {
// 				st.Errorf("\nunexpected panic: %[1]T: %[1]v", result.recovered)
// 			}
// 		}()

// 		f(st)
// 	})
// 	result.result = TestResult(outcome)
// 	result.Report = stdout

// 	if assert {
// 		result.testOutcome(t, expect.result)
// 		result.testPanic(t, nil)
// 		result.testReport(t, stdout)
// 	}

// 	return result
// }

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
func runInternal(t *testing.T, f func(*testing.T)) ([]string, []string, bool) {
	t.Helper()

	result := false
	stdout, stderr := CaptureOutput(t, func() {
		t.Helper()
		sut := []testing.InternalTest{{
			Name: t.Name(),
			F:    f,
		}}
		result = testing.RunTests(matchAll, sut)
	})

	// this is a bit of a hack: when run with -v, the internal test run
	// emits a RUN line and a PASS line when the test passes, which we don't
	// want to see in the captured output.  So we remove them here.
	//
	// (if the internal test fails there is no additional output to worry about)
	for i := len(stdout) - 1; i >= 0; i-- {
		s := strings.TrimSpace(stdout[i])
		if strings.HasPrefix(s, "=== RUN") || strings.HasPrefix(s, "--- PASS: ") {
			stdout = append(stdout[:i], stdout[i+1:]...)
		}
	}

	return stdout, stderr, result
}
