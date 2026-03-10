package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Error fails the test with the given message",
			Act: func() {
				Error("test error message")
			},
			Assert: func(result *R) {
				result.Expect("test error message")
			},
		},
		{Scenario: "Errorf fails the test with a formatted message",
			Act: func() {
				Errorf("test error message %d", 42)
			},
			Assert: func(result *R) {
				result.Expect("test error message 42")
			},
		},
	}...))
}

func TestFail(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "string message",
			Act: func() {
				Fail("it failed")
			},
			Assert: func(result *R) {
				result.Expect("it failed")
			},
		},
		{Scenario: "multi-line failure report",
			Act: func() {
				Fail([]string{"it failed", "on multiple lines"})
			},
			Assert: func(result *R) {
				result.Expect([]string{
					"it failed",
					"on multiple lines",
				})
			},
		},
		{Scenario: "no failure report",
			Act: func() {
				Fail()
			},
			Assert: func(result *R) {
				result.Expect("test failed")
			},
		},
		{Scenario: "fail if condition is true",
			Act: func() {
				FailIf(true, "condition was true")
			},
			Assert: func(result *R) {
				result.Expect("condition was true")
			},
		},
		{Scenario: "fail if condition is not true",
			Act: func() {
				FailIfNot(false, "condition was false")
			},
			Assert: func(result *R) {
				result.Expect("condition was false")
			},
		},

		{Scenario: "applies opt.Required",
			Act: func() {
				Fail(opt.Required(), "failed and is required")
				Fail("this should not be evaluated")
			},
			Assert: func(result *R) {
				result.Expect("failed and is required")
			},
		},
	}...))
}
