package test_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestExpect_BeNil(t *testing.T) {
	With(t)

	// verify that the BeNil matcher factory function works as expected;
	// the matcher itself is covered by it's own tests

	Run(HelperTests([]HelperScenario{
		{Scenario: "should be nil",
			Act: func() { var subject any; Expect(subject).Should(BeNil()) },
		},
		{Scenario: "should not be nil",
			Act: func() { Expect(errors.New("error")).ShouldNot(BeNil()) },
		},
	}...))
}

func TestExpect_IsNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "nil",
			Act: func() { var subject any; Expect(subject).IsNil() },
		},
		{Scenario: "nil error",
			Act: func() { var err error; Expect(err).IsNil() },
		},

		{Scenario: "non-nil error",
			Act: func() { Expect(errors.New("error")).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got error: error")
			},
		},
		{Scenario: "non-nilable subject",
			Act: func() { Expect(0).IsNil() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'int' are not nilable",
				)
			},
		},
		{Scenario: "with custom failure report",
			Act: func() {
				Expect(errors.New("not nil")).IsNil(opt.FailureReport(func(a ...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect("custom failure report")
			},
		},
	}...))
}

func TestExpect_IsNotNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "non-nil error",
			Act: func() { Expect(errors.New("error")).IsNotNil() },
		},
		{Scenario: "non-nilable subject",
			Act: func() { Expect(0).IsNotNil() },
		},

		{Scenario: "nil",
			Act: func() { var subject any; Expect(subject).IsNotNil() },
			Assert: func(result *R) {
				result.Expect("expected not nil")
			},
		},
		{Scenario: "with custom failure report",
			Act: func() {
				var subject any
				Expect(subject).IsNotNil(opt.FailureReport(func(a ...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect("custom failure report")
			},
		},
	}...))
}

// ============================================================================
// MARK: examples
// ============================================================================

func Example_expectation_IsNil() {
	var (
		value int = 42
		ref   *int
	)

	// this test will pass
	Expect(ref).IsNil()

	// this test will fail
	ref = &value
	Expect(ref).IsNil()

	// Output:
	// expected nil, got 42 [*int]
}

func Example_expectation_IsNotNil() {
	var (
		value int = 42
		ref   *int
	)

	// this test will pass
	ref = &value
	Expect(ref, "non-nil ref").IsNotNil()

	// this test will fail
	ref = nil
	Expect(ref, "nil ref").IsNotNil()

	// Output:
	// nil ref:
	//   expected not nil
}
