package expect_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/test"
)

type Counter struct{}

func (Counter) Count() int {
	return 42
}

func TestExpectType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "expecting int got int",
			Act: func() {
				result, ok := expect.Type[int](1)
				Expect(result).To(Equal(1))
				Expect(ok).To(BeTrue())
			},
		},
		{Scenario: "expecting int got string",
			Act: func() {
				result, ok := expect.Type[int]("string")
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
				result, ok := expect.Type[int](false, "named value")
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
				result, ok := expect.Type[error](nil)
				Expect(ok).To(BeFalse())
				Expect(result).Should(BeNil())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: error",
					"got          : nil",
				)
			},
		},
		{Scenario: "expecting error, got error",
			Act: func() {
				result, ok := expect.Type[error](errors.New("an error occurred"))
				Expect(ok).To(BeTrue())
				Expect(result.Error()).To(Equal("an error occurred"))
			},
		},
		{Scenario: "expecting interface, got struct implementing interface",
			Act: func() {
				result, ok := expect.Type[interface{ Count() int }](Counter{})
				Expect(ok).To(BeTrue())
				Expect(result).IsNotNil()
			},
		},
		{Scenario: "expecting interface, got struct not implementing interface",
			Act: func() {
				result, ok := expect.Type[error](Counter{})
				Expect(ok).To(BeFalse())
				Expect(result).IsNil()
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: error",
					"got          : expect_test.Counter",
				)
			},
		},
	}...))
}

func ExampleType() {
	test.Example()

	// expect.Type returns the value as the expected type and true if the
	// value is of that type
	var got any = 1 / 2.0

	result, ok := expect.Type[float64](got)

	fmt.Printf("ok is %v\n", ok)
	fmt.Printf("result: type is: %T\n", result)
	fmt.Printf("result: value is: %v\n", result)

	// expect.Type returns the zero value of the expected type and false if the
	// value is not of that type (the return values can be ignored if the
	// test is only concerned with checking the type)
	got = "1 / 2.0"

	_, _ = expect.Type[float64](got)

	//Output:
	// ok is true
	// result: type is: float64
	// result: value is: 0.5
	//
	// expected type: float64
	// got          : string
}
