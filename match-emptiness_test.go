package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestBeEmpty(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Should(BeEmpty) with empty will pass",
			Act: func() {
				Expect([]int{}).Should(BeEmpty())
			},
		},
		{Scenario: "Should(BeEmpty) with non-empty will fail",
			Act: func() {
				Expect([]int{1}).Should(BeEmpty())
			},
			Assert: func(result *R) {
				result.Expect(test.Failed, opt.IgnoreReport(true))
			},
		},
		{Scenario: "Should(BeEmpty) with nil will fail",
			Act: func() {
				Expect([]int(nil)).Should(BeEmpty())
			},
			Assert: func(result *R) {
				result.Expect(test.Failed, opt.IgnoreReport(true))
			},
		},
		{Scenario: "ShouldNot(BeEmpty) with non-empty will pass",
			Act: func() {
				Expect([]int{1}).ShouldNot(BeEmpty())
			},
		},
		{Scenario: "ShouldNot(BeEmpty) with empty will fail",
			Act: func() {
				Expect([]int{}).ShouldNot(BeEmpty())
			},
			Assert: func(result *R) {
				result.Expect(test.Failed, opt.IgnoreReport(true))
			},
		},
		{Scenario: "ShouldNot(BeEmpty) with nil will pass",
			Act: func() {
				Expect([]int(nil)).ShouldNot(BeEmpty())
			},
		},
	}...))
}

func TestBeEmptyOrNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Should(BeEmptyOrNil) with nil",
			Act: func() {
				var nilSlice []int
				Expect(nilSlice).Should(BeEmptyOrNil())
			},
		},
	}...))
}

// =================================================================
// MARK: examples
// =================================================================

func ExampleBeEmpty() {
	test.Example()

	// these tests all pass:
	Expect([]int{}).Should(BeEmpty())
	Expect("").Should(BeEmpty())
	Expect(map[int]string{}).Should(BeEmpty())
	Expect([0]int{}).Should(BeEmpty())

	// these tests all fail:
	Expect([]int{1}).Should(BeEmpty())

	var ch chan int
	Expect(ch, "nil channel").Should(BeEmpty())

	var ns []int
	Expect(ns, "nil slice").Should(BeEmpty())

	var nm map[int]struct{}
	Expect(nm, "nil map").Should(BeEmpty())

	// Output:
	// expected: <empty slice>
	// got     : [ 1
	//
	// nil channel:
	//   expected: <empty chan>
	//   got     : <nil>
	//
	// nil slice:
	//   expected: <empty slice>
	//   got     : <nil>
	//
	// nil map:
	//   expected: <empty map>
	//   got     : <nil>
}

func ExampleBeEmptyOrNil() {
	test.Example()

	// these tests all pass:
	Expect([]int{}).Should(BeEmptyOrNil())
	Expect("").Should(BeEmptyOrNil())
	Expect(map[int]string{}).Should(BeEmptyOrNil())
	Expect([0]int{}).Should(BeEmptyOrNil())

	var ch chan int
	Expect(ch, "nil channel").Should(BeEmptyOrNil())

	var ns []int
	Expect(ns, "nil slice").Should(BeEmptyOrNil())

	var nm map[int]struct{}
	Expect(nm, "nil map").Should(BeEmptyOrNil())

	// this test will fail:
	Expect([]int{1}).Should(BeEmptyOrNil())

	// Output:
	// expected: <empty slice>
	// got     : [ 1
}
