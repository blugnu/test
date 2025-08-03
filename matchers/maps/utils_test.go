package maps_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/maps"
)

func Test_appendToReport(t *testing.T) {
	With(t)

	Run(Test("appending map with 'any' values", func() {
		m := map[string]any{
			"key1": "value1",
			"key2": 42,
			"key3": []string{"item1", "item2"},
		}

		result := maps.AppendToReport([]string{}, "map:", m)
		expected := []string{
			"map:",
			`  "key1" => "value1"`,
			`  "key2" => 42`,
			`  "key3" => [item1 item2]`,
		}

		Expect(result).To(EqualSlice(expected))
	}))
}
