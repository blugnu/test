package expect_test

import (
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/opt"
)

func TestEqual(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Equal",
			Act: func() {
				expect.Equal(42, 42)
			},
		},
		{Scenario: "Equal fails when values are inequal",
			Act: func() {
				expect.Equal(42, 43)
			},
			Assert: func(result *R) {
				result.Expect(
					`equal_test.go`,
					`  expected 43, got 42`,
				)
			},
		},
		{Scenario: "with a subject name",
			Act: func() {
				expect.Equal(42, 43, "answer")
			},
			Assert: func(result *R) {
				result.Expect(
					`equal_test.go`,
					`  answer:`,
					`    expected 43, got 42`,
				)
			},
		},
		{Scenario: "with custom failure message",
			Act: func() {
				expect.Equal(42, 43, opt.OnFailure("values are not equal"))
			},
			Assert: func(result *R) {
				result.Expect(
					`equal_test.go`,
					`  values are not equal`,
				)
			},
		},
	}...))
}
