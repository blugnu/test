package require_test

import (
	"testing"

	bg "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/require"
	"github.com/blugnu/test/test"
)

//nolint:dupl	// code is similar but not identical
func TestFalse(t *testing.T) {
	bg.With(t)

	bg.Run(bg.HelperTests([]bg.HelperScenario{
		{Scenario: "require.False when false",
			Act: func() { require.False(false) },
		},
		{Scenario: "require.True when false",
			Act: func() {
				require.False(true)
				require.False(true) // not reached
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"bool_test.go",
					"expected false",
				)
				bg.Expect(result.Report).Should(bg.HaveLen(2)) // no other failures reported
			},
		},

		// supported options
		{Scenario: "require.False with name",
			Act: func() { require.False(true, "this will fail") },
			Assert: func(result *bg.R) {
				result.Expect(
					"this will fail:",
					"  expected false",
				)
			},
		},
		{Scenario: "require.False with custom failure report",
			Act: func() {
				require.False(true, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "require.False with name and custom failure report",
			Act: func() {
				require.False(true, "this will fail", opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
	}...))
}

//nolint:dupl	// code is similar but not identical
func TestTrue(t *testing.T) {
	bg.With(t)

	bg.Run(bg.HelperTests([]bg.HelperScenario{
		{Scenario: "require.True when true",
			Act: func() { require.True(true) },
		},
		{Scenario: "require.True when false",
			Act: func() {
				require.True(false)
				require.True(false) // not reached
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"bool_test.go",
					"expected true",
				)
				bg.Expect(result.Report).Should(bg.HaveLen(2)) // no other failures reported
			},
		},

		// supported options
		{Scenario: "require.True with name",
			Act: func() { require.True(false, "this will fail") },
			Assert: func(result *bg.R) {
				result.Expect(
					"this will fail:",
					"  expected true",
				)
			},
		},
		{Scenario: "require.True with custom failure report",
			Act: func() {
				require.True(false, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "require.True with name and custom failure report",
			Act: func() {
				require.True(false, "this will fail", opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
	}...))
}

func ExampleFalse() {
	test.Example()

	a := 1

	// only one test failure is reported by these tests

	require.False(a == 1)    // will fail and stop the test here
	expect.False(a+1-1 == 1) // this would also fail but is not evaluated

	// Output:
	// expected false
}

func ExampleTrue() {
	test.Example()

	a := 1

	// only one test failure is reported by these tests

	require.True(a == 0) // will fail and stop the test here
	expect.True(a < 0)   // this would also fail but is not evaluated

	// Output:
	// expected true
}
