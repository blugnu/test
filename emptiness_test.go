package test_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestEmptiness(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "empty slice is empty",
			Act: func() {
				Expect([]int{}).Should(BeEmpty())
			},
		},
		{Scenario: "nil slice is empty",
			Act: func() {
				var nilSlice []int
				Expect(nilSlice).Should(BeEmptyOrNil())
			},
		},
		{Scenario: "nil slice has len 0",
			Act: func() {
				var nilSlice []int
				Expect(nilSlice).Should(HaveLen(0))
			},
		},
	}...))
}
