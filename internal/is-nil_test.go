package internal_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal"
)

func TestIsNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// MARK: nil cases

		{Scenario: "nil channel",
			Act: func() {
				var x chan int = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},
		{Scenario: "nil function",
			Act: func() {
				var x func() = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},
		{Scenario: "nil interface",
			Act: func() {
				var x any = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},
		{Scenario: "nil map",
			Act: func() {
				var x map[string]int = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},
		{Scenario: "nil pointer",
			Act: func() {
				var x *int = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},
		{Scenario: "nil slice",
			Act: func() {
				var x []int = nil
				Expect(internal.IsNil(x)).To(BeTrue())
			},
		},

		// MARK: non-nil cases

		{Scenario: "non-nil channel",
			Act: func() {
				x := make(chan int)
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
		{Scenario: "non-nil function",
			Act: func() {
				x := func() {}
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
		{Scenario: "non-nil interface",
			Act: func() {
				var x any = 42
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
		{Scenario: "non-nil map",
			Act: func() {
				x := map[string]int{}
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
		{Scenario: "non-nil pointer",
			Act: func() {
				var x *int = new(int)
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
		{Scenario: "non-nil slice",
			Act: func() {
				x := []int{}
				Expect(internal.IsNil(x)).To(BeFalse())
			},
		},
	}...))
}
