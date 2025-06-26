package test_test

import (
	"runtime"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestNilPanic(t *testing.T) {
	With(t)

	result := NilPanic()
	Expect(result).To(Equal(panics.Expected{R: &runtime.PanicNilError{}}))
}

func TestPanic(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "with no args",
			Act: func() {
				result := Panic()
				Expect(result).To(Equal(panics.Expected{}))
			},
		},
		{Scenario: "with nil arg",
			Act: func() {
				result := Panic(nil)
				Expect(result).To(Equal(panics.Expected{R: opt.NoPanicExpected(true)}))
			},
		},
		{Scenario: "with non-nil arg",
			Act: func() {
				result := Panic("some string")
				Expect(result).To(Equal(panics.Expected{R: "some string"}))
			},
		},
		{Scenario: "with multiple args",
			Act: func() {
				result := Panic("one", "two")
				Expect(result).To(Equal(panics.Expected{}))
			},
			Assert: func(result *R) {
				result.ExpectInvalid("Panic: expected at most one argument, got 2")
			},
		},
	})
}

func ExamplePanic() {
	test.Example()

	defer Expect(Panic("some string")).DidOccur()

	panic("some other string")

	// Output:
	// unexpected panic:
	//   expected : string("some string")
	//   recovered: string("some other string")
}
