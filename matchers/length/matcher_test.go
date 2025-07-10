package length_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestLength(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// passing tests
		{Scenario: "empty string",
			Act: func() { Expect("").Should(HaveLen(0)) },
		},
		{Scenario: "empty array",
			Act: func() {
				Expect([0]int{}).Should(HaveLen(0))
			},
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).Should(HaveLen(0)) },
		},
		{Scenario: "empty map",
			Act: func() {
				Expect(map[string]int{}).Should(HaveLen(0))
			},
		},
		{Scenario: "empty channel",
			Act: func() { Expect(make(chan int)).Should(HaveLen(0)) },
		},
		{Scenario: "non-empty string",
			Act: func() { Expect("abc").Should(HaveLen(3)) },
		},
		{Scenario: "non-empty array",
			Act: func() {
				Expect([3]int{1, 2, 3}).Should(HaveLen(3))
			},
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{1, 2, 3}).Should(HaveLen(3)) },
		},
		{Scenario: "non-empty map",
			Act: func() {
				Expect(map[string]int{"a": 1, "b": 2, "c": 3}).Should(HaveLen(3))
			},
		},
		{Scenario: "non-empty channel",
			Act: func() {
				ch := make(chan int, 3)
				ch <- 1
				ch <- 2
				ch <- 3
				Expect(ch).Should(HaveLen(3))
			},
		},

		// nil tests
		{Scenario: "nil any",
			Act: func() {
				var sut any
				Expect(sut).Should(HaveLen(0))
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"length.Matcher: requires a value that is a string, slice, channel, or map: got <nil>",
				)
			},
		},
		{Scenario: "nil slice",
			Act: func() {
				var sut []int
				Expect(sut).Should(HaveLen(0))
			},
		},
		{Scenario: "nil map",
			Act: func() {
				var sut map[string]int
				Expect(sut).Should(HaveLen(0))
			},
		},
		{Scenario: "nil channel",
			Act: func() {
				var sut chan int
				Expect(sut).Should(HaveLen(0))
			},
		},

		// failing tests on valid types
		{Scenario: "string of unexpected length",
			Act: func() { Expect("abc").Should(HaveLen(4)) },
			Assert: func(result *R) {
				result.Expect(
					"expected: len() == 4",
					"got     : len() == 3",
				)
			},
		},
		{Scenario: "array of unexpected length",
			Act: func() {
				Expect([3]int{1, 2, 3}).Should(HaveLen(4))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: len() == 4",
					"got     : len() == 3",
				)
			},
		},
		{Scenario: "slice of unexpected length",
			Act: func() { Expect([]int{1, 2, 3}).Should(HaveLen(4)) },
			Assert: func(result *R) {
				result.Expect(
					"expected: len() == 4",
					"got     : len() == 3",
				)
			},
		},
		{Scenario: "map of unexpected length",
			Act: func() {
				Expect(map[string]int{"a": 1, "b": 2, "c": 3}).Should(HaveLen(4))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: len() == 4",
					"got     : len() == 3",
				)
			},
		},

		// invalid type
		{Scenario: "int",
			Act: func() { Expect(42).Should(HaveLen(0)) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"length.Matcher: requires a value that is a string, slice, channel, or map: got int",
				)
			},
		},
	}...))
}
