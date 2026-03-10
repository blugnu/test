package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestBooleans(t *testing.T) {
	With(t)

	Run(
		HelperTests([]HelperScenario{
			{Scenario: "expect.False when false",
				Act: func() { expect.False(false) },
			},
			{Scenario: "expect.True when true",
				Act: func() { expect.True(true) },
			},
			{Scenario: "BeFalse when false",
				Act: func() { Expect(false).To(BeFalse()) },
			},
			{Scenario: "ToNot BeFalse when true",
				Act: func() { Expect(true).ToNot(BeFalse()) },
			},
			{Scenario: "BeTrue when true",
				Act: func() { Expect(true).To(BeTrue()) },
			},
			{Scenario: "ToNot BeTrue when false",
				Act: func() { Expect(false).ToNot(BeTrue()) },
			},

			// supported options
			{Scenario: "expect.False with name",
				Act: func() { expect.False(true, "this will fail") },
				Assert: func(result *R) {
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
				Assert: func(result *R) {
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
				Assert: func(result *R) {
					result.Expect(
						"this will fail:",
						"  custom failure report",
					)
				},
			},
			{Scenario: "expect.True with name",
				Act: func() { expect.True(false, "this will fail") },
				Assert: func(result *R) {
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
				Assert: func(result *R) {
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
				Assert: func(result *R) {
					result.Expect(
						"this will fail:",
						"  custom failure report",
					)
				},
			},
			{Scenario: "BeFalse with custom failure report",
				Act: func() {
					Expect(true).To(BeFalse(), opt.FailureReport(func(...any) []string {
						return []string{"custom failure report"}
					}))
				},
				Assert: func(result *R) {
					result.Expect(
						"custom failure report",
					)
				},
			},
			{Scenario: "BeTrue with custom failure report",
				Act: func() {
					Expect(false).To(BeTrue(), opt.FailureReport(func(...any) []string {
						return []string{"custom failure report"}
					}))
				},
				Assert: func(result *R) {
					result.Expect(
						"custom failure report",
					)
				},
			},
		}...),
	)
}

// =================================================================
// MARK: examples
// =================================================================

// ExampleBeFalse demonstrates the use of the BeFalse matcher, showing
// how to use the matcher to assert that a boolean value is false and
// providing an example of the expected default output when the matcher
// fails.
func ExampleBeFalse() {
	test.Example()

	Expect(true).To(BeFalse())

	// Output:
	// expected false
}

// ExampleBeTrue demonstrates the use of the BeTrue matcher, showing
// how to use the matcher to assert that a boolean value is true and
// providing an example of the expected default output when the matcher
// fails.
func ExampleBeTrue() {
	test.Example()

	Expect(false).To(BeTrue())

	// Output:
	// expected true
}
