package test

import (
	"testing"
)

func TestStrings(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "Strings(string)",
			exec: func(t *testing.T) {
				// ACT
				result := Strings(t, "string")

				// ASSERT
				That(t, result.testable.got).Equals([]string{"string"})
			},
		},
		{scenario: "Strings([]byte))",
			exec: func(t *testing.T) {
				// ACT
				result := Strings(t, []byte("bytes"))

				// ASSERT
				That(t, result.testable.got).Equals([]string{"bytes"})
			},
		},
		{scenario: "Strings([]string)",
			exec: func(t *testing.T) {
				// ACT
				result := Strings(t, []string{"a", "b"})

				// ASSERT
				That(t, result.testable.got).Equals([]string{"a", "b"})
			},
		},
		{scenario: "Strings(int)",
			exec: func(t *testing.T) {
				// ARRANGE
				defer ExpectPanic(ErrInvalidArgument).Assert(t)

				// ACT
				result := Strings(t, 42)

				// ASSERT
				That(t, result.testable.got).Equals([]string{"42"})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {

			tc.exec(t)
		})
	}
}

func TestStringsTests(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// DoesNotContain
		{scenario: "Strings([a,b]).DoesNotContain(c)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).DoesNotContain("c")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings(slice).Equals(slice)",
			act: func(t T) {
				s := []string{"a", "b"}
				Strings(t, s).Equals(s)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b]).Equals([a,b])",
			act: func(t T) {
				g := []string{"a", "b"}
				w := []string{"a", "b"}
				Strings(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b]).Equals([c,d])",
			act: func(t T) {
				g := []string{"a", "b"}
				w := []string{"c", "d"}
				Strings(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains([]string{
					`got:`,
					`  [0]: "a"`,
					`  [1]: "b"`,
					`--`,
					`[0] wanted: "c"`,
					`       got: "a"`,
					`[1] wanted: "d"`,
					`       got: "b"`,
				})
			},
		},
		{scenario: "Strings([a,b]).Equals([b,a])",
			act: func(t T) {
				g := []string{"a", "b"}
				w := []string{"b", "a"}
				Strings(t, g).Equals(w) // TODO: Equals(w, test.AnyOrder)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains([]string{
					`got:`,
					`  [0]: "a"`,
					`  [1]: "b"`,
					`--`,
					`[0] wanted: "b"`,
					`       got: "a"`,
					`[1] wanted: "a"`,
					`       got: "b"`,
				})
			},
		},
		{scenario: "Strings(a,b).Equals(a,b,c)",
			act: func(t T) {
				g := []string{"a", "b"}
				w := []string{"a", "b", "c"}
				Strings(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains([]string{
					`wanted: 3 elements`,
					`got   : 2`,
				})
			},
		},
		{scenario: "Strings(a,b).Contains(a)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).Contains("a")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings(a,b,c,d).Contains(b,c)",
			act: func(t T) {
				g := []string{"a", "b", "c", "d"}
				Strings(t, g).Contains([]string{"b", "c"})
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b,,c,d]).Contains([b,,c])",
			act: func(t T) {
				g := []string{"a", "b", "", "c", "d"}
				Strings(t, g).Contains([]string{"b", "", "c"})
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b,' ',c,d]).Contains([b,,c])",
			act: func(t T) {
				g := []string{"a", "b", " ", "c", "d"}
				Strings(t, g).Contains([]string{"b", "", "c"})
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b]).Contains([c,d])",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).Contains([]string{"c", "d"})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).Contains([c,d])/strings/contains")
				test.Report.Contains([]string{
					`wanted: [`,
					`c`,
					`d`,
					`]`,
					`got: [`,
					`a`,
					`b`,
					`]`,
				})
			},
		},
		{scenario: "Strings([a,b]).Contains(c)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).Contains("c")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).Contains(c)/strings/contains")
				test.Report.Contains([]string{
					`wanted: "c"`,
					`got   : [`,
					`a`,
					`b`,
					`]`,
				})
			},
		},
		{scenario: "Strings([a,b]).Contains([]{c})",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).Contains([]string{"c"})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).Contains([]{c})/strings/contains")
				test.Report.Contains([]string{
					`wanted: [`,
					`c`,
					`]`,
					`got: [`,
					`a`,
					`b`,
					`]`,
				})
			},
		},
		{scenario: "Strings([a,b]).Contains(nil)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).Contains(nil)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).Contains(nil)/strings/contains")
				test.Report.Contains("Contains(nil) is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
			},
		},
		{scenario: "Strings([]).Contains([])",
			act: func(t T) {
				Strings(t, []string{}).Contains([]string{})
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([]).Contains([])/strings/contains")
				test.Report.Contains("Contains(<empty slice>) is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
			},
		},
		{scenario: "Strings([]).Contains(\"\")",
			act: func(t T) {
				Strings(t, []string{}).Contains("")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([]).Contains(\"\")/strings/contains")
				test.Report.Contains("Contains(\"\") is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
			},
		},
		{scenario: "Strings([]).Contains(42)",
			act: func(t T) {
				Strings(t, []string{}).Contains(42)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([]).Contains(42)/strings/contains")
				test.Report.Contains("Contains(int) is invalid: Contains() accepts string or []string")
			},
		},
		{scenario: "Strings([]).Contains(nil)",
			act: func(t T) {
				Strings(t, []string{}).Contains(nil)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([]).Contains(nil)/strings/contains")
				test.Report.Contains("Contains(nil) is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
			},
		},
		{scenario: "Strings([a,b]).DoesNotContain(a)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).DoesNotContain("a")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).DoesNotContain(a)/strings/does_not_contain")
				test.Report.Contains([]string{
					currentFilename(),
					`wanted: does not contain: "a"`,
					`got   : [`,
					`a`,
					`b`,
					`]`,
				})
			},
		},
		{scenario: "Strings([a,b]).DoesNotContain(empty string)",
			act: func(t T) {
				g := []string{"a", "b"}
				Strings(t, g).DoesNotContain("")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b]).DoesNotContain(empty_string)/strings/does_not_contain")
				test.Report.Contains([]string{
					currentFilename(),
					"DoesNotContain() invalid test: specified string is empty or consists entirely of whitespace",
				})
			},
		},
		{scenario: "Strings([abc]).ContainsMatch(abc)",
			act: func(t T) {
				Strings(t, "abc").ContainsMatch("abc")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([abc]).ContainsMatch(^ab$)",
			act: func(t T) {
				Strings(t, "abc").ContainsMatch("^ab$")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([abc]).ContainsMatch(^ab$)/strings/contains_match")
				test.Report.Contains([]string{
					`wanted: contains match: ^ab$`,
					`got   : [`,
					`abc`,
					`]`,
				})
			},
		},
		{scenario: "Strings([abc]).ContainsMatch(invalid regex)",
			act: func(t T) {
				Strings(t, "abc").ContainsMatch("[")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([abc]).ContainsMatch(invalid_regex)/strings/contains_match")
				test.Report.Contains("invalid test: ContainsMatch([): error parsing regexp: missing closing ]: `[`")
			},
		},
		{scenario: "Strings([]).IsEmpty()",
			act: func(t T) {
				Strings(t, []string{}).IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Strings([a,b],\"got\").IsEmpty()",
			act: func(t T) {
				Strings(t, []string{"a", "b"}, "got").IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b],\"got\").IsEmpty()/got/is_empty")
				test.Report.Contains([]string{
					`wanted: <empty slice>`,
					`got   : [`,
					`a`,
					`b`,
					`]`,
				})
			},
		},
		{scenario: "Strings([]).IsNil()",
			act: func(t T) {
				Strings(t, []string{}).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([]).IsNil()/strings/is_nil")
				test.Report.Contains([]string{
					`wanted: nil`,
					`got   : <empty slice>`,
				})
			},
		},
		{scenario: "Strings([a,b,c]).IsNil()",
			act: func(t T) {
				Strings(t, []string{"a", "b", "c"}).IsNil()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Strings([a,b,c]).IsNil()/strings/is_nil")
				test.Report.Contains([]string{
					`wanted: nil`,
					`got   : [`,
					`a`,
					`b`,
					`c`,
					`]`,
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

func TestStringsTrimmed(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func() StringsTest
		assert   func(*testing.T, StringsTest)
	}{
		{scenario: "Strings([\" a \",\" b\",\"c \"]).Trimmed()",
			act: func() StringsTest {
				return Strings(t, []string{" a ", " b", "c "}).Trimmed()
			},
			assert: func(t *testing.T, st StringsTest) {
				SlicesEqual(t, []string{"a", "b", "c"}, st.got)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(t, tc.act())
		})
	}
}
