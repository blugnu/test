package slices_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/slices"
	"github.com/blugnu/test/opt"
)

func TestAppendToReport(t *testing.T) {
	With(t)

	type testcase struct {
		Scenario string
		Exec     func()
	}
	RunScenarios(
		func(tc *testcase, _ int) {
			tc.Exec()
		},
		[]testcase{
			{Scenario: "nil slice",
				Exec: func() {
					result := slices.AppendToReport([]string{}, []string(nil), "prefix:")
					Expect(result).To(EqualSlice([]string{"prefix: nil"}))
				},
			},
			{Scenario: "empty slice",
				Exec: func() {
					result := slices.AppendToReport([]string{}, []string{}, "prefix:")
					Expect(result).To(EqualSlice([]string{"prefix: <empty slice>"}))
				},
			},
			{Scenario: "non-empty slice",
				Exec: func() {
					result := slices.AppendToReport([]string{}, []string{"a", "b", "c"}, "prefix:")
					Expect(result).To(EqualSlice([]string{
						"prefix:",
						"| \"a\"",
						"| \"b\"",
						"| \"c\"",
					}))
				},
			},
			{Scenario: "non-empty slice (opt.PrefixInlineWithFirstItem)",
				Exec: func() {
					result := slices.AppendToReport([]string{}, []string{"a", "b", "c"}, "prefix:", opt.PrefixInlineWithFirstItem(true))
					Expect(result).To(EqualSlice([]string{
						"prefix: | \"a\"",
						"        | \"b\"",
						"        | \"c\"",
					}))
				},
			},
			{Scenario: "not a slice",
				Exec: func() {
					result := slices.AppendToReport([]string{}, "not a slice", "prefix:")
					Expect(result).To(EqualSlice([]string{"prefix: <not a slice>"}))
				},
			},
		},
	)
}
