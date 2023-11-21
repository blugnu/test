package test

import "testing"

func TestEqual(t *testing.T) {
	// ARRANGE
	type result struct {
		outcome HelperResult
		wanted  string
		got     string
	}
	testcases := []struct {
		name string
		sut  func(*testing.T)
		result
	}{
		{name: "a == b", sut: func(st *testing.T) { Equal(st, 1, 1) }, result: result{outcome: ShouldPass}},
		{name: "a != b", sut: func(st *testing.T) { Equal(st, 1, 2) }, result: result{outcome: ShouldFail, wanted: "wanted: 1", got: "got   : 2"}},
		{name: "FormatHex", sut: func(st *testing.T) { Equal(st, 1, 255, FormatHex) }, result: result{outcome: ShouldFail, wanted: "wanted: 1", got: "got   : ff"}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			stdout, _ := Helper(t, func(st *testing.T) {
				tc.sut(st)
			}, tc.result.outcome)

			// ASSERT
			stdout.Contains(t, tc.result.wanted)
			stdout.Contains(t, tc.result.got)
		})
	}
}
