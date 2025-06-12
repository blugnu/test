package test

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
	RunScenarios(
		func(tc *testcase, _ int) {
			// ACT
			result := tc.String()

			// ASSERT
			Expect(result).To(Equal(tc.result))
		},
		[]testcase{
			{TestOutcome: TestPassed, result: "TestPassed"},
			{TestOutcome: TestFailed, result: "TestFailed"},
			{TestOutcome: TestPanicked, result: "TestPanicked"},
			{TestOutcome: 99, result: "TestOutcome(99)"},
		})
}

func TestTest(t *testing.T) {
	With(t)

	Run("Test call made from a non-Test function", func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic, but did not get one")
			}
		}()

		test.Example()
		_ = Test(func() {})
	})

	Run("test panics", func() {
		result := Test(func() {
			panic("whoops!")
		})
		result.Expect(TestPanicked, "whoops!")
	})

	Run("additional output to stdout", func() {
		result := Test(func() {
			fmt.Println("")
			fmt.Println("some preamble output")
			fmt.Println("")
			Expect(true).To(BeFalse())
			fmt.Println("additional output")
		})

		// the report consists only of test report output
		result.Expect(TestFailed, "expected false, got true")
	})
}

func TestR_Expect(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
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
	})
}

func Test_runInternal(t *testing.T) {
	With(t)

	Run("runInternalMatchAll coverage", func() {
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
	})

	Run("cleans verbose output", func() {
		// This test simulates the output of a successful (passing) test that is run
		// with the -v flag.  We ensure that the RUN: and PASS: lines are removed from
		// the output to avoid false negatives in tests that expect an empty output for
		// a passing test.

		t := T().(*testing.T)

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
	})
}

func Test_testFilename(t *testing.T) {
	With(t)

	Run("called from a test file", func() {
		// ACT
		result := testFilename()

		// ASSERT
		Expect(result).To(Equal("test_test.go"))
	})

	Run("called from a non-test file (simulated)", func() {
		defer Restore(Original(&isTestFile).ReplacedBy(func(s string) bool { return false }))

		// ACT
		result := testFilename()

		// ASSERT
		Expect(result).To(Equal("<unknown test file>"))
	})
}
