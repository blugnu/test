package report_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

func Test_AppendToReport(t *testing.T) {
	With(t)

	Run(Test("map is nil", func() {
		m := map[string]any(nil)
		result := report.AppendMap([]string{}, m, opt.Name("map:"))
		expected := []string{
			"map: <nil>",
		}
		Expect(result).To(EqualSlice(expected))
	}))

	Run(Test("map is empty", func() {
		m := map[string]any{}
		result := report.AppendMap([]string{}, m, opt.Name("map:"))
		expected := []string{
			"map: <empty map>",
		}
		Expect(result).To(EqualSlice(expected))
	}))

	Run(Test("map with 'any' values, quoted strings", func() {
		m := map[string]any{
			"key1": "value1",
			"key2": 42,
			"key3": []string{"item1", "item2"},
		}

		result := report.AppendMap([]string{}, m, opt.Name("map:"))
		expected := []string{
			`map: ( "key1" => "value1"`,
			`     ( "key2" => 42`,
			`     ( "key3" => [ "item1"`,
			`                 [ "item2"`,
		}

		Expect(result).To(EqualSlice(expected))
	}))

	Run(Test("map with 'any' values, plain strings", func() {
		m := map[string]any{
			"key1": "value1",
			"key2": 42,
			"key3": []string{"item1", "item2"},
		}

		result := report.AppendMap([]string{}, m, opt.UnquotedStrings, opt.Name("map:"))
		expected := []string{
			"map: ( key1 => value1",
			"     ( key2 => 42",
			"     ( key3 => [ item1",
			"               [ item2",
		}

		Expect(result).To(EqualSlice(expected))
	}))

	Run(Test("map with large slice", func() {
		m := map[string]any{
			"key1": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}

		result := report.AppendMap([]string{}, m, opt.MaxItems(4), opt.UnquotedStrings, opt.Name("map:"))
		expected := []string{
			"map: ( key1 => [ 1",
			"               [ 2",
			"               [ 3",
			"               [ 4",
			"               ... and 6 more",
		}

		Expect(result).To(EqualSlice(expected))
	}))
}
