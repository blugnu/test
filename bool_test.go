package test

import (
	"testing"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestBooleans(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "ExpectFalse when false",
			Act: func() { ExpectFalse(false) },
		},
		{Scenario: "ExpectTrue when true",
			Act: func() { ExpectTrue(true) },
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
		{Scenario: "ExpectFalse with name",
			Act: func() { ExpectFalse(true, "this will fail") },
			Assert: func(result *R) {
				result.Expect(
					"this will fail:",
					"  expected false, got true",
				)
			},
		},
		{Scenario: "ExpectFalse with custom failure report",
			Act: func() {
				ExpectFalse(true, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "ExpectFalse with name and custom failure report",
			Act: func() {
				ExpectFalse(true, "this will fail", opt.FailureReport(func(...any) []string {
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
		{Scenario: "ExpectTrue with name",
			Act: func() { ExpectTrue(false, "this will fail") },
			Assert: func(result *R) {
				result.Expect(
					"this will fail:",
					"  expected true, got false",
				)
			},
		},
		{Scenario: "ExpectTrue with custom failure report",
			Act: func() {
				ExpectTrue(false, opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
		{Scenario: "ExpectTrue with name and custom failure report",
			Act: func() {
				ExpectTrue(false, "this will fail", opt.FailureReport(func(...any) []string {
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
	})
}

func ExampleExpectFalse() {
	test.Example()

	ExpectFalse(true)

	// Output:
	// expected false, got true
}

func ExampleExpectTrue() {
	test.Example()

	ExpectTrue(false)

	// Output:
	// expected true, got false
}

func ExampleBeFalse() {
	test.Example()

	Expect(true).To(BeFalse())

	// Output:
	// expected false, got true
}

func ExampleBeTrue() {
	test.Example()

	Expect(false).To(BeTrue())

	// Output:
	// expected true, got false
}
