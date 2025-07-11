package equal_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestEqual(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expected equal and was equal",
			Act: func() { Expect(1).To(Equal(1)) },
		},

		{Scenario: "Equal(int)",
			Act: func() { Expect(1).To(Equal(2)) },
			Assert: func(result *R) {
				result.Expect(
					"expected 2, got 1",
				)
			},
		},
		{Scenario: "Equal(string)",
			Act: func() { Expect("the quick brown fox").To(Equal("jumped over the lazy dog")) },
			Assert: func(result *R) {
				result.Expect(
					"expected: \"jumped over the lazy dog\"",
					"got     : \"the quick brown fox\"",
				)
			},
		},
		{Scenario: "Equal(struct)",
			Act: func() {
				type foo struct {
					name string
				}
				Expect(foo{"ford"}).To(Equal(foo{"arthur"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: equal_test.foo{name:"arthur"}`,
					`got     : equal_test.foo{name:"ford"}`,
				)
			},
		},
		{Scenario: "expected to not equal when not equal",
			Act: func() {
				Expect(1).ToNot(Equal(2))
			},
		},
		{Scenario: "expected to not equal when equal",
			Act: func() {
				Expect(1).ToNot(Equal(1))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected to not equal: 1",
				)
			},
		},
		{Scenario: "comparing values of type T implementing Equal(T)",
			Act: func() {
				Expect(equatable{false}).To(Equal(equatable{true}))
			},
		},
		{Scenario: "using custom comparison function",
			Act: func() {
				Expect(1).To(Equal(2), func(a, b int) bool { return true })
			},
		},
	}...))
}
