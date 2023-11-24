package test

import "testing"

func TestMaps(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		name    string
		want    map[string]int
		got     map[string]int
		outcome HelperResult
	}{
		{name: "same keys and values",
			want:    map[string]int{"a": 1, "b": 2},
			got:     map[string]int{"a": 1, "b": 2},
			outcome: ShouldPass,
		},
		{name: "same keys different values",
			want:    map[string]int{"a": 1, "b": 2},
			got:     map[string]int{"a": 2, "b": 1},
			outcome: ShouldFail,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT & ASSERT
			Helper(t, func(st *testing.T) {
				Maps(st, tc.want, tc.got)
			}, tc.outcome)
		})
	}
}
