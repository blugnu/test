package test

import (
	"strings"
	"testing"
)

func TestValueTest(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// these tests should pass
		{scenario: "Value(1).DoesNotEqual(2)",
			act: func(t T) {
				Value(t, 1).DoesNotEqual(2)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Value(1).Equals(1)",
			act: func(t T) {
				Value(t, 1).Equals(1)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Value(1).Equals(2,custom{w=2*g})",
			act: func(t T) {
				Value(t, 1).Equals(2, func(g, w int) bool { return w == 2*g })
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Value[*int](nil).IsNil()",
			act: func(t T) {
				Value[*int](t, nil).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Value[*testing.T](non-nil).IsNotNil()",
			act: func(t T) {
				Value(t, t).IsNotNil()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},

		// these tests should fail
		{scenario: "Value(1).DoesNotEqual(1)",
			act: func(t T) { Value(t, 1).DoesNotEqual(1) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).DoesNotEqual(1)/value/does_not_equal")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: not: 1",
					"got        : 1",
				})
				test.Report.DoesNotContain("method: test.ShallowEquality")
			}},
		{scenario: "Value(1).DoesNotEqual(1,test.ShallowEquality)",
			act: func(t T) { Value(t, 1).DoesNotEqual(1, ShallowEquality) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).DoesNotEqual(1,test.ShallowEquality)/value/does_not_equal")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: not: 1",
					"got        : 1",
					"method     : test.ShallowEquality",
				})
			}},
		{scenario: "Value(1).DoesNotEqual(2,custom{w=2*g})",
			act: func(t T) { Value(t, 1).DoesNotEqual(2, func(g, w int) bool { return w == 2*g }) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).DoesNotEqual(2,custom{w=2*g})/value/does_not_equal")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: not: 2",
					"got        : 1",
					"method     : func(got, wanted) bool {...}",
				})
			}},
		{scenario: "Value(1).Equals(2)",
			act: func(t T) { Value(t, 1).Equals(2) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).Equals(2)/value/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 2",
					"got   : 1",
				})
				test.Report.DoesNotContain("method: test.ShallowEquality")
			}},
		{scenario: "Value(1).Equals(2,test.ShallowEquality)",
			act: func(t T) { Value(t, 1).Equals(2, ShallowEquality) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).Equals(2,test.ShallowEquality)/value/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 2",
					"got   : 1",
					"method: test.ShallowEquality",
				})
			}},
		{scenario: "Value(1).Equals(1,custom{w=g*2})",
			act: func(t T) { Value(t, 1).Equals(1, func(g, w int) bool { return w == g*2 }) },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).Equals(1,custom{w=g*2})/value/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: 1",
					"got   : 1",
					"method: func(got, wanted) bool {...}",
				})
			}},
		{scenario: "Value(1,format_func).Equals(2)",
			act: func(t T) {
				Value(t, 1, func(v int) string {
					s := "one"
					if v == 2 {
						s = "two"
					}
					return s
				}).Equals(2)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1,format_func).Equals(2)/value/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: two",
					"got   : one",
				})
			}},
		{scenario: "Value(1).IsNil()",
			act: func(t T) { Value(t, 1).IsNil() },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).IsNil()/value/is_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"invalid test",
				})
			}},
		{scenario: "Value(1).IsNotNil()",
			act: func(t T) { Value(t, 1).IsNotNil() },
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(1).IsNotNil()/value/is_not_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"invalid test",
				})
			}},
		{scenario: "Value(\"prefix:body\").Equals(\"body\",\"begins with\",strings.HasPrefix)",
			act: func(t T) {
				Value(t, "prefix:body").Equals("body", "begins with", strings.HasPrefix)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(\"prefix:body\").Equals(\"body\",\"begins_with\",strings.HasPrefix)/value/begins_with")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: body",
					"got   : prefix:body",
					"method: func(got, wanted) bool {...}",
				})
			},
		},
		{scenario: "Value(\"prefix:body\").DoesNotEqual(\"prefix:\",\"does not begin with\",func(got,wanted))",
			act: func(t T) {
				Value(t, "prefix:body").DoesNotEqual("prefix:", "does not begin with", func(got, wanted string) bool {
					return got != wanted
				})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Value(\"prefix:body\").DoesNotEqual(\"prefix:\",\"does_not_begin_with\",func(got,wanted))/value/does_not_begin_with")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: not: prefix",
					"got        : prefix:body",
					"method     : func(got, wanted) bool {...}",
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
