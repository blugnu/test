package slices_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestEqualsSlice(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: []string / To
		{Scenario: "expected equal slices to be equal",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(EqualSlice([]string{"a", "b"}))
			},
		},
		{Scenario: "expected equal slices to not be equal",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).ToNot(EqualSlice([]string{"a", "b"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string not equal to:`,
					`| "a"`,
					`| "b"`,
				)
			},
		},
		{Scenario: "expected empty slice to be equal to another empty slice",
			Act: func() {
				Expect([]string{}).To(EqualSlice([]string{}))
			},
		},
		{Scenario: "expected empty slice to be equal to non-empty slice",
			Act: func() {
				Expect([]int{}).To(EqualSlice([]int{42}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int equal to:`,
					`| 42`,
					`got: <empty slice>`,
				)
			},
		},
		{Scenario: "expected slice to be equal to another slice that is not equal",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(EqualSlice([]string{"c", "d"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string equal to:`,
					`| "c"`,
					`| "d"`,
					`got:`,
					`| "a"`,
					`| "b"`,
				)
			},
		},
		{Scenario: "same items in different order when order is significant",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(EqualSlice([]string{"b", "a"}), opt.ExactOrder(false))
			},
		},
		{Scenario: "custom comparison function (T)",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(EqualSlice([]string{"A", "B"}), func(a, b string) bool {
					return true
				})
			},
		},
		{Scenario: "custom comparison function (any)",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(EqualSlice([]string{"A", "B"}), func(a, b any) bool {
					return true
				})
			},
		},
	}...))
}
