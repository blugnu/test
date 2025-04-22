package test

import (
	"testing"
)

func TestBooleans(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		// ExpectFalse scenarios
		{Scenario: "ExpectFalse(false)",
			Act: func() { ExpectFalse(false) },
			Assert: func(result R) {
				result.Assert()
			},
		},
		{Scenario: "ExpectFalse(true)",
			Act: func() { ExpectFalse(true) },
			Assert: func(result R) {
				result.Assert([]string{
					"expected false, got true",
				})
			},
		},

		// ExpectTrue scenarios
		{Scenario: "ExpectTrue(true)",
			Act: func() { ExpectTrue(true) },
			Assert: func(result R) {
				result.Assert()
			},
		},
		{Scenario: "ExpectTrue(false)",
			Act: func() { ExpectTrue(false) },
			Assert: func(result R) {
				result.Assert("expected true, got false")
			},
		},

		// BeFalse scenarios
		{Scenario: "false.BeFalse()",
			Act: func() { Expect(false).To(BeFalse()) },
			Assert: func(result R) {
				result.Assert()
			},
		},
		{Scenario: "true.BeFalse()",
			Act: func() { Expect(true).To(BeFalse()) },
			Assert: func(result R) {
				result.Assert("expected false, got true")
			},
		},

		// BeTrue scenarios
		{Scenario: "true.BeTrue()",
			Act: func() { Expect(true).To(BeTrue()) },
			Assert: func(result R) {
				result.Assert()
			},
		},
		{Scenario: "false.BeTrue()",
			Act: func() { Expect(false).To(BeTrue()) },
			Assert: func(result R) {
				result.Assert("expected true, got false")
			},
		},

		// coverage for OneLineError
		{Scenario: "coverage/OneLineError",
			Act:    func() { BeTrue().OneLineError() },
			Assert: func(result R) {},
		},
	})
}

func ExampleExpectFalse() {
	With(ExampleTestRunner{})

	got := true
	ExpectFalse(got)

	// Output:
	// expected false, got true
}
