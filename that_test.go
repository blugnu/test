package test

import "testing"

func TestThat(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// tests that should pass
		{
			scenario: "That(1).DeepEquals(1)",
			act: func(t T) {
				got := 1
				wanted := 1
				That(t, got).DeepEquals(wanted)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "That(1).Equals(2,()=>true)",
			act: func(t T) {
				That(t, 1).Equals(2, func(int, int) bool { return true })
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil any).IsNil()",
			act: func(t T) {
				That[any](t, nil).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil ptr).IsNil()",
			act: func(t T) {
				var g *int
				That(t, g).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil chan).IsNil()",
			act: func(t T) {
				var g chan int
				That(t, g).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil func).IsNil()",
			act: func(t T) {
				var g func()
				That(t, g).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil map).IsNil()",
			act: func(t T) {
				var g map[int]int
				That(t, g).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{
			scenario: "That(nil slice).IsNil()",
			act: func(t T) {
				var g []int
				That(t, g).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},

		// tests that should fail
		{
			scenario: "That(1).DeepEquals(2)",
			act: func(t T) {
				got := 1
				wanted := 2
				That(t, got).DeepEquals(wanted)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).DeepEquals(2)/got/deep_equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 2",
					"got   : 1",
				})
			},
		},
		{scenario: "That(1).Equals(1,()=>false)",
			act: func(t T) {
				That(t, 1).Equals(1, func(int, int) bool { return false })
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).Equals(1,()=>false)/got/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 1",
					"got   : 1",
					"method: <comparison func>",
				})
			},
		},
		{scenario: "That(1).Equals(2)",
			act: func(t T) {
				That(t, 1).Equals(2)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).Equals(2)/got/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 2",
					"got   : 1",
					"method: reflect.DeepEqual",
				})
			},
		},
		{scenario: "That(1).Equals(2,\"foo\")",
			act: func(t T) {
				That(t, 1).Equals(2, "foo")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).Equals(2,\"foo\")/got/foo")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 2",
					"got   : 1",
					"method: reflect.DeepEqual",
				})
			},
		},
		{
			scenario: "That(1).IsNil()",
			act: func(t T) {
				That(t, 1).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"invalid test: values of type 'int' are not nilable",
				})
			},
		},
		{
			scenario: "That(1).IsNotNil()",
			act: func(t T) {
				That(t, 1).IsNotNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(1).IsNotNil()/got/is_not_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"invalid test: values of type 'int' are not nilable",
				})
			},
		},
		{
			scenario: "That(non-nil ptr).IsNil()",
			act: func(t T) {
				That(t, new(int)).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(non-nil_ptr).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: nil (*int)",
					"got   : not nil",
				})
			},
		},
		{
			scenario: "That(non-nil chan).IsNil()",
			act: func(t T) {
				That(t, make(chan int)).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(non-nil_chan).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: nil (chan int)",
					"got   : not nil",
				})
			},
		},
		{
			scenario: "That(non-nil func).IsNil()",
			act: func(t T) {
				That(t, func() {}).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(non-nil_func).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: nil (func())",
					"got   : not nil",
				})
			},
		},
		{
			scenario: "That(non-nil map).IsNil()",
			act: func(t T) {
				That(t, make(map[int]int)).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(non-nil_map).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: nil (map[int]int)",
					"got   : not nil",
				})
			},
		},
		{
			scenario: "That(non-nil slice).IsNil()",
			act: func(t T) {
				That(t, make([]int, 0)).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(non-nil_slice).IsNil()/got/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: nil ([]int)",
					"got   : not nil",
				})
			},
		},
		{
			scenario: "That(nil slice).IsNotNil()",
			act: func(t T) {
				That(t, []int(nil)).IsNotNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("That(nil_slice).IsNotNil()/got/is_not_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: not nil ([]int)",
					"got   : nil",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}
