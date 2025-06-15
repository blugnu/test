package test

import (
	"testing"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

func TestKeysOfMap(t *testing.T) {
	With(t)

	sut := map[string]int{
		"ford":   42,
		"arthur": 23,
	}
	result := KeysOfMap(sut)
	Expect(result).To(EqualSlice([]string{"ford", "arthur"}), opt.AnyOrder())
}

func TestValuesOfMap(t *testing.T) {
	With(t)

	sut := map[string]int{
		"ford":   42,
		"arthur": 23,
	}
	result := ValuesOfMap(sut)
	Expect(result).To(EqualSlice([]int{23, 42}), opt.AnyOrder())
}

func TestContainMap(t *testing.T) {
	With(t)

	Expect(map[string]int{
		"ford":   42,
		"arthur": 23,
		"marvin": 2,
	}).To(ContainMap(map[string]int{
		"arthur": 23,
		"ford":   42,
	}))

	RunTestScenarios([]TestScenario{
		{Scenario: "with nil map",
			Act: func() {
				var nilmap map[string]int
				Expect(map[string]int{}).To(ContainMap(nilmap))
			},
			Assert: func(result *R) {
				result.Expect(
					"ContainMap() called with nil map",
					"Did you mean Expect(map).IsNil() or Expect(map).IsEmpty()?",
				)
			},
		},
		{Scenario: "with empty map",
			Act: func() {
				Expect(map[string]int{}).To(ContainMap(map[string]int{}))
			},
			Assert: func(result *R) {
				result.Expect(
					"ContainMap() called with empty map",
					"Did you mean Expect(map).To(EqualMap(<empty map>)) or Expect(map).IsEmpty()?",
				)
			},
		},
	})
}

func TestContainMapEntry(t *testing.T) {
	With(t)

	Expect(map[string]int{
		"ford":   42,
		"arthur": 23,
	}).To(ContainMapEntry("ford", 42))
}

func TestEqualMap(t *testing.T) {
	With(t)

	Expect(map[string]int{
		"ford":   42,
		"arthur": 23,
	}).To(EqualMap(map[string]int{
		"ford":   42,
		"arthur": 23,
	}))
}

func ExampleEqualMap() {
	test.Example()

	sut := map[string]int{
		"ford":   42,
		"arthur": 23,
	}

	// this test will pass
	Expect(sut).To(EqualMap(map[string]int{
		"ford":   42,
		"arthur": 23,
	}))

	// this test will fail
	sut = map[string]int{"marvin": 99}
	Expect(sut).To(EqualMap(map[string]int{"trillian": 24}))

	// Output:
	// expected map:
	//   "trillian" => 24
	// got:
	//   "marvin" => 99
}
