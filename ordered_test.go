package test_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestBeGreaterThan(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expected greater than and was greater than",
			Act: func() {
				Expect(2).To(BeGreaterThan(1))
			},
		},
		{Scenario: "expected not greater than and was equal",
			Act: func() {
				Expect(2).ToNot(BeGreaterThan(2))
			},
		},
	}...))
}

func TestBeLessThan(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expected less than and was greater than",
			Act: func() {
				Expect(1).To(BeLessThan(2))
			},
		},
		{Scenario: "expected not less than and was equal",
			Act: func() {
				Expect(2).ToNot(BeLessThan(2))
			},
		},
	}...))
}
