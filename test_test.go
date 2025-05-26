package test

import (
	"fmt"
	"testing"
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

		With(ExampleTestRunner{})
		_ = Test(func() {})
	})

	Run("test panics", func() {
		result := Test(func() {
			panic("whoops!")
		})
		result.Expect(TestPanicked, "whoops!")
	})
}

func TestR_Expect(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "no arguments/test passed",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestPassed,
				}
				sut.Expect()
			},
		},
		{Scenario: "no arguments/test failed with no report",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
				}
				sut.Expect()
			},
			Assert: func(result *R) {
				result.Expect(
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestFailed",
				)
			},
		},
		{Scenario: "no arguments/test failed with report",
			Act: func() {
				sut := R{
					t:       T(),
					Outcome: TestFailed,
					Report:  []string{"report"},
				}
				sut.Expect()
			},
			Assert: func(result *R) {
				result.Expect(
					// failed on outcome
					TestFilename(),
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestFailed",
					// failed on test report
					TestFilename(),
					"test report:",
					"  expected: <no report>",
					"  got:",
					"  | report",
				)
			},
		},
		{Scenario: "expect that a test that panicked did panic, no report expectations",
			Act: func() {
				sut := R{
					t:         T(),
					Outcome:   TestPanicked,
					Recovered: "recovered",
				}
				sut.Expect(TestPanicked)
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
					// failed on outcome
					TestFilename(),
					"test outcome:",
					"  expected: TestPassed",
					"  got     : TestPanicked",
					// failed on recovered value
					TestFilename(),
					"unexpected panic:",
					"  recovered: string(recovered)",
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
		//
		// FUTURE: perhaps this should be handled by R.Expect() instead of runInternal()?

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
		Expect(stdout).IsEmpty()
		Expect(stderr).IsNil()
	})
}
