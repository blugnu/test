package test_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"
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
					"expected type: int",
					"got          : string",
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
					"  expected type: int",
					"  got          : bool",
				})
			},
		},
		{Scenario: "expecting error, got nil",
			Act: func() {
				result, ok := ExpectType[error](nil)
				Expect(ok).To(BeFalse())
				Expect(result).Should(BeNil())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: error",
					"got          : <nil>",
				)
			},
		},
		{Scenario: "expecting error, got error",
			Act: func() {
				result, ok := ExpectType[error](errors.New("an error occurred"))
				Expect(ok).To(BeTrue())
				Expect(result.Error()).To(Equal("an error occurred"))
			},
		},
	}...))
}

func TestRequireType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expecting int got int",
			Act: func() {
				result := RequireType[int](1)
				Expect(result).To(Equal(1))
			},
		},
		{Scenario: "expecting int got string",
			Act: func() {
				RequireType[int]("string")
				Expect(false, "this will not be evaluated").To(BeTrue())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: int",
					"got          : string",
				)
			},
		},
	}...))
}

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
	// expected type: float64
	// got          : string
}

func ExampleRequireType() {
	test.Example()

	// RequireType returns the value as the expected type when it
	// is of that type
	var got any = 1 / 2.0
	result := RequireType[float64](got)

	fmt.Printf("result: type is: %T\n", result)
	fmt.Printf("result: value is: %v\n", result)

	// RequireType terminates the current test if the value is not
	// of the required type
	got = "1 / 2.0"
	RequireType[float64](got)
	Expect(false, "this will not be evaluated").To(BeTrue())

	//Output:
	// result: type is: float64
	// result: value is: 0.5
	//
	// expected type: float64
	// got          : string
}
