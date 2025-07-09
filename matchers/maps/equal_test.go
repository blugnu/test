package maps_test

import (
	"testing"

	. "github.com/blugnu/test"
)

type implementsEqual struct {
	isEqual bool
}

func (v implementsEqual) Equal(other implementsEqual) bool {
	return v.isEqual
}

func TestEqualMap(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expecting equal maps, got equal",
			Act: func() {
				w := map[string]int{"a": 1, "b": 2}
				g := map[string]int{"a": 1, "b": 2}
				Expect(g).To(EqualMap(w))
			},
		},
		{Scenario: "expecting a non-empty map to be equal to an empty map",
			Act: func() {
				w := map[string]int{}
				g := map[string]int{"a": 1}
				Expect(g).To(EqualMap(w))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got:",
					`  "a" => 1`,
				)
			},
		},
		{Scenario: "expecting an empty map to be equal to a non-empty map",
			Act: func() {
				w := map[string]int{"a": 1}
				g := map[string]int{}
				Expect(g).To(EqualMap(w))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected map:",
					`  "a" => 1`,
					"got: <empty map>",
				)
			},
		},
		{Scenario: "expecting with same keys but different values to be equal",
			Act: func() {
				w := map[string]int{"a": 10}
				g := map[string]int{"a": 1}
				Expect(g).To(EqualMap(w))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected map:",
					`  "a" => 10`,
					"got:",
					`  "a" => 1`,
				)
			},
		},
		{Scenario: "expecting with same keys but different values to NOT be equal",
			Act: func() {
				w := map[string]int{"a": 10, "b": 20}
				g := map[string]int{"a": 1, "b": 2}
				Expect(g).ToNot(EqualMap(w))
			},
		},
		{Scenario: "expecting equal maps to not be equal",
			Act: func() {
				w := map[string]int{"a": 1}
				g := map[string]int{"a": 1}
				Expect(g).ToNot(EqualMap(w))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map not equal to:",
					`  "a" => 1`,
				)
			},
		},
		{Scenario: "expecting empty map to not be empty",
			Act: func() {
				w := map[string]int{}
				g := map[string]int{}
				Expect(g).ToNot(EqualMap(w))
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected: <empty map>",
				)
			},
		},
		{Scenario: "value type implments Equal method",
			Act: func() {
				w := map[string]implementsEqual{"value": {isEqual: true}}
				g := map[string]implementsEqual{"value": {isEqual: true}}
				Expect(g).To(EqualMap(w))
			},
		},
		{Scenario: "using custom compare function (type-safe)",
			Act: func() {
				w := map[string]int{"value": 1}
				g := map[string]int{"value": 2}
				Expect(g).To(EqualMap(w), func(a, b int) bool { return true })
			},
		},
		{Scenario: "using custom compare function (any)",
			Act: func() {
				w := map[string]string{"value": "a"}
				g := map[string]string{"value": "b"}
				Expect(g).To(EqualMap(w), func(a, b any) bool { return true })
			},
		},
	}...))
}
