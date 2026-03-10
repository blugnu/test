package require_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/require"
)

func TestMap(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Equals halts test when maps are inequal",
			Act: func() {
				require.Map(map[int]int{1: 1}).Equals(map[int]int{1: 2})
				Fail("not reached")
			},
			Assert: func(result *R) {
				result.Expect(
					`map_test.go`,
					`expected map[int]int: ( 1 => 2`,
					`got: ( 1 => 1`,
				)
				Expect(result.Report).ToNot(ContainItem("not reached"))
			},
		},
	}...))
}
