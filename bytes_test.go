package test

import (
	"testing"

	"github.com/blugnu/test/opt"
)

func TestSliceOfBytes(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "when got and expected are equal",
			Act: func() { Expect([]byte{1, 2, 3}).To(EqualBytes([]byte{1, 2, 3})) },
		},
		{Scenario: "expected not equal and was not equal",
			Act: func() { Expect([]byte{1, 2, 3}).ToNot(EqualBytes([]byte{3, 2, 1})) },
		},
		{Scenario: "equal slices of custom byte type",
			Act: func() {
				type MyByte byte
				Expect([]MyByte{1, 2, 3}).To(EqualBytes([]MyByte{1, 2, 3}))
			},
		},

		// supported options
		{Scenario: "custom failure report",
			Act: func() {
				Expect([]byte{1}).To(EqualBytes([]byte{2}), opt.FailureReport(func(...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect(
					"custom failure report",
				)
			},
		},
	})
}

func ExampleEqualBytes() {
	With(ExampleTestRunner{})

	a := []byte{0x01, 0x02, 0x03}
	b := []byte{0x01, 0x03, 0x02}

	Expect(a).To(EqualBytes(b))

	// Output:
	// bytes not equal:
	//   differences at: [1, 2]
	// expected: 01 03 02
	//         |    ** **
	// got     : 01 02 03
}
