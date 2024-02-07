package test

import "testing"

func TestMaps(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenraio string
		act      func(T)
		assert   func(HelperTest)
	}{
		{scenraio: "same keys and values",
			act: func(t T) {
				w := map[string]int{"a": 1, "b": 2}
				g := map[string]int{"a": 1, "b": 2}
				Map(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenraio: "same keys different values",
			act: func(t T) {
				w := map[string]int{"a": 1, "b": 2}
				g := map[string]int{"a": 1, "b": 3}
				Map(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("same_keys_different_values")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: map[a:1 b:2]",
					"got   : map[a:1 b:3]",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenraio, func(t *testing.T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}
