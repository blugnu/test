package test_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func TestHelperTests(t *testing.T) {
	With(t)

	Run(Test("no Act function defined", func() {
		result := TestHelper(func() {
			Run(HelperTests([]HelperScenario{{}}...))
		})

		result.ExpectInvalid("no Act function defined")
	}))

	Run(Test("result not checked", func() {
		result := TestHelper(func() {
			Run(HelperTests([]HelperScenario{{Act: func() {}, Assert: func(*R) {}}}...))
		})

		result.ExpectWarning("HelperTests() result not tested (scenario: testcase-001)")
	}))
}
