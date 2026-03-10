package expect_test

import (
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
)

func TestSlice(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Equals",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).Equals([]int{1, 2, 3})
			},
		},
		{Scenario: "Equals fails when slices are inequal",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).Equals([]int{1, 2, 4})
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected []int equal to: [ 1`,
					`                         [ 2`,
					`                         [ 4`,
					`got: [ 1`,
					`     [ 2`,
					`     [ 3`,
				)
			},
		},

		{Scenario: "Contains",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).Contains(1, 3)
			},
		},
		{Scenario: "Contains fails when slice does not contain all elements",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).Contains(2, 4)
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected []int containing items: [ 2`,
					`                                 [ 4`,
					`got: [ 1`,
					`     [ 2`,
					`     [ 3`,
				)
			},
		},

		{Scenario: "ContainsAnyOf",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).ContainsAny(3, 4)
			},
		},
		{Scenario: "ContainsAnyOf fails when slice does not contain any elements",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).ContainsAny(4, 5)
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected []int containing any of: [ 4`,
					`                                  [ 5`,
					`got: [ 1`,
					`     [ 2`,
					`     [ 3`,
				)
			},
		},

		{Scenario: "ContainsSlice",
			Act: func() {
				expect.Slice([]int{1, 2, 3, 4, 5}).ContainsSlice(2, 3, 4)
			},
		},
		{Scenario: "ContainsSlice fails when slice does not contain the sub-slice",
			Act: func() {
				expect.Slice([]int{1, 2, 3, 4, 5}).ContainsSlice(3, 5, 4)
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected []int containing slice: [ 3`,
					`                                 [ 5`,
					`                                 [ 4`,
					`got: [ 1`,
					`     [ 2`,
					`     [ 3`,
					`     [ 4`,
					`     [ 5`,
				)
			},
		},

		{Scenario: "HasLength",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).HasLength(3)
			},
		},
		{Scenario: "HasLength fails when slice length is incorrect",
			Act: func() {
				expect.Slice([]int{1, 2, 3}).HasLength(4)
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected 4, got 3`,
				)
			},
		},

		{Scenario: "IsEmpty",
			Act: func() {
				expect.Slice([]int{}).IsEmpty()
			},
		},
		{Scenario: "IsEmpty fails when slice is not empty",
			Act: func() {
				expect.Slice([]int{1}).IsEmpty()
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected: <empty slice>`,
				)
			},
		},

		{Scenario: "IsNotEmpty",
			Act: func() {
				expect.Slice([]int{1}).IsNotEmpty()
			},
		},
		{Scenario: "IsNotEmpty fails when slice is empty",
			Act: func() {
				expect.Slice([]int{}).IsNotEmpty()
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected: <not empty>`,
				)
			},
		},
	}...))
}
