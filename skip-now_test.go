package test_test

import (
	"testing"

	. "github.com/blugnu/test"
)

func Test_SkipNow(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{{
		Act: func() {
			Error("this is reported")
			SkipNow("this is logged as the skipped reason")

			Error("This line should not be reached")
		},
		Assert: func(result *R) {
			result.Expect(
				"this is reported",
				"<== SKIPPED: this is logged as the skipped reason",
			)
		},
	}}...))
}
