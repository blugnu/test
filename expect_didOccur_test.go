package test_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestExpect_DidOccur(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// panics
		{Scenario: "panic was expected and occurred",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
				panic(ErrInvalidArgument)
			},
		},
		{Scenario: "panic was expected and did not occur",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidOccur()
			},
			Assert: func(result *R) {
				result.Expect(TestFailed, opt.IgnoreReport(true))
			},
		},

		// errors
		{Scenario: "error was expected and occurred",
			Act: func() { Expect(errors.New("error")).DidOccur() },
		},
		{Scenario: "error was expected and did not occur",
			Act: func() {
				var err error
				Expect(err).DidOccur()
			},
			Assert: func(result *R) {
				result.Expect("expected error, got nil")
			},
		},

		// unsupported types
		{Scenario: "not an error or panic",
			Act: func() {
				Expect(42).DidOccur()
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.DidOccur: may only be used with Panic() or error values",
				)
			},
		},
	}...))
}

func TestExpect_DidNotOccur(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		// panics
		{Scenario: "panic(), no panic occurs",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidNotOccur()
			},
		},
		{Scenario: "panic(r), no panic occurs",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidNotOccur()
			},
		},
		{Scenario: "panic(nil)",
			Act: func() {
				defer Expect(Panic(nil)).DidNotOccur()
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"DidNotOccur: may not be used with Panic(nil); did you mean NilPanic()?",
				)
			},
		},
		{Scenario: "panic was expected to not occur",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidNotOccur()
				panic(ErrInvalidArgument)
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: panic with *errors.errorString(invalid argument): should not have occurred",
				)
			},
		},
		{Scenario: "panic was not expected and occurred",
			Act: func() {
				defer Expect(Panic(ErrInvalidArgument)).DidNotOccur()
				panic(ErrInvalidOperation)
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  recovered: *errors.errorString(invalid operation)",
				)
			},
		},

		// errors
		{Scenario: "error was not expected and did not occur",
			Act: func() { var err error = nil; Expect(err).DidNotOccur() },
		},
		{Scenario: "error was not expected and occurred",
			Act: func() {
				Expect(errors.New("error")).DidNotOccur()
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: <no error>",
					`got     : *errors.errorString(error)`,
				)
			},
		},

		// unsupported types
		{Scenario: "not an error or panic",
			Act: func() {
				Expect(42).DidNotOccur()
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.DidNotOccur: may only be used with Panic() or error values",
				)
			},
		},
	}...))
}
