package require_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/require"
)

func TestSlice(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Equal halts test when slices are inequal",
			Act: func() {
				require.Slice([]int{1, 2, 3}).Equal(1, 2, 4)
				Fail("not reached")
			},
			Assert: func(result *R) {
				result.Expect(
					`slice_test.go`,
					`expected []int equal to: [ 1`,
					`                         [ 2`,
					`                         [ 4`,
					`got: [ 1`,
					`     [ 2`,
					`     [ 3`,
				)
				Expect(result.Report).ToNot(ContainItem("not reached"))
			},
		},
	}...))
}
