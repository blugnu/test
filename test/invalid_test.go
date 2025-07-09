package test_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestInvalid(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no test frame",
			Act: func() {
				defer Expect(Panic("INVALID TEST")).DidOccur()
				With(test.NilFrame())
				test.Invalid()
			},
		},
		{Scenario: "no message",
			Act: func() { test.Invalid() },
			Assert: func(result *R) {
				result.ExpectInvalid()
			},
		},
		{Scenario: "with string message",
			Act: func() { test.Invalid("invalid test") },
			Assert: func(result *R) {
				result.ExpectInvalid("invalid test")
			},
		},
		{Scenario: "with []string message",
			Act: func() {
				test.Invalid(
					"this is an invalid test with more information than can",
					"fit in a single line, so it is split into multiple lines",
					"to ensure that the message is still readable and informative",
				)
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"this is an invalid test with more information than can",
					"fit in a single line, so it is split into multiple lines",
					"to ensure that the message is still readable and informative",
				)
			},
		},
	}...))
}

func TestError(t *testing.T) {
	With(t)

	errTest := errors.New("test error")

	Run(HelperTests([]HelperScenario{
		{Scenario: "no test frame",
			Act: func() {
				defer Expect(Panic(errTest)).DidOccur() // FUTURE: DidOccur(WithString("with message")) if possible?
				With(test.NilFrame())
				test.Error(errTest, "with message")
			},
		},
		{Scenario: "no test frame, no message",
			Act: func() {
				defer Expect(Panic(errTest)).DidOccur()
				With(test.NilFrame())
				test.Error(errTest)
			},
		},
		{Scenario: "no message",
			Act: func() { test.Error(errTest) },
			Assert: func(result *R) {
				result.ExpectInvalid("ERROR: test error")
			},
		},
		{Scenario: "one-line message",
			Act: func() { test.Error(errTest, "with message") },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"ERROR: test error",
					"with message",
				)
			},
		},
		{Scenario: "multi-line message",
			Act: func() { test.Error(errTest, "with message line 1", "and line 2") },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"ERROR: test error",
					"with message line 1",
					"and line 2",
				)
			},
		},
	}...))
}

func TestWarning(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no test frame",
			Act: func() {
				defer Expect(Panic("WARNING: something you should know")).DidOccur() // FUTURE: DidOccur(WithString("with message")) if possible?
				With(test.NilFrame())
				test.Warning("something you should know")
			},
		},
		{Scenario: "with test frame",
			Act: func() { test.Warning("something you should know") },
			Assert: func(result *R) {
				result.ExpectWarning("something you should know")
			},
		},
	}...))
}

func ExampleError() {
	test.Example()

	// test.Error is used to mark a test as invalid with a specific error.
	// It will fail the current test with an error message that includes
	// the error and any additional message(s).
	ErrPreconditionsNotMet := errors.New("pre-conditions not met")
	test.Error(ErrPreconditionsNotMet, "explanation of the error")

	// Output:
	// <== INVALID TEST
	// ERROR: pre-conditions not met
	// explanation of the error
}
func ExampleInvalid() {
	test.Example()

	// test.Invalid is used to mark a test as invalid.
	// It will fail the current test with an error message that includes
	// the provided message(s).  If no message is provided, it will fail
	// the test as invalid without any additional information.

	test.Invalid("This test is invalid because of some reason")

	// Output:
	// <== INVALID TEST
	// This test is invalid because of some reason
}
