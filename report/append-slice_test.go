package report_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

func TestAppendArray(t *testing.T) {
	With(t)

	type testcase struct {
		Scenario string
		Exec     func()
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			tc.Exec()
		}),
		Cases([]testcase{
			{Scenario: "empty array",
				Exec: func() {
					result := report.AppendArray([]string{}, [0]string{}, opt.Name("title:"))
					Expect(result).To(EqualSlice([]string{"title: <empty array>"}))
				},
			},
			{Scenario: "non-empty array",
				Exec: func() {
					result := report.AppendArray([]string{}, [3]string{"a", "b", "c"}, opt.Name("title:"))
					Expect(result).To(EqualSlice([]string{
						`title: [ "a"`,
						`       [ "b"`,
						`       [ "c"`,
					}))
				},
			},
			{Scenario: "not an array",
				Exec: func() {
					result := TestHelper(func() {
						report.AppendArray([]string{}, "not an array", opt.Name("title:"))
					})
					result.ExpectInvalid(
						"report.AppendArray requires an array argument",
					)
				},
			},
		}),
	))
}

func TestAppendSlice(t *testing.T) {
	With(t)

	type testcase struct {
		Scenario string
		Exec     func()
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			tc.Exec()
		}),
		Cases([]testcase{
			{Scenario: "nil slice",
				Exec: func() {
					result := report.AppendSlice([]string{}, []string(nil), opt.Name("title:"))
					Expect(result).To(EqualSlice([]string{"title: <nil>"}))
				},
			},
			{Scenario: "empty slice",
				Exec: func() {
					result := report.AppendSlice([]string{}, []string{}, opt.Name("title:"))
					Expect(result).To(EqualSlice([]string{"title: <empty slice>"}))
				},
			},
			{Scenario: "non-empty slice",
				Exec: func() {
					result := report.AppendSlice([]string{}, []string{"a", "b", "c"}, opt.Name("title:"))
					Expect(result).To(EqualSlice([]string{
						`title: [ "a"`,
						`       [ "b"`,
						`       [ "c"`,
					}))
				},
			},
		}),
	))
}

func TestAppendSliceOrArray(t *testing.T) {
	With(t)

	Run(Test("with value that is not a slice or array", func() {
		result := TestHelper(func() {
			report.AppendSliceOrArray([]string{}, "not a slice or array", opt.Name("title:"))
		})
		result.ExpectInvalid(
			"report.AppendSliceOrArray requires a slice or array argument",
		)
	}))
}
