package ordered_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestGreaterThan(t *testing.T) { //nolint:dupl  // incorrectly flags this entire test as a duplicate of TestLessThan (it is not)
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "expected greater than and was greater than",
			Act: func() {
				Expect(2).To(BeGreaterThan(1))
			},
		},
		{Scenario: "expected greater than and was equal",
			Act: func() {
				Expect(2).To(BeGreaterThan(2))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: greater than 2",
					"got     : 2",
				)
			},
		},
		{Scenario: "expected greater than and was less",
			Act: func() {
				Expect(1).To(BeGreaterThan(2))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: greater than 2",
					"got     : 1",
				)
			},
		},
		{Scenario: "expected not greater than and was greater than",
			Act: func() {
				Expect(2).ToNot(BeGreaterThan(1))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: not greater than 1",
					"got     : 2",
				)
			},
		},
		{Scenario: "custom comparison function",
			Act: func() {
				Expect(1).To(BeGreaterThan(2), func(a, b int) bool {
					return true
				})
			},
		},
	})
}

func TestLessThan(t *testing.T) { //nolint:dupl  // incorrectly flags this entire test as a duplicate of TestGreaterThan (it is not)
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "expected less than and was less than",
			Act: func() {
				Expect(1).To(BeLessThan(2))
			},
		},
		{Scenario: "expected less than and was equal",
			Act: func() {
				Expect(2).To(BeLessThan(2))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: less than 2",
					"got     : 2",
				)
			},
		},
		{Scenario: "expected less than and was greater",
			Act: func() {
				Expect(2).To(BeLessThan(1))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: less than 1",
					"got     : 2",
				)
			},
		},
		{Scenario: "expected not less than and was less than",
			Act: func() {
				Expect(1).ToNot(BeLessThan(2))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: not less than 2",
					"got     : 1",
				)
			},
		},
		{Scenario: "custom comparison function",
			Act: func() {
				Expect(2).To(BeLessThan(1), func(a, b int) bool {
					return true
				})
			},
		},
	})
}
