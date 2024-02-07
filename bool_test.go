package test

import (
	"testing"
)

func TestBool(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scneario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// these tests should pass
		{scneario: "Bool(true).IsTrue()",
			act: func(t T) {
				Bool(t, true, "got").IsTrue()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scneario: "Bool(false).IsFalse()",
			act: func(t T) {
				Bool(t, false, "got").IsFalse()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scneario: "Bool(false).Equals(false)",
			act: func(t T) {
				Bool(t, false, "got").Equals(false)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scneario: "Bool(true).Equals(true)",
			act: func(t T) {
				Bool(t, true, "got").Equals(true)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scneario: "IsFalse(false)",
			act: func(t T) {
				IsFalse(t, false, "got")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scneario: "IsTrue(true)",
			act: func(t T) {
				IsTrue(t, true, "got")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},

		// these should fail
		{scneario: "Bool(true).IsFalse()",
			act: func(t T) { Bool(t, true, "got").IsFalse() },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("FAIL: TestBool/Bool(true).IsFalse()/got/is_false")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: false",
					"got   : true",
				})
			},
		},
		{scneario: "Bool(false).IsTrue()",
			act: func(t T) { Bool(t, false, "got").IsTrue() },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("FAIL: TestBool/Bool(false).IsTrue()/got/is_true")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: true",
					"got   : false",
				})
			},
		},
		{scneario: "Bool(true).Equals(false)",
			act: func(t T) { Bool(t, true, "got").Equals(false) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("FAIL: TestBool/Bool(true).Equals(false)/got/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: false",
					"got   : true",
				})
			},
		},
		{scneario: "Bool(false).Equals(true)",
			act: func(t T) { Bool(t, false, "got").Equals(true) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("FAIL: TestBool/Bool(false).Equals(true)/got/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: true",
					"got   : false",
				})
			},
		},
		{scneario: "IsFalse(true)",
			act: func(t T) { IsFalse(t, true) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("IsFalse(true)/is_false")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: false",
					"got   : true",
				})
			},
		},
		{scneario: "IsTrue(false)",
			act: func(t T) { IsTrue(t, false) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("IsTrue(false)/is_true")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: true",
					"got   : false",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scneario, func(t T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}
