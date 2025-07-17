package ordered_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestBetween(t *testing.T) {
	With(t)

	Run(Test("Closed Interval (default)", func() {
		Run(HelperTests([]HelperScenario{
			{Scenario: "expected between and was between",
				Act: func() {
					Expect(2).To(BeBetween(1).And(3))
				},
			},
			{Scenario: "expected between and was equal to min",
				Act: func() {
					Expect(1).To(BeBetween(1).And(3))
				},
			},
			{Scenario: "expected between and was equal to max",
				Act: func() {
					Expect(3).To(BeBetween(1).And(3))
				},
			},
			{Scenario: "expected between and was less than min",
				Act: func() {
					Expect(0).To(BeBetween(1).And(3))
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between [1, 3]: 1 <= x <= 3",
						"got     : 0",
					)
				},
			},
			{Scenario: "expected between and was greater than max",
				Act: func() {
					Expect(4).To(BeBetween(1).And(3))
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between [1, 3]: 1 <= x <= 3",
						"got     : 4",
					)
				},
			},
		}...))
	}))

	Run(Test("Open Interval", func() {
		Run(HelperTests([]HelperScenario{
			{Scenario: "expected between and was between",
				Act: func() {
					Expect(2).To(BeBetween(1).And(3), opt.IntervalOpen)
				},
			},
			{Scenario: "expected between and was equal to min",
				Act: func() {
					Expect(1).To(BeBetween(1).And(3), opt.IntervalOpen)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3): 1 < x < 3",
						"got     : 1",
					)
				},
			},
			{Scenario: "expected between and was equal to max",
				Act: func() {
					Expect(3).To(BeBetween(1).And(3), opt.IntervalOpen)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3): 1 < x < 3",
						"got     : 3",
					)
				},
			},
			{Scenario: "expected between and was less than min",
				Act: func() {
					Expect(0).To(BeBetween(1).And(3), opt.IntervalOpen)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3): 1 < x < 3",
						"got     : 0",
					)
				},
			},
			{Scenario: "expected between and was greater than max",
				Act: func() {
					Expect(4).To(BeBetween(1).And(3), opt.IntervalOpen)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3): 1 < x < 3",
						"got     : 4",
					)
				},
			},
		}...))
	}))

	Run(Test("Open Min Interval", func() {
		Run(HelperTests([]HelperScenario{
			{Scenario: "expected between and was between",
				Act: func() {
					Expect(2).To(BeBetween(1).And(3), opt.IntervalOpenMin)
				},
			},
			{Scenario: "expected between and was equal to min",
				Act: func() {
					Expect(1).To(BeBetween(1).And(3), opt.IntervalOpenMin)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3]: 1 < x <= 3",
						"got     : 1",
					)
				},
			},
			{Scenario: "expected between and was equal to max",
				Act: func() {
					Expect(3).To(BeBetween(1).And(3), opt.IntervalOpenMin)
				},
			},
			{Scenario: "expected between and was less than min",
				Act: func() {
					Expect(0).To(BeBetween(1).And(3), opt.IntervalOpenMin)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3]: 1 < x <= 3",
						"got     : 0",
					)
				},
			},
			{Scenario: "expected between and was greater than max",
				Act: func() {
					Expect(4).To(BeBetween(1).And(3), opt.IntervalOpenMin)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between (1, 3]: 1 < x <= 3",
						"got     : 4",
					)
				},
			},
		}...))
	}))

	Run(Test("Open Max Interval", func() {
		Run(HelperTests([]HelperScenario{
			{Scenario: "expected between and was between",
				Act: func() {
					Expect(2).To(BeBetween(1).And(3), opt.IntervalOpenMax)
				},
			},
			{Scenario: "expected between and was equal to min",
				Act: func() {
					Expect(1).To(BeBetween(1).And(3), opt.IntervalOpenMax)
				},
			},
			{Scenario: "expected between and was equal to max",
				Act: func() {
					Expect(3).To(BeBetween(1).And(3), opt.IntervalOpenMax)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between [1, 3): 1 <= x < 3",
						"got     : 3",
					)
				},
			},
			{Scenario: "expected between and was less than min",
				Act: func() {
					Expect(0).To(BeBetween(1).And(3), opt.IntervalOpenMax)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between [1, 3): 1 <= x < 3",
						"got     : 0",
					)
				},
			},
			{Scenario: "expected between and was greater than max",
				Act: func() {
					Expect(4).To(BeBetween(1).And(3), opt.IntervalOpenMax)
				},
				Assert: func(result *R) {
					result.Expect(
						"expected: between [1, 3): 1 <= x < 3",
						"got     : 4",
					)
				},
			},
		}...))
	}))

	Run(Test("Unsupported Interval", func() {
		result := TestHelper(func() {
			Expect(2).To(BeBetween(1).And(3), opt.IntervalClosure(99))
		})

		result.ExpectInvalid("unsupported option: IntervalClosure(99)")
	}))

	Run(Test("Not Between", func() {
		result := TestHelper(func() {
			Expect(2).ToNot(BeBetween(1).And(3))
		})

		result.Expect(
			"not between [1, 3]: 1 <= x <= 3",
			"got     : 2",
		)
	}))

	Run(Test("Min > Max", func() {
		result := TestHelper(func() {
			Expect(2).ToNot(BeBetween(3).And(1))
		})

		result.Expect(
			"not between [1, 3]: 1 <= x <= 3",
			"got     : 2",
		)
	}))
}
