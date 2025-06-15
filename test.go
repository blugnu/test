package test

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/matchers/slices"
	"github.com/blugnu/test/opt"
)

type TestOutcome int

const (
	TestPassed TestOutcome = iota
	TestFailed
	TestPanicked
)

func (to TestOutcome) String() string {
	switch to {
	case TestPassed:
		return "TestPassed"
	case TestFailed:
		return "TestFailed"
	case TestPanicked:
		return "TestPanicked"
	default:
		return fmt.Sprintf("TestOutcome(%d)", to)
	}
}

// R is a struct that contains the result of executing a test function.
type R struct {
	// value recovered if the test execution caused a panic
	Recovered any

	// captured stdout output from the test function (the test failure report)
	Report []string

	// captured stderr output from the test function (logs emitted by the test)
	Log []string

	// test outcome
	// (TestPassed, TestFailed, TestPanicked)
	Outcome TestOutcome

	// names of any tests that failed
	FailedTests []string

	t       TestingT
	checked bool
}

// Expect verifies that the test result (R) matches the expected outcome.
//
// When called without arguments it verifies that the test outcome is
// TestPassed with an empty test report.
//
// If called with a specific TestOutcome, it verifies that the test
// outcome matches the expected outcome.  If the expected outcome
func (r *R) Expect(exp ...any) {
	// TODO: cognitive complexity 27 > 15
	r.t.Helper()
	T().Helper()

	r.checked = true

	// in some cases there are multiple tests performed where if the first
	// test fails all following tests are certain to run, producing multiple
	// failures where only the first is helpful, so we set the IsRequired
	// option in addition to any options supplied
	opts := []any{opt.IsRequired(true)}

	var expectedReport []string
	for _, v := range exp {
		switch s := v.(type) {
		case string:
			expectedReport = append(expectedReport, s)
		case []string:
			expectedReport = append(expectedReport, s...)
		default:
			opts = append(opts, v)
		}
	}

	var rpt = make([]string, 0, len(r.Report))
	var sig = false
	for _, s := range r.Report {
		if sig = sig || strings.TrimSpace(s) != ""; sig {
			rpt = append(rpt, s)
		}
	}
	if r.Report = rpt; len(rpt) == 0 {
		r.Report = nil
	}

	expectedOutcome := TestPassed
	if opt.IsSet(opts, TestPanicked) {
		expectedOutcome = TestPanicked
	} else if opt.IsSet(opts, TestFailed) || len(expectedReport) > 0 {
		expectedOutcome = TestFailed
	}

	// setup an expectation for the failed tests; exactly what is expected
	// depends on the test outcome (it may not be tested at all).
	//
	// Instantiating the expectation here avoids having to do it in each case
	// that it is used.
	expectFailedTests := Expect(r.FailedTests, "failed tests")

	Expect(r.Outcome, "test outcome").To(Equal(expectedOutcome),
		append(opts, opt.FailureReport(func(...any) []string {
			report := []string{
				fmt.Sprintf("expected: %s", expectedOutcome),
				fmt.Sprintf("got     : %s", r.Outcome),
			}

			switch {
			case r.Outcome == TestPanicked:
				return append(report, fmt.Sprintf("recovered:\n  %[1]T(%[1]v)", r.Recovered))
			case len(r.Report) > 0:
				return slices.AppendToReport(report, r.Report, "with report:", opt.QuotedStrings(false))
			default:
				return report
			}
		}))...,
	)
	switch {
	case expectedOutcome == TestPanicked:
		Expect(r.Recovered, "recovered").IsNotNil(opts...)
		if len(expectedReport) > 0 {
			s := fmt.Sprintf("%v", r.Recovered)
			Expect(s, "recovered").To(ContainString(expectedReport[0]),
				append(opts, strings.Contains, opt.UnquotedStrings())...,
			)
		}

	case len(expectedReport) > 0:
		// first we check that the test we are running was identified as failed
		expectFailedTests.To(ContainItem(r.t.Name()), opts...)

		testfile := testFilename()
		if len(r.Report) > 0 {
			// if a report is expected we expect the first line of the report
			// to contain the name of the test file that was executing at the time
			Expect(r.Report[0]).To(ContainString(testfile),
				append(opts, opt.UnquotedStrings())...,
			)
		}

		// now we check that the report contains the expected lines
		Expect(r.Report).To(ContainSlice(expectedReport),
			append(opts, strings.Contains, opt.UnquotedStrings())...,
		)

	case expectedOutcome == TestFailed && opt.IsSet(opts, opt.IgnoreReport(true)):
		expectFailedTests.To(ContainItem(r.t.Name()), opts...)

	default:
		optx := []any{}
		for _, o := range opts {
			if o == opt.IsRequired(true) {
				continue
			}
			optx = append(optx, o)
		}
		expectFailedTests.Should(BeEmptyOrNil(), optx...)
		Expect(r.Report, "test report").Should(BeEmptyOrNil(),
			append(opts, opt.FailureReport(func(...any) []string {
				report := make([]string, 2, len(r.Report)+2)
				report[0] = "expected: <no report>"
				report[1] = "got:"
				for _, s := range r.Report {
					report = append(report, "| "+s)
				}
				return report
			}))...,
		)
	}
}

