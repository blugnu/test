package test

import (
	"os"
	"testing"
)

func TestHelperTest(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// HadOutcome tests
		{scenario: "HelperTest{result:true}.HadOutcome(true)",
			act: func(t T) {
				HelperTest{t: t, result: true}.HadOutcome(true)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{result:false}.HadOutcome(false)",
			act: func(t T) {
				HelperTest{t: t, result: false}.HadOutcome(false)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{result:true}.HadOutcome(false)",
			act: func(t T) {
				HelperTest{t: t, result: true}.HadOutcome(false)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{result:true}.HadOutcome(false)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: FAIL",
					"got   : PASS",
				})
			},
		},
		{scenario: "HelperTest{result:false}.HadOutcome(true)",
			act: func(t T) {
				HelperTest{t: t, result: false}.HadOutcome(true)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{result:false}.HadOutcome(true)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: PASS",
					"got   : FAIL",
				})
			},
		},

		// DidFail tests
		{scenario: "HelperTest{result:false}.DidFail()",
			act: func(t T) {
				HelperTest{t: t, result: false}.DidFail()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{result:true}.DidFail()",
			act: func(t T) {
				HelperTest{t: t, result: true}.DidFail()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{result:true}.DidFail()")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: FAIL",
					"got   : PASS",
				})
			},
		},

		// DidPass tests
		{scenario: "HelperTest{result:true}.DidPass()",
			act: func(t T) {
				HelperTest{t: t, result: true}.DidPass()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{result:false}.DidPass()",
			act: func(t T) {
				HelperTest{t: t, result: false}.DidPass()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{result:false}.DidPass()")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: PASS",
					"got   : FAIL",
				})
			},
		},

		// DidNotPanic tests
		{scenario: "HelperTest{recovered:nil}.DidNotPanic()",
			act: func(t T) {
				HelperTest{t: t, recovered: nil}.DidNotPanic()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{recovered:42}.DidNotPanic()",
			act: func(t T) {
				HelperTest{t: t, recovered: 42}.DidNotPanic()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{recovered:42}.DidNotPanic()")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: (did not panic)",
					"got   : panic: 42",
				})
			},
		},

		// DidPanic tests
		{scenario: "HelperTest{recovered:42}.DidPanic(42)",
			act: func(t T) {
				HelperTest{t: t, recovered: 42}.DidPanic(42)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "HelperTest{recovered:42}.DidPanic(64)",
			act: func(t T) {
				HelperTest{t: t, recovered: 42}.DidPanic(64)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{recovered:42}.DidPanic(64)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: panic: int: 64",
					"got   : panic: int: 42",
				})
			},
		},
		{scenario: "HelperTest{recovered:nil}.DidPanic(42)",
			act: func(t T) {
				HelperTest{t: t, recovered: nil}.DidPanic(42)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("HelperTest{recovered:nil}.DidPanic(42)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: panic: int: 42",
					"got   : (did not panic)",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}

func TestHelper(t *testing.T) {
	t.Run("matchAll", func(t *testing.T) {
		ok, err := matchAll("any", "none")
		if !ok || err != nil {
			t.Errorf("matchAll() = %v, %v, want %v, %v", ok, err, true, nil)
		}
	})

	testcases := []struct {
		scenario string
		act      func(T) HelperTest
		assert   func(HelperTest)
	}{
		{scenario: "Helper(func{}, \"\")",
			act: func(t T) HelperTest {
				return Helper(t, func(st *testing.T) {})
			},
			assert: func(test HelperTest) {
				test.DidPass()
			},
		},
		{scenario: "Helper(func{t.Fail()})",
			// tests that a Helper that fails does not fail the test if no outcome is specified;
			// the expectation is that tests will be explicitly performed by the caller,
			// using the returned HelperTest
			act: func(t T) HelperTest {
				return Helper(t, func(st *testing.T) { st.Fail() })
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Helper(func{t.Fail()})")
				test.Report.Contains("FAIL")
			},
		},
		{scenario: "Helper(func{t.Errorf(msg)})",
			act: func(t T) HelperTest {
				return Helper(t, func(st *testing.T) { st.Error("report message") })
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Helper(func{t.Errorf(msg)})")
				test.Report.Contains("report message")
			},
		},
		{scenario: "Helper(func{panic(42)})",
			act: func(t T) HelperTest {
				return Helper(t, func(st *testing.T) { panic(42) })
			},
			assert: func(test HelperTest) {
				test.DidPanic(42)
				test.Report.Contains("unexpected panic: int: 42")
			},
		},
		{scenario: "Helper(func{}) (verbose test output)",
			act: func(t T) HelperTest {
				return Helper(t, func(st *testing.T) {
					// simulates the additional output from a test that passes
					// when run with go test -v
					os.Stdout.WriteString("=== RUN   TheNameOfSomeTest\n")
					os.Stdout.WriteString("--- PASS: TheNameOfSomeTest (0.00s)\n")
				})
			},
			assert: func(test HelperTest) {
				test.Report.IsEmpty()
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(tc.act(t))
		})
	}
}
