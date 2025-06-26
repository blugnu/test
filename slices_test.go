package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestSlices(t *testing.T) {
	With(t)

	Run("ContainItem", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "slice contains expected item",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainItem("a"))
				},
			},
			{Scenario: "slice does not contain expected item",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainItem("c"))
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
			{Scenario: "slice contains item that matches using custom comparison function",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainItem("B"), opt.CaseSensitive(false))
				},
			},
		})
	})

	Run("ContainItems", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "slice contains expected items",
				Act: func() {
					Expect([]string{"a", "b", "c"}).To(ContainItems([]string{"a", "c"}))
				},
			},
			{Scenario: "slice does not contain expected items",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainItems([]string{"c", "d"}))
				},
				Assert: func(result *R) {
					result.Expect(
						`expected: []string containing items:`,
						`| "c"`,
						`| "d"`,
						`got:`,
						`| "a"`,
						`| "b"`,
					)
				},
			},
			{Scenario: "slice contains items that match using custom comparison function",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainItems([]string{"c", "d"}), func(a, b string) bool { return true })
				},
			},
		})
	})

	Run("ContainSlice", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "slice contains expected slice",
				Act: func() {
					Expect([]string{"a", "b", "c"}).To(ContainSlice([]string{"a", "b"}))
				},
			},
			{Scenario: "slice does not contain expected slice",
				Act: func() {
					Expect([]string{"a", "b", "c"}).To(ContainSlice([]string{"a", "c"}))
				},
				Assert: func(result *R) {
					result.Expect(
						`expected: []string containing slice:`,
						`| "a"`,
						`| "c"`,
						`got:`,
						`| "a"`,
						`| "b"`,
						`| "c"`,
					)
				},
			},
			{Scenario: "slice contains slice that matches using custom comparison function",
				Act: func() {
					Expect([]string{"a", "b"}).To(ContainSlice([]string{"c", "d"}), func(a, b string) bool { return true })
				},
			},
		})
	})

	Run("EqualSlice", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "slice is equal to another slice",
				Act: func() {
					Expect([]string{"a", "b"}).To(EqualSlice([]string{"a", "b"}))
				},
			},
			{Scenario: "slice is not equal to another slice",
				Act: func() {
					Expect([]string{"a", "b"}).To(EqualSlice([]string{"b", "a"}))
				},
				Assert: func(result *R) {
					result.Expect(
						`expected: []string equal to:`,
						`| "b"`,
						`| "a"`,
						`got:`,
						`| "a"`,
						`| "b"`,
					)
				},
			},
			{Scenario: "slice is equal to another slice that matches using custom comparison function",
				Act: func() {
					Expect([]string{"a", "b"}).To(EqualSlice([]string{"c", "d"}), func(a, b string) bool { return true })
				},
			},
			{Scenario: "slice is equal to another slice when order is not significant",
				Act: func() {
					Expect([]string{"a", "b"}).To(EqualSlice([]string{"b", "a"}), opt.ExactOrder(false))
				},
			},
		})
	})
}

func ExampleContainItem() {
	test.Example()

	sut := []string{"a", "b"}

	// these tests will pass
	Expect(sut).To(ContainItem("a"))
	Expect(sut).To(ContainItem("A"), opt.CaseSensitive(false))

	// this test will fail
	Expect(sut).To(ContainItem("c"))

	// Output:
	// expected: []string containing: "c"
	// got:
	// | "a"
	// | "b"
}
