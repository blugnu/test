package expect_test

import (
	"testing"

	bg "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

//nolint:dupl // code is similar but not identical
func TestFalse(t *testing.T) {
	bg.With(t)

	bg.Run(bg.HelperTests([]bg.HelperScenario{
		{Scenario: "expect.False when false",
			Act: func() { expect.False(false) },
		},
		{Scenario: "expect.True when false",
			Act: func() {
				expect.False(true)
				expect.False(true)
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"bool_test.go",
					"expected false",
					"bool_test.go",
					"expected false",
				)
				bg.Expect(result.Report).Should(bg.HaveLen(4)) // both failures reported
			},
		},

		// supported options
		{Scenario: "expect.False with name",
			Act: func() { expect.False(true, "this will fail") },
			Assert: func(result *bg.R) {
				result.Expect(
					"this will fail:",
					"  expected false",
				)
			},
		},
		{Scenario: "expect.False with custom failure report",
			Act: func() {
				expect.False(true, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "expect.False with name and custom failure report",
			Act: func() {
				expect.False(true, "this will fail", opt.FailureReport(func(...any) []string {
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

//nolint:dupl // code is similar but not identical
func TestTrue(t *testing.T) {
	bg.With(t)

	bg.Run(bg.HelperTests([]bg.HelperScenario{
		{Scenario: "expect.True when true",
			Act: func() { expect.True(true) },
		},
		{Scenario: "expect.True when false",
			Act: func() {
				expect.True(false)
				expect.True(false)
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"bool_test.go",
					"expected true",
					"bool_test.go",
					"expected true",
				)
				bg.Expect(result.Report).Should(bg.HaveLen(4)) // both failures reported
			},
		},

		// supported options
		{Scenario: "expect.True with name",
			Act: func() { expect.True(false, "this will fail") },
			Assert: func(result *bg.R) {
				result.Expect(
					"this will fail:",
					"  expected true",
				)
			},
		},
		{Scenario: "expect.True with custom failure report",
			Act: func() {
				expect.True(false, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *bg.R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "expect.True with name and custom failure report",
			Act: func() {
				expect.True(false, "this will fail", opt.FailureReport(func(...any) []string {
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
