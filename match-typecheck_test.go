package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestShould_BeOfType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "int is int",
			Act: func() { Expect(42).Should(BeOfType[int]()) },
		},
		{Scenario: "int is not string",
			Act: func() { Expect(42).Should(BeOfType[string]()) },
			Assert: func(result *R) {
				result.Expect(
					"expected type: string",
					"got          : int",
				)
			},
		},
	}...))
}

func TestShouldNot_BeOfType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "int is not string",
			Act: func() { Expect(42).ShouldNot(BeOfType[string]()) },
		},
		{Scenario: "int is string",
			Act: func() { Expect(42).ShouldNot(BeOfType[int]()) },
			Assert: func(result *R) {
				result.Expect(
					"should not be of type: int",
				)
			},
		},
	}...))
}

// =================================================================
// MARK: examples
// =================================================================

func ExampleBeOfType() {
	test.Example()

	var i int = 42

	Expect(i).Should(BeOfType[int]())       // passes
	Expect(i).ShouldNot(BeOfType[string]()) // passes

	Expect(i).Should(BeOfType[bool]()) // fails

	// Output:
	// expected type: bool
	// got          : int
}
