package test_test

import (
	"fmt"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/emptiness"
	"github.com/blugnu/test/matchers/matcher"
	"github.com/blugnu/test/test"
)

func TestExpectType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expecting int got int",
			Act: func() {
				result, ok := ExpectType[int](1)
				Expect(result).To(Equal(1))
				Expect(ok).To(BeTrue())
			},
		},
		{Scenario: "expecting int got string",
			Act: func() {
				result, ok := ExpectType[int]("string")
				Expect(result).To(Equal(0))
				Expect(ok).To(BeFalse())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: int",
					"got     : string",
				)
			},
		},
		{Scenario: "expecting named int got bool",
			Act: func() {
				result, ok := ExpectType[int](false, "named value")
				Expect(result).To(Equal(0))
				Expect(ok).To(BeFalse())
			},
			Assert: func(result *R) {
				result.Expect([]string{
					"named value:",
					"  expected: int",
					"  got     : bool",
				})
			},
		},
		{Scenario: "expecting interface implementation",
			Act: func() {
				ExpectType[matcher.ForAny](emptiness.Matcher{})
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"ExpectType: cannot be used to test for interfaces",
				)
			},
		},
	}...))
}

func ExampleExpectType() {
	test.Example()

	// ExpectType returns the value as the expected type and true if the
	// value is of that type
	var got any = 1 / 2.0
	result, ok := ExpectType[float64](got)

	fmt.Printf("ok is %v\n", ok)
	fmt.Printf("result: type is: %T\n", result)
	fmt.Printf("result: value is: %v\n", result)

	// ExpectType returns the zero value of the expected type and false if the
	// value is not of that type (the return values can be ignored if the
	// test is only concerned with checking the type)
	got = "1 / 2.0"
	ExpectType[float64](got)

	//Output:
	// ok is true
	// result: type is: float64
	// result: value is: 0.5
	//
	// expected: float64
	// got     : string
}
