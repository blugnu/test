package test

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// Equals tests
		{scenario: "Equals (same slices)",
			act: func(t T) {
				s := []string{"a", "b"}
				Slice(t, s).Equals(s)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Equals (equal slices)",
			act: func(t T) {
				Slice(t, []string{"a", "b"}).Equals([]string{"a", "b"})
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Equals (both nil)",
			act: func(t T) {
				Slice(t, []string(nil)).Equals([]string(nil))
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Equals (same values different order)",
			act: func(t T) {
				Slice(t, []string{"a", "b"}).Equals([]string{"b", "a"})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Equals_(same_values_different_order)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []string{\"b\", \"a\"}",
					"got   : []string{\"a\", \"b\"}",
				})
			},
		},
		{scenario: "Equals (ShallowEquality, default formatting)",
			act: func(t T) {
				type test struct{ int }
				a1 := &test{1}
				a2 := &test{2}
				b1 := &test{1}
				b2 := &test{2}
				a := []*test{a1, a2}
				b := []*test{b1, b2}
				Slice(t, a).Equals(b)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Equals_(ShallowEquality,_default_formatting)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []*test.test{(*test.test)",
					"got   : []*test.test{(*test.test)",
					"method: test.ShallowEquality",
				})
			},
		},
		{scenario: "Equals (ShallowEquality, FormatString)",
			act: func(t T) {
				type test struct{ int }
				a1 := &test{1}
				a2 := &test{2}
				b1 := &test{1}
				b2 := &test{2}
				a := []*test{a1, a2}
				b := []*test{b1, b2}
				Slice(t, a, FormatString).Equals(b)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Equals_(ShallowEquality,_FormatString)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [%!s(*test.test=&{1}) %!s(*test.test=&{2})]",
					"got   : [%!s(*test.test=&{1}) %!s(*test.test=&{2})]",
					"method: test.ShallowEquality",
				})
			},
		},
		{scenario: "Slice([]*int).Equals(test.DeepEquality) (not equal)",
			act: func(t T) {
				a := 42
				b := 24
				g := []*int{&a}
				w := []*int{&b}
				Slice(t, g).Equals(w, DeepEquality)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Slice([]*int).Equals(test.DeepEquality)_(not_equal)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []*int{(*int)",
					"got   : []*int{(*int)",
					"method: test.DeepEquality",
				})
			},
		},
		{scenario: "Slice([]*int, format func).Equals(test.DeepEquality) (not equal)",
			act: func(t T) {
				a := 42
				b := 24
				g := []*int{&a}
				w := []*int{&b}
				Slice(t, g, func(v []*int) string { return fmt.Sprintf("[%d]", *v[0]) }).Equals(w, DeepEquality)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Slice([]*int,_format_func).Equals(test.DeepEquality)_(not_equal)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [24]",
					"got   : [42]",
					"method: test.DeepEquality",
				})
			},
		},
		{scenario: "Slice([]*int).Equals(test.DeepEquality) (equal)",
			act: func(t T) {
				a := 42
				b := 42
				g := []*int{&a}
				w := []*int{&b}
				Slice(t, g).Equals(w, DeepEquality)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Slice([]{1,2,3}).Equals([]{2,4,6},custom{w=2*g})",
			act: func(t T) {
				g := []int{1, 2, 3}
				w := []int{2, 4, 6}
				Slice(t, g).Equals(w, func(g, w int) bool { return w == 2*g })
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Slice([]{1,2,3},\"doubled\").Equals([]{2,4,4},custom{w=2*g})",
			act: func(t T) {
				g := []int{1, 2, 3}
				w := []int{2, 4, 4}
				Slice(t, g, "doubled").Equals(w, func(g, w int) bool { return w == 2*g })
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Slice([]{1,2,3},\"doubled\").Equals([]{2,4,4},custom{w=2*g})/doubled/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []int{2, 4, 4}",
					"got   : []int{1, 2, 3}",
					"method: func(got, wanted) bool {...}",
				})
			},
		},
		{scenario: "Equals (different elements)",
			act: func(t T) {
				Slice(t, []string{"a", "b"}, FormatDecl).Equals([]string{"a", "b", "c"})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Equals_(different_elements)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []string{\"a\", \"b\", \"c\"}",
					"got   : []string{\"a\", \"b\"}",
				})
			},
		},

		// IsEmpty tests
		{scenario: "IsEmpty (is empty)",
			act: func(t T) {
				Slice(t, []string{}).IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "IsEmpty (not empty)",
			act: func(t T) {
				Slice(t, []string{"not empty"}).IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("IsEmpty_(not_empty)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: empty []string",
					"got   : []string{\"not empty\"}",
				})
			},
		},

		// IsNotEmpty tests
		{scenario: "IsNotEmpty (not empty)",
			act: func(t T) {
				Slice(t, []string{"not empty"}).IsNotEmpty()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "IsNotEmpty (empty)",
			act: func(t T) {
				Slice(t, []string{}).IsNotEmpty()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("IsNotEmpty_(empty)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: non-empty []string",
					"got   : empty",
				})
			},
		},

		// SlicesEqual tests
		// Equals tests
		{scenario: "SlicesEqual(got == wanted)",
			act: func(t T) {
				SlicesEqual(t, []string{"a", "b"}, []string{"a", "b"})
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "SlicesEqual(nil,nil)",
			act: func(t T) {
				SlicesEqual(t, []string(nil), []string(nil))
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "SlicesEqual([a,b],[b,a],FormatString)",
			act: func(t T) {
				SlicesEqual(t, []string{"a", "b"}, []string{"b", "a"}, FormatString)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("SlicesEqual([a,b],[b,a],FormatString)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [b a]",
					"got   : [a b]",
				})
			},
		},
		{scenario: "SlicesEqual(got != wanted, test.ShallowEquality, FormatDecl)",
			act: func(t T) {
				a := 42
				b := 42
				g := []*int{&a}
				w := []*int{&b}
				SlicesEqual(t, g, w, ShallowEquality, FormatDecl)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("SlicesEqual(got_!=_wanted,_test.ShallowEquality,_FormatDecl)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []*int{(*int)",
					"got   : []*int{(*int)",
					"method: test.ShallowEquality",
				})
			},
		},
		{scenario: "SlicesEqual(!=,test.ShallowEquality,format func)",
			act: func(t T) {
				a := 42
				b := 42
				g := []*int{&a}
				w := []*int{&b}
				SlicesEqual(t, g, w, ShallowEquality, func(v []*int) string { return fmt.Sprintf("[%p]", v[0]) })
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("SlicesEqual(!=,test.ShallowEquality,format_func)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [0x",
					"got   : [0x",
					"method: test.ShallowEquality",
				})
			},
		},
		{scenario: "SlicesEqual(got == wanted, test.DeepEquality)",
			act: func(t T) {
				a := 42
				b := 42
				g := []*int{&a}
				w := []*int{&b}
				SlicesEqual(t, g, w, DeepEquality)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "SlicesEqual([a,b],[a,b,c],FormatString)",
			act: func(t T) {
				SlicesEqual(t, []string{"a", "b"}, []string{"a", "b", "c"}, FormatString)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("SlicesEqual([a,b],[a,b,c],FormatString)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [a b c]",
					"got   : [a b]",
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
