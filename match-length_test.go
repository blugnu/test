package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestHaveLen(t *testing.T) {
	With(t)

	// HaveLen returns a matcher from the matchers/length package which is
	// thoroughly tested in that package; here we provide a simple test to
	// verify that the HaveLen factory function is wired up correctly.

	Run(HelperTests([]HelperScenario{
		{Scenario: "slice with items",
			Act: func() {
				sut := []int{1, 2}
				Expect(sut).Should(HaveLen(2))
			},
		},
		{Scenario: "nil slice",
			Act: func() {
				var sut []int
				Expect(sut).Should(HaveLen(0))
			},
		},
	}...))
}

// =================================================================
// MARK: examples
// =================================================================

func ExampleHaveLen() {
	test.Example()

	// these tests all pass:
	Expect([]int{}).Should(HaveLen(0))
	Expect(map[int]string{}).Should(HaveLen(0))
	Expect("").Should(HaveLen(0))
	Expect([0]int{}).Should(HaveLen(0))

	var ch chan int
	Expect(ch, "nil channel").Should(HaveLen(0))

	var ns []int
	Expect(ns, "nil slice").Should(HaveLen(0))

	var nm map[int]struct{}
	Expect(nm, "nil map").Should(HaveLen(0))

	// these tests all fail:
	Expect([]int{1, 2, 3}, "slice").Should(HaveLen(2))
	Expect(map[int]struct{}{1: {}, 2: {}, 3: {}}, "map").Should(HaveLen(2))
	Expect("abc", "string").Should(HaveLen(2))

	// Output:
	// slice:
	//   expected: len() == 2
	//   got     : len() == 3
	//
	// map:
	//   expected: len() == 2
	//   got     : len() == 3
	//
	// string:
	//   expected: len() == 2
	//   got     : len() == 3
}
