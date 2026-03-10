package expect_test

import (
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/test/expect"
)

func TestMap(t *testing.T) {
	With(t)

	// NOTE: several scenarios use only one key in the map under test to avoid
	//       tests failing due to key sorting differences in the failure report.

	Run(HelperTests([]HelperScenario{
		// MARK: DoesNotHaveKeys
		{Scenario: "DoesNotHaveKeys",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).DoesNotHaveKeys("c", "d")
			},
		},
		{Scenario: "DoesNotHaveKeys fails when map has one of the keys",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "d": 2}).DoesNotHaveKeys("a", "b")
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to not have keys: [ "a"`,
					`                               [ "b"`,
					`got: [ "a"`,
				)
			},
		},

		// MARK: DoesNotHaveValues
		{Scenario: "DoesNotHaveValues",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).DoesNotHaveValues(3, 4)
			},
		},
		{Scenario: "DoesNotHaveValues fails when map has one of the values",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).DoesNotHaveValues(2, 3)
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to not have values: [ 2`,
					`                                 [ 3`,
					`got: [ 2`,
				)
			},
		},

		// MARK: Equal
		{Scenario: "Equal",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).Equal(map[string]int{"a": 1, "b": 2})
			},
		},
		{Scenario: "Equal fails when maps are inequal",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).Equal(map[string]int{"a": 2})
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map[string]int: ( "a" => 2`,
					`got: ( "a" => 1`,
				)
			},
		},

		// MARK: HasAnyKey
		{Scenario: "HasAnyKey",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasAnyKey("b", "c")
			},
		},
		{Scenario: "HasAnyKey fails when map does not have any keys",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasAnyKey("b", "c")
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain at least one of these keys: [ "b"`,
					`                                                    [ "c"`,
					`got: [ "a"`,
				)
			},
		},

		// MARK: HasAnyValue
		{Scenario: "HasAnyValue",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasAnyValue(2, 3)
			},
		},
		{Scenario: "HasAnyValue fails when map does not have any values",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasAnyValue(2, 3)
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain at least one of these values: [ 2`,
					`                                                      [ 3`,
					`got: [ 1`,
				)
			},
		},

		// MARK: HasKey
		{Scenario: "HasKey",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasKey("a")
			},
		},
		{Scenario: "HasKey fails when map does not have key",
			// NOTE: only one key in the map under test to avoid the test failing
			//       due to key sorting differences in the failure report.
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasKey("b")
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain key: [ "b"`,
					`got keys: [ "a"`,
				)
			},
		},

		// MARK: HasKeys
		{Scenario: "HasKeys",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasKeys("a", "b")
			},
		},
		{Scenario: "HasKeys fails when map does not have all keys",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasKeys("a", "b")
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain keys: [ "a"`,
					`                              [ "b"`,
					`got keys: [ "a"`,
				)
			},
		},

		// MARK: HasLength
		{Scenario: "HasLength",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasLength(2)
			},
		},
		{Scenario: "HasLength fails when map length is incorrect",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasLength(2)
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					"expected map with: 2 keys",
					"got              : 1",
				)
			},
		},

		// MARK: HasValue
		{Scenario: "HasValue",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasValue(1)
			},
		},
		{Scenario: "HasValue fails when map does not have value",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasValue(2)
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain value: 2`,
					`got: ( "a" => 1`,
				)
			},
		},

		// MARK: HasValues
		{Scenario: "HasValues",
			Act: func() {
				expect.Map(map[string]int{"a": 1, "b": 2}).HasValues(1, 2)
			},
		},
		{Scenario: "HasValues fails when map does not have all values",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).HasValues(1, 2)
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					`expected map to contain values: [ 1`,
					`                                [ 2`,
					`got: ( "a" => 1`,
				)
			},
		},

		// MARK: IsEmpty
		{Scenario: "IsEmpty",
			Act: func() {
				expect.Map(map[string]int{}).IsEmpty()
			},
		},
		{Scenario: "IsEmpty fails when map is nil",
			Act: func() {
				var m map[string]int
				expect.Map(m).IsEmpty()
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					"expected: <empty map>",
					"got: <nil>",
				)
			},
		},
		{Scenario: "IsEmpty fails when map is not empty",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).IsEmpty()
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					"expected: <empty map>",
					`got: ( "a" => 1`,
				)
			},
		},

		// MARK: IsEmptyOrNil
		{Scenario: "IsEmptyOrNil with empty map",
			Act: func() {
				expect.Map(map[string]int{}).IsEmptyOrNil()
			},
		},
		{Scenario: "IsEmptyOrNil with nil map",
			Act: func() {
				var m map[string]int
				expect.Map(m).IsEmptyOrNil()
			},
		},
		{Scenario: "IsEmptyOrNil fails when map is not empty",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).IsEmptyOrNil()
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					"expected: <empty map> or <nil>",
					`got: ( "a" => 1`,
				)
			},
		},

		// MARK: IsNotEmpty
		{Scenario: "IsNotEmpty",
			Act: func() {
				expect.Map(map[string]int{"a": 1}).IsNotEmpty()
			},
		},
		{Scenario: "IsNotEmpty fails when map is empty",
			Act: func() {
				expect.Map(map[string]int{}).IsNotEmpty()
			},
			Assert: func(result *R) {
				result.Expect(
					"map_test.go",
					"expected map to not be empty",
					"got: <empty map>",
				)
			},
		},
	}...))
}
