package slices_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestContainsItems(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: []string / To
		{Scenario: "expected to contain items that are present",
			Act: func() {
				Expect([]int{1, 2, 3}).
					To(ContainItems([]int{3, 1}))
			},
		},
		{Scenario: "expected to contain items that are not present",
			Act: func() {
				Expect([]int{1, 2, 3}).
					To(ContainItems([]int{4, 5}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int containing items:`,
					`| 4`,
					`| 5`,
					`got:`,
					`| 1`,
					`| 2`,
					`| 3`,
				)
			},
		},
		{Scenario: "expected to not contain items that are present",
			Act: func() {
				Expect([]int{1, 2, 3}).
					ToNot(ContainItems([]int{1, 2}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int not containing items:`,
					`| 1`,
					`| 2`,
					`got:`,
					`| 1`,
					`| 2`,
					`| 3`,
				)
			},
		},
		{Scenario: "expected to contain items that are duplicated",
			Act: func() {
				Expect([]int{1, 2, 1, 2}).
					To(ContainItems([]int{1, 2, 1, 2}))
			},
		},
		{Scenario: "using a custom type-safe comparison function",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItems([]int{3, 4}), func(a, b int) bool {
					return true
				})
			},
		},
		{Scenario: "using a custom any-based comparison function",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItems([]int{3, 4}), func(a, b any) bool {
					return true
				})
			},
		},
	}...))
}
