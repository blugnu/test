package require_test

import (
	"fmt"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/require"
	"github.com/blugnu/test/test"
)

func TestExpectType(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "subsequent tests are evaluated when test passes",
			Act: func() {
				result := require.Type[int](1)
				Expect(result).To(Equal(1))

				Expect(false).To(BeTrue())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected true",
				)
			},
		},
		{Scenario: "subsequent tests are not evaluated when test fails",
			Act: func() {
				result := require.Type[int]("string")
				Expect(result).To(Equal(0))
				Expect(false).To(BeTrue())
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

func ExampleType() {
	test.Example()

	// Type returns the value as the expected type when it
	// is of that type
	var got any = 1 / 2.0

	result := require.Type[float64](got)

	fmt.Printf("result: type is: %T\n", result)
	fmt.Printf("result: value is: %v\n", result)

	// Type terminates the current test if the value is not
	// of the required type
	got = "1 / 2.0"

	require.Type[float64](got)

	Expect(false, "this will not be evaluated").To(BeTrue())

	//Output:
	// result: type is: float64
	// result: value is: 0.5
	//
	// expected type: float64
	// got          : string
}
