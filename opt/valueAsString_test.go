package opt_test

import (
	"fmt"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestValueAsString(t *testing.T) {
	With(t)

	type testcase struct {
		value  any
		opts   []any
		result string
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// act
			result := opt.ValueAsString(tc.value, tc.opts...)

			// assert
			Expect(result).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{value: nil, result: "nil"},
			{value: "test 1", result: "\"test 1\""},
			{value: "test 2", opts: []any{opt.QuotedStrings(true)}, result: "\"test 2\""},
			{value: "test 3", opts: []any{opt.QuotedStrings(false)}, result: "test 3"},
			{value: "test 4", opts: []any{opt.UnquotedStrings()}, result: "test 4"},
			{value: []int{1, 2}, result: "[1 2]"},
			{value: []int{1, 2}, opts: []any{opt.AsDeclaration(true)}, result: "[]int{1, 2}"},
			{value: 42, result: "42"},
			{value: 42, opts: []any{opt.QuotedStrings(true)}, result: "42"},
		}),
	))
}

func ExampleValueAsString() {
	// non-string values are returned as unquoted strings
	fmt.Println(opt.ValueAsString(42))

	// string values are returned as quoted strings
	fmt.Println(opt.ValueAsString("example"))

	// to suppress the quotes, use opt.UnquotedString(true)
	fmt.Println(opt.ValueAsString("example", opt.UnquotedStrings()))

	// Output:
	// 42
	// "example"
	// example
}
