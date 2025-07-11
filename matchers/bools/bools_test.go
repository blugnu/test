package bools_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestBooleans(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: BeFalse
		{Scenario: "BeFalse when false",
			Act: func() { Expect(false).To(BeFalse()) },
		},
		{Scenario: "BeFalse when true",
			Act: func() { Expect(true).To(BeFalse()) },
			Assert: func(result *R) {
				result.Expect(
					"expected false, got true",
				)
			},
		},
		{Scenario: "ToNot BeFalse when false",
			Act: func() { Expect(false).ToNot(BeFalse()) },
			Assert: func(result *R) {
				result.Expect(
					"did not expect false",
				)
			},
		},
		{Scenario: "ToNot BeFalse when true",
			Act: func() { Expect(true).ToNot(BeFalse()) },
		},

		// MARK: BeTrue
		{Scenario: "BeTrue when true",
			Act: func() { Expect(true).To(BeTrue()) },
		},
		{Scenario: "BeTrue when false",
			Act: func() { Expect(false).To(BeTrue()) },
			Assert: func(result *R) {
				result.Expect(
					"expected true, got false",
				)
			},
		},
		{Scenario: "ToNot BeTrue when true",
			Act: func() { Expect(true).ToNot(BeTrue()) },
			Assert: func(result *R) {
				result.Expect(
					"did not expect true",
				)
			},
		},
	}...))
}
