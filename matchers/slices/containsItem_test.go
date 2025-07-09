package slices_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestContainsItem(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: []string / To
		{Scenario: "expected slice of string to contain an item that is present",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainItem("a"))
			},
		},
		{Scenario: "expected slice of string to contain an item, case insensitive",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainItem("A"), opt.CaseSensitive(false))
			},
		},
		{Scenario: "expected empty slice of string to contain an item",
			Act: func() {
				Expect([]string{}).To(ContainItem("a"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing: "a"`,
					`got: <empty slice>`,
				)
			},
		},
		{Scenario: "expected slice of string to contain an item that differs in case",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainItem("A"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing: "A"`,
					`got:`,
					` | "a"`,
					` | "b"`,
				)
			},
		},
		{Scenario: "expected slice of string to contain an item that is not present",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainItem("c"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing: "c"`,
					`got:`,
					` | "a"`,
					` | "b"`,
				)
			},
		},
		{Scenario: "expected slice of string to contain an item that is not present, case insensitive",
			Act: func() {
				s := []string{"a", "b"}
				Expect(s).To(ContainItem("c"), opt.CaseSensitive(false))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []string containing: "c"`,
					`got:`,
					` | "a"`,
					` | "b"`,
					`(case insensitive comparison)`,
				)
			},
		},

		// MARK: []int / To
		{Scenario: "expected slice of int to contain an item that is present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItem(1))
			},
		},
		{Scenario: "expected empty slice of int to contain an item",
			Act: func() {
				Expect([]int{}).To(ContainItem(1))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int containing: 1`,
					`got: <empty slice>`,
				)
			},
		},
		{Scenario: "expected slice of int to contain an item that is not present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItem(3))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int containing: 3`,
					`got:`,
					`| 1`,
					`| 2`,
				)
			},
		},

		// MARK: []int / ToNot
		{Scenario: "expected slice of int to not contain an item that is not present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).ToNot(ContainItem(3))
			},
		},
		{Scenario: "expected slice of int to not contain an item that is present",
			Act: func() {
				s := []int{1, 2}
				Expect(s).ToNot(ContainItem(1))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: []int not containing: 1`,
					`got:`,
					`| 1`,
					`| 2`,
				)
			},
		},
		{Scenario: "using a custom compare function (type-safe)",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItem(3), func(a, b int) bool {
					return true
				})
			},
		},
		{Scenario: "using a custom compare function (any)",
			Act: func() {
				s := []int{1, 2}
				Expect(s).To(ContainItem(3), func(a, b any) bool {
					return true
				})
			},
		},
	}...))
}
