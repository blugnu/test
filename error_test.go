package test_test

import (
	"testing"

	. "github.com/blugnu/test"
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
