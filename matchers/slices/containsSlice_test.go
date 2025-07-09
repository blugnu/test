package slices_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestContainsSlice(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: []string / To
		{Scenario: "expected slice of string to contain a slice that is present",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainSlice([]string{"a"}))
			},
		},
		{Scenario: "expected empty slice of string to contain a slice",
			Act: func() {
				Expect([]string{}).To(ContainSlice([]string{"a"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing slice:`,
					`| "a"`,
					`got: <empty slice>`,
				)
			},
		},
		{Scenario: "expected slice of string to contain a slice that is not present",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainSlice([]string{"c"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing slice:`,
					`| "c"`,
					`got:`,
					`| "a"`,
					`| "b"`,
				)
			},
		},

		// MARK: []int / To
		{Scenario: "expected slice of int to contain a slice that is present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainSlice([]int{1}))
			},
		},

		{Scenario: "expected slice is not present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainSlice([]int{3, 4}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int containing slice:`,
					`| 3`,
					`| 4`,
					`got:`,
					`| 1`,
					`| 2`,
				)
			},
		},
		{Scenario: "expected slice is only partially present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainSlice([]int{1, 2, 3}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int containing slice:`,
					`| 1`,
					`| 2`,
					`| 3`,
					`got:`,
					`| 1`,
					`| 2`,
				)
			},
		},
		{Scenario: "expected slice is present and preceded by partial match",
			Act: func() {
				s := []int{1, 2, 1, 2, 3}
				Expect(s).To(ContainSlice([]int{1, 2, 3}))
			},
		},
		{Scenario: "expected to not contain slice that is present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).ToNot(ContainSlice([]int{1}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int not containing slice:`,
					`| 1`,
					`got:`,
					`| 1`,
					`| 2`,
				)
			},
		},
		{Scenario: "using a custom type-safe comparison function",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainSlice([]int{3, 4}), func(a, b int) bool {
					return true
				})
			},
		},
		{Scenario: "using a custom any-based comparison function",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainSlice([]int{3, 4}), func(a, b any) bool {
					return true
				})
			},
		},
	}...))
}
