package require_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/require"
)

func TestEqual(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Equal halts test when values are inequal",
			Act: func() {
				require.Equal(42, 43)
				Fail("not reached")
			},
			Assert: func(result *R) {
				result.Expect(
					`equal_test.go`,
					`  expected 43, got 42`,
				)
				Expect(result.Report).ToNot(ContainItem("not reached"))
			},
		},
	}...))
}
