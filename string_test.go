package test

import "testing"

func TestString(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "String(string)",
			exec: func(t *testing.T) {
				// ACT
				result := String(t, "string")

				// ASSERT
				That(t, result.testable.got).Equals("string")
			},
		},
		{scenario: "String([]byte))",
			exec: func(t *testing.T) {
				// ACT
				result := String(t, []byte("bytes"))

				// ASSERT
				That(t, result.testable.got).Equals("bytes")
			},
		},
		{scenario: "Strings(int)",
			exec: func(t *testing.T) {
				// ACT
				result := String(t, 42)

				// ASSERT
				That(t, result.testable.got).Equals("42")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}

func TestStringTests(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// Contains
		{scenario: "String(abc).Contains(<empty string>)",
			act: func(t T) {
				String(t, "abc").Contains("")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).Contains(<empty_string>)/string/contains")
				test.Report.Contains([]string{
					`Contains(<empty string>) is invalid: did you mean IsEmpty()?`,
				})
			},
		},
		{scenario: "String(abc).Contains(a)",
			act: func(t T) {
				String(t, "abc").Contains("a")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "String(abc).Contains(d)",
			act: func(t T) {
				String(t, "abc").Contains("d")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).Contains(d)/string/contains")
				test.Report.Contains([]string{
					`wanted: string containing: "d"`,
					`got   : "abc"`,
				})
			},
		},

		// DoesNotContain
		{scenario: "String(abc).DoesNotContain(<empty string>)",
			act: func(t T) {
				String(t, "abc").DoesNotContain("")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).DoesNotContain(<empty_string>)/string/does_not_contain")
				test.Report.Contains([]string{
					`DoesNotContain(<empty string>) is invalid test: did you mean IsNotEmpty()?`,
				})
			},
		},

		{scenario: "String(abc).DoesNotContain(d)",
			act: func(t T) {
				String(t, "abc").DoesNotContain("d")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "String(abc).DoesNotContain(a)",
			act: func(t T) {
				String(t, "abc").DoesNotContain("a")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).DoesNotContain(a)/string/does_not_contain")
				test.Report.Contains([]string{
					`found: "a"`,
					`at   : 0`,
					`got  : "abc"`,
					`        ^`,
				})
			},
		},
		{scenario: "String(abc).DoesNotContain(b)",
			act: func(t T) {
				String(t, "abc").DoesNotContain("b")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).DoesNotContain(b)/string/does_not_contain")
				test.Report.Contains([]string{
					`found: "b"`,
					`at   : 1`,
					`got  : "abc"`,
					`         ^`,
				})
			},
		},
		{scenario: "String(abc).DoesNotContain(c)",
			act: func(t T) {
				String(t, "abc").DoesNotContain("c")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).DoesNotContain(c)/string/does_not_contain")
				test.Report.Contains([]string{
					`found: "c"`,
					`at   : 2`,
					`got  : "abc"`,
					`          ^`,
				})
			},
		},

		// Equals
		{scenario: "String(abc).Equals(abc)",
			act: func(t T) {
				String(t, "abc").Equals("abc")
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "String(abc).Equals(ABC)",
			act: func(t T) {
				String(t, "abc").Equals("ABC")
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).Equals(ABC)/string/equals")
				test.Report.Contains([]string{
					`wanted: "ABC"`,
					`got   : "abc"`,
				})
			},
		},

		// IsEmpty
		{scenario: "String(abc).IsEmpty()",
			act: func(t T) {
				String(t, "abc").IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String(abc).IsEmpty()/string/is_empty")
				test.Report.Contains([]string{
					`wanted: <empty string>`,
					`got   : "abc"`,
				})
			},
		},
		{scenario: "String().IsEmpty()",
			act: func(t T) {
				String(t, "").IsEmpty()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},

		// IsNotEmpty
		{scenario: "String(abc).IsNotEmpty()",
			act: func(t T) {
				String(t, "abc").IsNotEmpty()
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "String().IsNotEmpty()",
			act: func(t T) {
				String(t, "").IsNotEmpty()
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("String().IsNotEmpty()/string/is_not_empty")
				test.Report.Contains([]string{
					`wanted: <non-empty string>`,
					`got   : ""`,
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
