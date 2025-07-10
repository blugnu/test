package test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/matchers/slices"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type TestOutcome int

const (
	TestPassed TestOutcome = iota
	TestFailed
	TestPanicked
)

const (
	cFailedTests = "failed tests"
	cTestOutcome = "test outcome"
	cTestReport  = "test report"
	cRecovered   = "recovered"
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

	// Stack is the stack trace captured when recovering from a panicked test
	// or nil if the test did not panic
	Stack []byte

	checked bool
	t       TestingT
}

// Expect verifies that the test result (R) matches the expected outcome.
//
// At least one argument must be provided to Expect() to specify the expected
// outcome of the test. The arguments can be:
//
// - a TestOutcome value (TestPassed, TestFailed, TestPanicked)
// - a string or slice of strings that are expected to be present in the test report
// - a combination of the above, with options to control the assertion behavior
//
// The function will check the test outcome, the test report, and any recovered value
// from a panic, and will fail the test if any of the expectations are not met.
//
// Currently, any additional output to stdout produced by the test function is
// ignored and must not be specified in any expected test report.
//
// If no arguments are provided, Expect() will fail the test with an error message
// indicating that at least one expected outcome or report line is required.
//
// If the test outcome is expected to be TestPanicked, the first argument must be a
// TestOutcome value (TestPanicked) with a single string argument that is expected to
// match the string representation (%v) of the value recovered from the panic.
func (r *R) Expect(exp ...any) {
	r.t.Helper()

	if len(exp) == 0 {
		// if no arguments are given, we expect the test to have passed
		test.Invalid("R.Expect: no arguments; an expected TestOutcome and/or test report are required")
	}

	// we have some expectations so can mark the result as having been checked
	r.checked = true

	expectedOutcome, expectedReport, opts := r.analyseArgs(exp...)

	// multiple assertions will be made and we want to fail early, so we
	// set the IsRequired option if not already set.
	if !opt.IsSet(opts, opt.IsRequired(true)) {
		opts = append(opts, opt.IsRequired(true))
	}

	r.assertOutcome(expectedOutcome, opts...)

	switch {
	case expectedOutcome == TestPanicked:
		r.assertPanicked(expectedReport, opts...)

	case len(expectedReport) > 0:
		if len(r.FailedTests) == 0 {
			test.Warning("test failed as expected, but no test report or failures were recorded")
		}

		// first we check that the test we are running was identified as failed
		Expect(r.FailedTests, cFailedTests).To(ContainItem(r.t.Name()), opts...)

		testfile := testFilename()
		if len(r.Report) > 0 {
			// FUTURE: this only reports the first line of the report when it fails,
			// which is not helpful in the case of a badly formatted report.
			//
			// This test should use a matcher that will report the entire report
			// when it fails, so that the test author can see what went wrong.
			//
			// e.g. a SliceBeginsWith() matcher (rather than explicitly testing only
			// the Report[0] item)

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
		Expect(r.FailedTests, cFailedTests).To(ContainItem(r.t.Name()), opts...)

	default:
		r.assertWhenTestPassed(opts...)
	}
}

// ExpectInvalid verifies that the test result (R) indicates an invalid test; that is,
// the test.Invalid() function was called during evaluation of the test, indicating
// some problem that makes the test result unreliable or meaningless.
//
// When called without arguments it verifies that the test outcome is failed and that
// the test report consists only of the '<== INVALID TEST' label.  i.e. if any additional
// report lines are present in the test report then the ExpectInvalid() call will itself
// fail.
//
// If the test report is not significant to a test, it must be explicitly ignored by
// passing the opt.IgnoreReport(true) option:
//
//	result.ExpectInvalid(opt.IgnoreReport(true))
//
// # Invalid Tests
//
// The Go standard library testing framework does not provide a way to mark a test
// as invalid; the test.Invalid() function fails a test and emits a message labelled
// with '<== INVALID TEST'.
//
// When testing for an invalid test:
//
// - the test outcome is expected to be TestFailed
// - the test report is expected to start with the '<== INVALID TEST' label
// - the test report is expected to contain any report lines specified
func (r *R) ExpectInvalid(report ...any) {
	r.t.Helper()
	r.Expect(append([]any{"<== INVALID TEST"}, report...)...)
}

// ExpectWarning verifies that the test result (R) contains a warning; that is,
// the test.Warning() function was called during evaluation of the test.
func (r *R) ExpectWarning(msg string) {
	r.t.Helper()
	r.Expect("<== WARNING: " + msg)
}

// analyseArgs processes the expected arguments passed to Expect().
// It separates the expected report strings from the options and determines
// the expected outcome of the test based on the options provided.
// It returns the expected outcome, the report strings, and the options.
func (r *R) analyseArgs(exp ...any) (TestOutcome, []string, []any) {
	var (
		report  []string
		opts    []any
		outcome = TestPassed
	)

	// separate the expected report strings from any options
	for _, v := range exp {
		switch s := v.(type) {
		case string:
			report = append(report, s)
		case []string:
			report = append(report, s...)
		default:
			opts = append(opts, v)
		}
	}

	if opt.IsSet(opts, TestPanicked) {
		outcome = TestPanicked
	} else if opt.IsSet(opts, TestFailed) || len(report) > 0 {
		outcome = TestFailed
	}

	return outcome, report, opts
}

// assertOutcome checks that the test outcome matches the expected outcome.
//
// The function is used internally by Expect() to verify the test outcome.
func (r *R) assertOutcome(expected TestOutcome, opts ...any) {
	T().Helper()

	Expect(r.Outcome, cTestOutcome).To(Equal(expected),
		append(opts, opt.FailureReport(func(...any) []string {
			report := []string{
				fmt.Sprintf("expected: %s", expected),
				fmt.Sprintf("got     : %s", r.Outcome),
			}

			switch {
			case r.Outcome == TestPanicked:
				report = append(report, "")
				report = append(report, "recovered:")
				report = append(report, fmt.Sprintf("  %[1]T(%[1]v)", r.Recovered))
				if trace := panics.StackTrace(r.Stack, opts...); trace != nil {
					report = append(report, "")
					report = append(report, "stack:")
					report = append(report, trace...)
				}
				return report

			case len(r.Report) > 0:
				return slices.AppendToReport(report, r.Report, "with report:", opt.QuotedStrings(false))

			default:
				return report
			}
		}))...,
	)
}

// assertPanicked checks that the test result indicates a panic and that the
// recovered value matches the expected report (if specified).
//
// The function is used internally by Expect() to verify the test outcome when a panic
// is expected.
//
// FUTURE: accept an error in the arguments and assert that the recovered value
// matches the error.  This would provide more robust assertions for error panics
//
// For now, if an expected report is given, we check that the recovered value
// matches the first line of the report. If the expected report has > 1 line, we
// fail the test as invalid, as we don't know how to handle that case yet.
func (r *R) assertPanicked(expectedReport []string, opts ...any) {
	T().Helper()

	if len(expectedReport) > 1 {
		test.Invalid(fmt.Sprintf("R.Expect: only 1 string may be specified to match a recovered value from an expected panic (got %d)", len(expectedReport)))
	}

	Expect(r.Recovered, cRecovered).IsNotNil(opts...)

	if len(expectedReport) == 0 {
		return
	}

	s := fmt.Sprintf("%v", r.Recovered)
	Expect(s, cRecovered).To(ContainString(expectedReport[0]),
		append(opts, strings.Contains, opt.UnquotedStrings())...,
	)
}

// assertWhenTestPassed checks that there were no failed tests and no report
func (r *R) assertWhenTestPassed(opts ...any) {
	T().Helper()

	// in this case, we want the failure report to provide details of
	// any unexpectedly failed tests and an unexpected report, so we
	// remove any IsRequired option to ensure that both expectations
	// are evaluated
	opts = opt.Unset(opts, opt.IsRequired(true))

	Expect(r.FailedTests, cFailedTests).Should(BeEmptyOrNil(), opts...)

	Expect(r.Report, cTestReport).Should(BeEmptyOrNil(),
		append(opts, opt.FailureReport(func(...any) []string {
			const preambleLen = 2

			report := make([]string, preambleLen, len(r.Report)+preambleLen)
			report[0] = "expected: <no report>"
			report[1] = "got:"
			for _, s := range r.Report {
				report = append(report, "| "+s)
			}

			return report
		}))...,
	)
}

// TestHelper runs a function that executes a function in an internal test runner,
// independent of the current test, returning an R value that captures the
// following:
//
//   - the test outcome (TestPassed, TestFailed, TestPanicked)
//   - names of any tests that failed
//   - stdout output
//   - stderr output
//   - any value recovered from a panic
//
// This function is intended to be used to test helper functions.  For example,
// it is used extensively in the blugnu/test package itself, to test the
// test framework.
func TestHelper(f func()) R {
	t, ok := T().(*testing.T)
	if !ok {
		panic("TestHelper() must be called with a *testing.T test frame")
	}

	t.Helper()

	var (
		recovered any
		stack     []byte
	)
	stdout, stderr, outcome := runInternal(t, func(internal *testing.T) {
		testframe.Push(internal)
		defer func() {
			recovered = recover()
			testframe.Pop()

			if recovered != nil {
				const bufsize = 65536

				buf := make([]byte, bufsize)
				n := runtime.Stack(buf, false)
				stack = buf[0 : n-1]
			}
		}()
		f()
	})

	if recovered != nil {
		outcome = TestPanicked
	}

	_, report, failed := analyseReport(stdout)

	return R{
		t:           t,
		FailedTests: failed,
		Recovered:   recovered,
		Log:         stderr,
		Report:      report,
		Stack:       stack,
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
// It is used to run a test function in a separate test runner in order to
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
	const skipFrames = 2
	const maxFrames = 64

	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(skipFrames, pcs)
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)

	frame, more := frames.Next()
	for more {
		if isTestFile(frame.File) {
			return filepath.Base(frame.File)
		}
		frame, more = frames.Next()
	}

	return "<unknown test file>"
}

// analyseReport removes any initial empty lines from the test report
// and extracts the names of any tests.
//
// It returns stdout output (additional to any test report output), the report
// and a slice of names of tests that failed.
//
// If any slice is empty, nil is returned instead of an empty slice.
func analyseReport(stdout []string) ([]string, []string, []string) {
	T().Helper()

	failed := []string{}
	namePattern := regexp.MustCompile(`--- FAIL: (.*) \((\d+\.\d+)s\)$`)
	locnPattern := regexp.MustCompile(`[\s]*(.*)_test\.go:[0-9]+:`)

	// FUTURE: properly support multiple test failures in a report.
	//
	// The implementation below assumes one test failure in the report, of
	// the form:
	//
	// --- FAIL: <test name> (<duration>s)
	//     <stdout output (if any)>
	//     <test filename>:<line no>:
	//         report line 1
	//         report line 2
	//
	// But if the output consists of multiple test failures it will be:
	//
	// --- FAIL: <test name> (<duration>s)
	//     <stdout output (if any)>
	//     <test filename>:<line no>:
	//         report line 1
	//         report line 2
	//         ...
	//     <test filename>:<line no>:
	//         report line 1
	//     ...

	output := make([]string, 0, len(stdout))
	report := make([]string, 0, len(stdout))
	inReport := false
	for _, s := range stdout {
		if match := namePattern.FindStringSubmatch(s); len(match) > 0 {
			failed = append(failed, match[1])
			continue
		}

		if match := locnPattern.FindStringSubmatch(s); len(match) > 0 {
			// if we find a location line we assume that this and following lines
			// are part of the report for the test that failed.
			inReport = true
		}

		switch inReport {
		case true:
			report = append(report, s)
		default:
			output = append(output, s)
		}
	}

	if len(failed) > 0 && len(report) == 0 {
		// if we have a failed test there MUST be some failure report, otherwise we have
		// been presented with a report that does not conform to the expected layout
		//
		// if that's the case, we will return the original stdout as the report,
		// with a warning that the report does not conform to the expected layout
		report = []string{"WARNING: check test location (missing a T().Helper() call?)"}
		report = slices.AppendToReport(report, stdout, "report:", opt.UnquotedStrings())
		return nil, report, nil
	}

	ornil := func(s []string) []string {
		if len(s) == 0 {
			return nil
		}
		return s
	}

	return ornil(output), ornil(report), ornil(failed)
}
