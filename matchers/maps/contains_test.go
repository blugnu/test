package maps_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestContainsMap(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "map contains expected map",
			Act: func() {
				m := map[string]int{"a": 1, "b": 2, "c": 3}
				s := map[string]int{"a": 1, "b": 2}
				Expect(m).To(ContainMap(s))
			},
		},
		{Scenario: "map does not contain expected map",
			Act: func() {
				m := map[string]int{"a": 1}
				s := map[string]int{"a": 1, "b": 2}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  "a" => 1`,
					`  "b" => 2`,
					"got:",
					`  "a" => 1`,
				)
			},
		},
		{Scenario: "map contains map with string values that differ in case",
			Act: func() {
				m := map[int]string{1: "ford"}
				s := map[int]string{1: "Ford"}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  1 => "Ford"`,
					"got:",
					`  1 => "ford"`,
				)
			},
		},
		{Scenario: "map contains map with string values that differ in case (case insensitive)",
			Act: func() {
				m := map[int]string{1: "ford", 2: "arthur"}
				s := map[int]string{1: "Ford", 2: "Arthur"}
				Expect(m).To(ContainMap(s), opt.CaseSensitive(false))
			},
		},
		{Scenario: "map contains map with []string values that differ in case",
			Act: func() {
				m := map[int][]string{1: {"ford", "arthur"}}
				s := map[int][]string{1: {"Ford", "arthur"}}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  1 => | "Ford"`,
					`       | "arthur"`,
					"got:",
					`  1 => | "ford"`,
					`       | "arthur"`,
				)
			},
		},
		{Scenario: "map contains map with []string values that differ in case (case insensitive)",
			Act: func() {
				m := map[int][]string{1: {"ford", "arthur"}}
				s := map[int][]string{1: {"Ford", "Arthur"}}
				Expect(m).To(ContainMap(s), opt.CaseSensitive(false))
			},
		},
		{Scenario: "map contains map with nil slice values",
			Act: func() {
				m := map[int][]int{1: nil}
				s := map[int][]int{1: nil}
				Expect(m).To(ContainMap(s))
			},
		},
		{Scenario: "map contains map with empty slice values",
			Act: func() {
				m := map[int][]int{1: {}}
				s := map[int][]int{1: {}}
				Expect(m).To(ContainMap(s))
			},
		},
		{Scenario: "expecting map with nil slice values to contain map with non-nil values",
			Act: func() {
				m := map[int][]int{1: nil}
				s := map[int][]int{1: {}}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  1 => <empty slice>`,
					"got:",
					`  1 => nil`,
				)
			},
		},
		{Scenario: "expecting map with non-nil slice values to contains map with nil values",
			Act: func() {
				m := map[int][]int{1: {}}
				s := map[int][]int{1: nil}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  1 => nil`,
					"got:",
					`  1 => <empty slice>`,
				)
			},
		},
		{Scenario: "map contains map with empty slice value but contains values with non-empty slice",
			Act: func() {
				m := map[int][]int{1: {}}
				s := map[int][]int{1: {1, 2}}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  1 => | 1`,
					`       | 2`,
					"got:",
					`  1 => <empty slice>`,
				)
			},
		},
		{Scenario: "expecting map with slice value in the wrong order",
			Act: func() {
				m := map[string][]int{"a": {1, 2}}
				s := map[string][]int{"a": {2, 1}}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  "a" => | 2`,
					`         | 1`,
					"got:",
					`  "a" => | 1`,
					`         | 2`,
				)
			},
		},
		{Scenario: "expecting map with slice values in any order that are equivalent",
			Act: func() {
				m := map[string][]int{"a": {1, 2}}
				s := map[string][]int{"a": {2, 1}}
				Expect(m).To(ContainMap(s), opt.AnyOrder())
			},
		},
		{Scenario: "expecting map with slice values in any order that are not equivalent",
			Act: func() {
				m := map[string][]int{"a": {1, 2}}
				s := map[string][]int{"a": {1, 1}}
				Expect(m).To(ContainMap(s), opt.AnyOrder())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  "a" => | 1`,
					`         | 1`,
					"got:",
					`  "a" => | 1`,
					`         | 2`,
				)
			},
		},
		{Scenario: "expecting map to contain a larger map",
			Act: func() {
				m := map[string]int{"a": 1, "b": 2}
				s := map[string]int{"a": 1, "b": 2, "c": 3}
				Expect(m).To(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map containing:",
					`  "a" => 1`,
					`  "b" => 2`,
					`  "c" => 3`,
					"got:",
					`  "a" => 1`,
					`  "b" => 2`,
				)
			},
		},
		{Scenario: "map not expected to contain map",
			Act: func() {
				m := map[string]int{"a": 1, "b": 2}
				s := map[string]int{"c": 1, "d": 2}
				Expect(m).ToNot(ContainMap(s))
			},
		},
		{Scenario: "map not expected to contain map, but does",
			Act: func() {
				m := map[string]int{"a": 1}
				s := map[string]int{"a": 1}
				Expect(m).ToNot(ContainMap(s))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: map not containing:",
					`  "a" => 1`,
				)
			},
		},

		{Scenario: "map with slice values contains expected map",
			Act: func() {
				m := map[string][]int{"a": {1, 2}, "b": {3, 4}}
				s := map[string][]int{"a": {1, 2}}
				Expect(m).To(ContainMap(s))
			},
		},
	}...))
}
