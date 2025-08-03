package typecheck_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
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
		{Scenario: "string is type based on string",
			Act: func() {
				type custom string

				Expect("hello").Should(BeOfType[custom]())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: typecheck_test.custom",
					"got          : string",
				)
			},
		},
		{Scenario: "type based on string is string",
			Act: func() {
				type custom string

				Expect(custom("hello")).Should(BeOfType[string]())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: string",
					"got          : typecheck_test.custom",
				)
			},
		},

		{Scenario: "interface implemented",
			Act: func() {
				var err = errors.New("an error")
				Expect(err).Should(BeOfType[error]())
			},
		},
		{Scenario: "interface not implemented",
			Act: func() {
				Expect(42).Should(BeOfType[error]())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: error",
					"got          : int",
				)
			},
		},
		{Scenario: "interface of wrong concrete type",
			Act: func() {
				var err error = errors.New("an error")
				Expect(err).Should(BeOfType[int]())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected type: int",
					"got          : *errors.errorString",
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
		{Scenario: "int is not int",
			Act: func() { Expect(42).ShouldNot(BeOfType[int]()) },
			Assert: func(result *R) {
				result.Expect("should not be of type: int")
			},
		},
		{Scenario: "int is not error",
			Act: func() { Expect(42).ShouldNot(BeOfType[error]()) },
		},
	}...))
}
