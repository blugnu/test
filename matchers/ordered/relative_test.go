package ordered_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/ordered"
)

func TestRelative(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: >

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

		// MARK: >=

		{Scenario: "expected greater than or equal to and was greater than",
			Act: func() {
				Expect(2).To(BeGreaterThan(1).OrEqual())
			},
		},
		{Scenario: "expected greater than or equal to and was equal",
			Act: func() {
				Expect(2).To(BeGreaterThan(2).OrEqual())
			},
		},
		{Scenario: "expected greater than or equal to and was less",
			Act: func() {
				Expect(1).To(BeGreaterThan(2).OrEqual())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: greater than or equal to 2",
					"got     : 1",
				)
			},
		},

		// MARK: <

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

		// MARK: <=

		{Scenario: "expected less than or equal to and was less than",
			Act: func() {
				Expect(1).To(BeLessThan(2).OrEqual())
			},
		},
		{Scenario: "expected less than or equal to and was equal",
			Act: func() {
				Expect(2).To(BeLessThan(2).OrEqual())
			},
		},
		{Scenario: "expected less than or equal to and was greater",
			Act: func() {
				Expect(2).To(BeLessThan(1).OrEqual())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: less than or equal to 1",
					"got     : 2",
				)
			},
		},

		// MARK: other

		{Scenario: "custom comparison function",
			Act: func() {
				match := ordered.RelativeMatcher[int]{
					Expected:   1,
					Comparison: ordered.LessThan,
				}
				Expect(2).To(match, func(a, b int) bool {
					return true
				})
			},
		},
		{Scenario: "unsupported comparison",
			Act: func() {
				match := ordered.RelativeMatcher[int]{
					Expected:   1,
					Comparison: ordered.Comparison(99),
				}
				Expect(1).To(match)
			},
			Assert: func(result *R) {
				result.ExpectInvalid("unknown comparison type: Comparison(99)")
			},
		},
	}...))
}
