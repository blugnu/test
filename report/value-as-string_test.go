package report_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type stringError string

func (ce stringError) Error() string {
	return string(ce)
}

type structError struct {
	s     string
	extra int
}

func (es structError) Error() string {
	return es.s
}

type prefixedError struct {
	s     string
	extra int
}

func (es prefixedError) Error() string {
	return "prefix: " + es.s
}

func byref[T any](v T) *T {
	return &v
}

func TestTypeAsString(t *testing.T) {
	With(t)

	Expect(report.TypeName("test")).To(Equal("string"))
	Expect(report.TypeName(42)).To(Equal("int"))
	Expect(report.TypeName[any](nil)).To(Equal("<nil>"))
	Expect(report.TypeName(struct{}{})).To(Equal("struct{}"))
	Expect(report.TypeName([]int{1, 2, 3})).To(Equal("[]int"))
	Expect(report.TypeName([]any{})).To(Equal("[]any"))
}

func TestValueAsString(t *testing.T) {
	With(t)

	type customFloat float64
	type customBool bool
	type customInt int
	type customString string

	type testcase struct {
		value  any
		opts   []any
		result string
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// act
			result := report.Value(tc.value, tc.opts...)

			// assert
			Expect(result).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{value: nil, result: "<nil>"},
			{value: "test 1", result: "\"test 1\""},
			{value: "test 2", opts: []any{opt.QuotedStrings}, result: "\"test 2\""},
			{value: "test 3", opts: []any{opt.UnquotedStrings}, result: "test 3"},
			{value: "test 4", opts: []any{opt.UnquotedStrings}, result: "test 4"},
			{value: [0]int{}, result: "[0]int (len=0)"},
			{value: [2]int{1, 2}, result: "[2]int (len=2) [1, 2]"},
			{value: []int{1, 2}, result: "[]int (len,cap=2) [1, 2]"},
			{value: []int{1, 2, 3, 4, 5}, opts: []any{opt.MaxItems(3)}, result: "[]int (len,cap=5) [1, 2, 3, ... and 2 more]"},
			{value: []int{1, 2}, opts: []any{opt.ValueAsDeclaration}, result: "[]int{1, 2}"}, // ValueAsDeclaration overrides the default formatting for slices
			{value: make([]int, 0), result: "[]int (len,cap=0)"},
			{value: make([]int, 0, 2), result: "[]int (len=0, cap=2)"},
			{value: make([]int, 1, 2), result: "[]int (len=1, cap=2) [0]"},
			{value: []string{"quoted 1", "quoted 2"}, result: "[]string (len,cap=2) [\"quoted 1\", \"quoted 2\"]"},
			{value: []string{"unquoted 1", "unquoted 2"}, opts: []any{opt.UnquotedStrings}, result: "[]string (len,cap=2) [\"unquoted 1\", \"unquoted 2\"]"}, // strings in array items are always quoted
			{value: []any{customString("custom string"), "string", 42}, result: "[]any (len,cap=3) [\"custom string\" [report_test.customString], \"string\", 42]"},
			{value: 42, result: "42"},
			{value: 42, opts: []any{opt.QuotedStrings}, result: "42"},
			{value: 42.1, result: "42.1"},
			{value: 42.1, opts: []any{opt.QuotedStrings}, result: "42.1"},
			{value: true, result: "true"},
			{value: true, opts: []any{opt.QuotedStrings}, result: "true"},
			{value: false, result: "false"},
			{value: false, opts: []any{opt.QuotedStrings}, result: "false"},
			{value: struct{}{}, result: "struct{}"},
			{value: func() bool { return true }, result: "func() bool"},
			{value: make(chan struct{}), result: "chan struct{} (len=0)"},
			{value: map[int]int{}, result: "map[int]int (len=0)"},
			{value: map[int]int{1: 1}, result: "map[int]int (len=1) [1 => 1]"},
			{value: []any{1, true, func() {}}, result: "[]any (len,cap=3) [1, true, func()]"},
			{value: testcase{}, result: "{value:<nil>, opts:[]any (len,cap=0), result:\"\"} [report_test.testcase]"},
			{value: &testcase{}, result: "{value:<nil>, opts:[]any (len,cap=0), result:\"\"} [*report_test.testcase]"},
			{value: byref(func() {}), result: "*func()"},
			{value: byref("string"), result: "\"string\" [*string]"},
			{value: (*string)(nil), result: "<nil> [*string]"},
			{value: (*string)(nil), opts: []any{opt.NoTypeNames}, result: "<nil>"},
			{value: (customString)("custom"), result: "\"custom\" [report_test.customString]"},
			{value: byref((customString)("custom")), result: "\"custom\" [*report_test.customString]"},
			{value: (customInt)(42), result: "42 [report_test.customInt]"},
			{value: byref((customInt)(42)), result: "42 [*report_test.customInt]"},
			{value: (customFloat)(3.14), result: "3.14 [report_test.customFloat]"},
			{value: byref((customFloat)(3.14)), result: "3.14 [*report_test.customFloat]"},
			{value: (customBool)(true), result: "true [report_test.customBool]"},
			{value: byref((customBool)(true)), result: "true [*report_test.customBool]"},
			{value: errors.New("this is an error"), result: `"this is an error" [*errors.errorString]`},
			{value: stringError("this is a custom error"), result: `"this is a custom error" [report_test.stringError]`},
			{value: structError{s: "this is a custom error", extra: 42}, result: `"this is a custom error" [report_test.structError]`},
			{value: prefixedError{s: "this is a custom error", extra: 42}, result: `"prefix: this is a custom error" [report_test.prefixedError]`},
			{opts: []any{opt.ValueAsDeclaration}, value: errors.New("this is an error"), result: `"this is an error" [*errors.errorString{s:"this is an error"}]`},
			{opts: []any{opt.ValueAsDeclaration}, value: stringError("this is a custom error"), result: `"this is a custom error" [report_test.stringError]`},
			{opts: []any{opt.ValueAsDeclaration}, value: structError{s: "this is a custom error", extra: 42}, result: `"this is a custom error" [report_test.structError{s:"this is a custom error", extra:42}]`},
			{opts: []any{opt.ValueAsDeclaration}, value: prefixedError{s: "this is a custom error", extra: 42}, result: `"prefix: this is a custom error" [report_test.prefixedError{s:"this is a custom error", extra:42}]`},
		}),
	))
}

func ExampleValueAsString() {
	// non-string values are returned as unquoted strings
	fmt.Println(report.Value(42))

	// string values are returned as quoted strings
	fmt.Println(report.Value("example"))

	// to suppress the quotes, use opt.PlainStrings
	fmt.Println(report.Value("example", opt.UnquotedStrings))

	// Output:
	// 42
	// "example"
	// example
}
