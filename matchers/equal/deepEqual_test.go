package equal_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestDeepEqual(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "DeepEqual(struct)/equal",
			Act: func() {
				type foo struct {
					bytes []byte
				}
				Expect(foo{[]byte{1}}).To(DeepEqual(foo{[]byte{1}}))
			},
		},

		{Scenario: "DeepEqual(int)",
			Act: func() { Expect(1).To(DeepEqual(2)) },
			Assert: func(result *R) {
				result.Expect(
					"expected 2, got 1",
				)
			},
		},
		{Scenario: "DeepEqual(string)",
			Act: func() { Expect("the quick brown fox").To(DeepEqual("jumped over the lazy dog")) },
			Assert: func(result *R) {
				result.Expect(
					"expected: \"jumped over the lazy dog\"",
					"got     : \"the quick brown fox\"",
				)
			},
		},
		{Scenario: "DeepEqual(struct)",
			Act: func() {
				type foo struct {
					name string
				}
				Expect(foo{"ford"}).To(DeepEqual(foo{"arthur"}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: equal_test.foo{name:"arthur"}`,
					`got     : equal_test.foo{name:"ford"}`,
				)
			},
		},
		{Scenario: "DeepEqual(struct)/inequal",
			Act: func() {
				type foo struct {
					bytes []byte
				}
				Expect(foo{[]byte{65}}).To(DeepEqual(foo{[]byte{97}}))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: equal_test.foo{bytes:[]uint8{0x61}}`,
					`got     : equal_test.foo{bytes:[]uint8{0x41}}`,
				)
			},
		},
		{Scenario: "expect nil to deep equal nil",
			Act: func() { var a any; Expect(a).To(DeepEqual(any(nil))) },
		},
		{Scenario: "expect non-nil to deep equal nil",
			Act: func() { var a any = 42; Expect(a).To(DeepEqual(any(nil))) },
			Assert: func(result *R) {
				result.Expect(
					"expected nil, got 42",
				)
			},
		},

		{Scenario: "comparing values of type T implementing Equal(T)",
			Act: func() {
				Expect(equatable{false}).To(DeepEqual(equatable{true}))
			},
		},
		{Scenario: "with custom comparison func (any)",
			Act: func() {
				Expect(1).To(DeepEqual(2), func(expected, got any) bool {
					return true
				})
			},
		},
		{Scenario: "with custom comparison func (T)",
			Act: func() {
				Expect(1).To(DeepEqual(2), func(expected, got int) bool {
					return true
				})
			},
		},
	}...))
}
