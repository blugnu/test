package test //nolint: testpackage // tests rely on access to private fields

import (
	"fmt"
	"testing"

	"github.com/blugnu/test/test"
)

func TestTestOutcome(t *testing.T) {
	With(t)

	type testcase struct {
		TestOutcome
		result string
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// ACT
			result := tc.String()

			// ASSERT
			Expect(result).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{TestOutcome: TestPassed, result: "TestPassed"},
			{TestOutcome: TestFailed, result: "TestFailed"},
			{TestOutcome: TestPanicked, result: "TestPanicked"},
			{TestOutcome: 99, result: "TestOutcome(99)"},
		}),
	))
}

func TestTestHelper(t *testing.T) {
	With(t)

	Run(Test("Test call made from a non-Test function", func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic, but did not get one")
			}
		}()

		test.Example()
		_ = TestHelper(func() {})
	}))

	Run(Test("test panics", func() {
		result := TestHelper(func() {
			panic("whoops!")
		})
		result.Expect(TestPanicked, "whoops!")
	}))

	Run(Test("additional output to stdout", func() {
		result := TestHelper(func() {
			fmt.Println("")
			fmt.Println("some preamble output")
			fmt.Println("")
			Expect(true).To(BeFalse())
			fmt.Println("additional output")
		})

		// the report consists only of test report output
		result.Expect(TestFailed, "expected false, got true")
	}))
}

func TestR_Expect(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no arguments",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestPassed,
				}
				sut.Expect()
			},
			Assert: func(result *R) {
				result.ExpectInvalid("R.Expect: no arguments; an expected TestOutcome and/or test report are required")
			},
		},
		{Scenario: "expected to panic (no report)",
			Act: func() {
				sut := R{
					t:         T(),
					Outcome:   TestPanicked,
					Recovered: "recovered",
				}
				sut.Expect(TestPanicked)
			},
		},
		{Scenario: "expected to panic (matched recovered string)",
			Act: func() {
				sut := R{
					t:         T(),
					Outcome:   TestPanicked,
					Recovered: "recovered",
				}
				sut.Expect(TestPanicked, "recovered")
			},
		},
		{Scenario: "expected to panic (too many strings specified)",
			Act: func() {
				sut := R{
					t:         T(),
					Outcome:   TestPanicked,
					Recovered: "recovered",
				}
				sut.Expect(TestPanicked, "recovered", "and a second string")
			},
			Assert: func(result *R) {
				result.ExpectInvalid("R.Expect: only 1 string may be specified to match a recovered value from an expected panic (got 2)")
			},
		},
		{Scenario: "expected to pass, but panicked",
			Act: func() {
				sut := R{
					t:         T(),
					Outcome:   TestPanicked,
					Recovered: "recovered",
				}
				sut.Expect(TestPassed)
			},
			Assert: func(result *R) {
				result.Expect(
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestPanicked",
					"recovered:",
					"  string(recovered)",
				)
			},
		},
		{Scenario: "expected to pass, but failed (no report)",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
				}
				sut.Expect(TestPassed)
			},
			Assert: func(result *R) {
				result.Expect(
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestFailed",
				)
			},
		},
		{Scenario: "expected to pass, but failed (with report)",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
					Report:  []string{"actual report"},
				}
				sut.Expect(TestPassed)
			},
			Assert: func(result *R) {
				result.Expect(
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestFailed",
					"with report:",
					"| actual report",
				)
			},
		},
		{Scenario: "expected to fail with no report, failed with report",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
					Report:  []string{"actual report"},
				}
				sut.Expect(TestFailed)
			},
			Assert: func(result *R) {
				result.Expect(
					"test report:",
					"  expected: <no report>",
					"  got:",
					"  | actual report",
				)
			},
		},
		{Scenario: "failed with no report",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
				}
				sut.Expect(TestFailed, "some expected failure report")
			},
			Assert: func(result *R) {
				result.ExpectWarning("test failed as expected, but no test report or failures were recorded")
			},
		},
	}...))
}

func Test_runInternal(t *testing.T) {
	With(t)

	Run(Test("runInternalMatchAll coverage", func() {
		// NOTE: This test is provided to provide coverage only.
		//
		// The runInternalMatchAll function is passed to the RunTests() function called
		// by runInternal() when testing a test.
		//
		// For some reason the function is never called (perhaps it would be if go test
		// were executed with a test pattern ?).  In any event, we 'test' it here to ensure
		// that it is covered.

		result, err := runInternalMatchAll("", "")

		Expect(err).IsNil()
		Expect(result).To(BeTrue())
	}))

	Run(Test("cleans verbose output", func() {
		// This test simulates the output of a successful (passing) test that is run
		// with the -v flag.  We ensure that the RUN: and PASS: lines are removed from
		// the output to avoid false negatives in tests that expect an empty output for
		// a passing test.

		t := RequireType[*testing.T](T())

		// ACT
		stdout, stderr, outcome := runInternal(t, func(st *testing.T) {
			fmt.Println("=== RUN: this output should be removed")
			fmt.Println("=== PAUSE: this output should be removed")
			fmt.Println("=== CONT: this output should be removed")
			fmt.Println("--- PASS: this output should be removed")
		})

		// ASSERT
		Expect(outcome).To(Equal(TestPassed))
		Expect(stdout).Should(BeEmpty())
		Expect(stderr).Should(BeNil())
	}))
}

func Test_testFilename(t *testing.T) {
	With(t)

	Run(Test("called from a test file", func() {
		// ACT
		result := testFilename()

		// ASSERT
		Expect(result).To(Equal("test-helper_test.go"))
	}))

	Run(Test("called from a non-test file (simulated)", func() {
		defer Restore(Original(&isTestFile).ReplacedBy(func(s string) bool { return false }))

		// ACT
		result := testFilename()

		// ASSERT
		Expect(result).To(Equal("<unknown test file>"))
	}))
}

func Test_analyseReport(t *testing.T) {
	With(t)

	Run(Test("empty report", func() {
		// ACT
		output, report, failed := analyseReport([]string{})

		// ASSERT
		Expect(output, "output").Should(BeEmptyOrNil())
		Expect(report, "report").Should(BeEmptyOrNil())
		Expect(failed, "failed").Should(BeEmptyOrNil())
	}))

	Run(Test("report with one test failure", func() {
		// ACT
		output, report, failed := analyseReport([]string{
			"--- FAIL: TestSomething (0.0s)",
			"    stdout output (if any)",
			"    something_test.go:112:",
			"        report line 1",
			"        report line 2",
		})

		// ASSERT
		Expect(output, "output").To(EqualSlice([]string{"    stdout output (if any)"}))
		Expect(report, "report").To(EqualSlice([]string{
			"    something_test.go:112:",
			"        report line 1",
			"        report line 2",
		}))
		Expect(failed, "failed").To(EqualSlice([]string{"TestSomething"}))
	}))

	Run(Test("misreported test location (simulates missing t.Helper)", func() {
		// ACT
		output, report, failed := analyseReport([]string{
			"--- FAIL: TestSomething (0.0s)",
			"    stdout output (if any)",
			"    something.go:112: test failed", // should be something_test.go:112, i.e. missing _test suffix
		})

		// ASSERT
		Expect(output, "output").Should(BeEmptyOrNil())
		Expect(report, "report").To(EqualSlice([]string{
			"WARNING: check test location (missing a T().Helper() call?)",
			"report:",
			"| --- FAIL: TestSomething (0.0s)",
			"|     stdout output (if any)",
			"|     something.go:112: test failed",
		}))
		Expect(failed, "failed").Should(BeEmptyOrNil())
	}))
}