func (r *R) ExpectInvalid(report ...any) {
	r.t.Helper()
	r.Expect(append([]any{"<== INVALID TEST"}, report...)...)
}

// Test runs a function that exercises a test function returning an R that captures the
// test result:
//
// - the outcome (TestPassed, TestFailed, TestPanicked)
// - names of any tests that failed
// - any output
// - any logs
// - any value recovered from a panic
func Test(f func()) R {
	t, ok := T().(*testing.T)
	if !ok {
		panic("Test() must be called with a *testing.T test frame")
	}

	t.Helper()

	var recovered any
	stdout, stderr, outcome := runInternal(t, func(internal *testing.T) {
		testframe.Push(internal)
		defer func() {
			recovered = recover()
			testframe.Pop()
		}()
		f()
	})

	if recovered != nil {
		outcome = TestPanicked
	}

	failed := []string{}
	rxFailed := regexp.MustCompile(`--- FAIL: (.*) \((\d+\.\d+)s\)$`)

	newStdout := make([]string, 0, len(stdout))
	for _, s := range stdout {
		if match := rxFailed.FindStringSubmatch(s); len(match) > 0 {
			failed = append(failed, match[1])
		} else {
			newStdout = append(newStdout, s)
		}
	}
	stdout = newStdout
	if len(failed) == 0 {
		failed = nil
	}

	return R{
		t:           t,
		FailedTests: failed,
		Recovered:   recovered,
		Log:         stderr,
		Report:      stdout,
		Outcome:     outcome,
	}
}

// runInternalMatchAll is a test matcher function that always returns true and no error.
// It is used by runInternal to ensure that the internal test is always executed by the
// InternalTest runner, regardless of any match conditions specified for the current
// *testing.T test run.
func runInternalMatchAll(pat, match string) (bool, error) {
	return true, nil
}

// runInternal is a helper function that runs a specified test function as an
// internal test.
//
// It is used by Test() to run a test function in a separate test runner in order to
// inspect the outcome of the test function without that affecting the state of the
// current test.
func runInternal(t *testing.T, f func(*testing.T)) ([]string, []string, TestOutcome) {
	t.Helper()

	result := TestFailed
	stdout, stderr := Record(func() {
		it := []testing.InternalTest{{
			Name: t.Name(),
			F:    f,
		}}
		if testing.RunTests(runInternalMatchAll, it) {
			result = TestPassed
		}
	})

	// this is a bit of a hack: when run with -v, the internal test runner
	// emits a RUN line and a PASS line when the test passes, and a SKIP line
	// if the test is skipped. We remove these lines as we don't want to see them
	// in the captured output.
	//
	// (if the internal test fails we retain all of the output)
	//
	// NOTE: This only applies when processing the output of an internal test; that
	// is, a test that tests a test.  Output from normal tests is not filtered.
	for i := len(stdout) - 1; i >= 0; i-- {
		s := strings.TrimSpace(stdout[i])
		if strings.HasPrefix(s, "=== RUN") ||
			strings.HasPrefix(s, "=== NAME") ||
			strings.HasPrefix(s, "=== PAUSE") ||
			strings.HasPrefix(s, "=== CONT") ||
			strings.HasPrefix(s, "--- PASS: ") ||
			strings.HasPrefix(s, "--- SKIP: ") {
			stdout = append(stdout[:i], stdout[i+1:]...)
		}
	}

	return stdout, stderr, result
}

var isTestFile = func(s string) bool {
	return strings.HasSuffix(s, "_test.go")
}

// testFilename returns the name of the first test file (_test.go) that is found
// in the call stack.
func testFilename() string {
	pcs := make([]uintptr, 20)
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)

	frame, more := frames.Next()
	for more {
		if isTestFile(frame.File) {
			_, filename := path.Split(frame.File)
			return filename
		}
		frame, more = frames.Next()
	}

	return "<unknown test file>"
}
