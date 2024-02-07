package test

import "testing"

func TestIsType(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T) (any, bool)
		assert   func(HelperTest, any, bool)
	}{
		{scenario: "IsType[string](string)",
			act: func(t T) (any, bool) {
				return IsType[string](t, "foo")
			},
			assert: func(test HelperTest, got any, ok bool) {
				test.DidPass()
				test.Report.IsEmpty()
				Equal(t, got, "foo")
				IsTrue(t, ok)
			},
		},
		{scenario: "IsType[string](int)",
			act: func(t T) (any, bool) {
				return IsType[string](t, 42)
			},
			assert: func(test HelperTest, got any, ok bool) {
				test.DidFail()
				test.Report.Contains("IsType[string](int)/is_of_type")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: string",
					"got   : int",
				})
				Equal(t, got, "")
				IsFalse(t, ok)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			var (
				got any
				ok  bool
			)

			//ACT
			test := Helper(t, func(st *testing.T) {
				got, ok = tc.act(st)
			})

			// ASSERT
			tc.assert(test, got, ok)
		})
	}
}
