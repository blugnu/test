package test

import (
	"math"
	"testing"

	"github.com/blugnu/test/opt"
)

type implementsCount[T int | int64 | uint | uint64] struct{ n T }

func (e implementsCount[T]) Count() T { return e.n }

type implementsLen[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLen[T]) Len() T { return e.n }

type implementsLength[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLength[T]) Length() T { return e.n }

// MARK: IsEmpty

func TestExpect_IsEmpty(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "nil",
			Act: func() { var a any; Expect(a).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : nil",
				)
			},
		},
		{Scenario: "nil slice",
			Act: func() { var s []int; Expect(s).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty slice>",
					"got     : nil slice",
				)
			},
		},
		{Scenario: "nil map",
			Act: func() { var m map[bool]struct{}; Expect(m).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got     : nil map",
				)
			},
		},
		{Scenario: "nil channel",
			Act: func() { var ch chan struct{}; Expect(ch).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty chan>",
					"got     : nil chan",
				)
			},
		},

		{Scenario: "empty string is empty",
			Act: func() { Expect("").IsEmpty() },
		},
		{Scenario: "empty array is empty",
			Act: func() { Expect([0]int{}).IsEmpty() },
		},
		{Scenario: "empty slice is empty",
			Act: func() { Expect([]int{}).IsEmpty() },
		},
		{Scenario: "empty channel is empty",
			Act: func() { Expect(make(chan struct{})).IsEmpty() },
		},
		{Scenario: "empty map is empty",
			Act: func() { Expect(map[string]struct{}{}).IsEmpty() },
		},
		{Scenario: "Count() int returns 0",
			Act: func() { Expect(implementsCount[int]{0}).IsEmpty() },
		},
		{Scenario: "Count() int64 returns 0",
			Act: func() { Expect(implementsCount[int64]{0}).IsEmpty() },
		},
		{Scenario: "Count() uint returns 0",
			Act: func() { Expect(implementsCount[uint]{0}).IsEmpty() },
		},
		{Scenario: "Count() uint64 returns 0",
			Act: func() { Expect(implementsCount[uint64]{0}).IsEmpty() },
		},
		{Scenario: "Len() int returns 0",
			Act: func() { Expect(implementsLen[int]{0}).IsEmpty() },
		},
		{Scenario: "Len() int64 returns 0",
			Act: func() { Expect(implementsLen[int64]{0}).IsEmpty() },
		},
		{Scenario: "Len() uint returns 0",
			Act: func() { Expect(implementsLen[uint]{0}).IsEmpty() },
		},
		{Scenario: "Len() uint64 returns 0",
			Act: func() { Expect(implementsLen[uint64]{0}).IsEmpty() },
		},
		{Scenario: "Length() int returns 0",
			Act: func() { Expect(implementsLength[int]{0}).IsEmpty() },
		},
		{Scenario: "Length() int64 returns 0",
			Act: func() { Expect(implementsLength[int64]{0}).IsEmpty() },
		},
		{Scenario: "Length() uint returns 0",
			Act: func() { Expect(implementsLength[uint]{0}).IsEmpty() },
		},
		{Scenario: "Length() uint64 returns 0",
			Act: func() { Expect(implementsLength[uint64]{0}).IsEmpty() },
		},
		{Scenario: "Length() uint64 returns > math.MaxInt",
			Act: func() { Expect(implementsLength[uint64]{math.MaxInt + 1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Length() == 9223372036854775808",
				)
			},
		},

		// non-empty
		{Scenario: "non-empty array",
			Act: func() { Expect([3]int{}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty array>",
					"got     : len() == 3",
				)
			},
		},
		{Scenario: "non-empty channel",
			Act: func() { ch := make(chan struct{}, 1); ch <- struct{}{}; Expect(ch).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty chan>",
					"got     : len() == 1",
				)
			},
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]struct{}{"set": {}}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got     : len() == 1",
				)
			},
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{0, 1, 2}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty slice>",
					"got     : len() == 3",
				)
			},
		},
		{Scenario: "non-empty string",
			Act: func() { Expect("abc").IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty string>",
					"got     : \"abc\"",
				)
			},
		},
		{Scenario: "non-empty []string",
			Act: func() { Expect([]string{"abc", "def"}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty []string>",
					`got: | "abc"`,
					`     | "def"`,
				)
			},
		},

		// Count/Len/Length int
		{Scenario: "non-empty Count() int",
			Act: func() { Expect(implementsCount[int]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Count() == 1",
				)
			},
		},
		{Scenario: "non-empty Len() int",
			Act: func() { Expect(implementsLen[int]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Len() == 1",
				)
			},
		},
		{Scenario: "non-empty Length() int",
			Act: func() { Expect(implementsLength[int]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Length() == 1",
				)
			},
		},

		// Count/Len/Length int64
		{Scenario: "non-empty Count() int64",
			Act: func() { Expect(implementsCount[int64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Count() == 1",
				)
			},
		},
		{Scenario: "non-empty Len() int64",
			Act: func() { Expect(implementsLen[int64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Len() == 1",
				)
			},
		},
		{Scenario: "non-empty Length() int64",
			Act: func() { Expect(implementsLength[int64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Length() == 1",
				)
			},
		},

		// Count/Len/Length uint
		{Scenario: "non-empty Count() uint",
			Act: func() { Expect(implementsCount[uint]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Count() == 1",
				)
			},
		},
		{Scenario: "non-empty Len() uint",
			Act: func() { Expect(implementsLen[uint]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Len() == 1",
				)
			},
		},
		{Scenario: "non-empty Length() uint",
			Act: func() { Expect(implementsLength[uint]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Length() == 1",
				)
			},
		},

		// Count/Len/Length uint64
		{Scenario: "non-empty Count() uint64",
			Act: func() { Expect(implementsCount[uint64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Count() == 1",
				)
			},
		},
		{Scenario: "non-empty Len() uint64",
			Act: func() { Expect(implementsLen[uint64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Len() == 1",
				)
			},
		},
		{Scenario: "non-empty Length() uint64",
			Act: func() { Expect(implementsLength[uint64]{1}).IsEmpty() },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty>",
					"got     : Length() == 1",
				)
			},
		},

		// invalid type
		{Scenario: "invalid type",
			Act: func() { Expect(1).IsEmpty() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"IsEmpty: requires a value that is a string, array, slice, channel or map,",
					"         or is of a type that implements a Count(), Len(), or Length()",
					"         function returning an int, int64, uint, or uint64.",
					"",
					"         A value of type int does not meet these criteria.",
				)
			},
		},

		// custom error report
		{Scenario: "custom error report",
			Act: func() {
				Expect([]int{1}).IsEmpty(opt.OnFailure("custom error report"))
			},
			Assert: func(result *R) {
				result.Expect("custom error report")
			},
		},
	})
}

func TestExpect_IsEmptyOrNil(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "nil",
			Act: func() { var a any; Expect(a).IsEmptyOrNil() },
		},
		{Scenario: "nil slice",
			Act: func() { var s []int; Expect(s).IsEmptyOrNil() },
		},
		{Scenario: "nil map",
			Act: func() { var m map[bool]struct{}; Expect(m).IsEmptyOrNil() },
		},
		{Scenario: "nil channel",
			Act: func() { var ch chan struct{}; Expect(ch).IsEmptyOrNil() },
		},
		{Scenario: "invalid type",
			Act: func() { var a struct{ value int }; Expect(a).IsEmptyOrNil() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"IsEmptyOrNil: requires a value that is a slice, channel, or map, or is of",
					"              a type that implements a Count(), Len(), or Length() function",
					"              returning an int, int64, uint, or uint64.",
					"",
					"              A value of type struct { value int } does not meet these criteria.",
				)
			},
		},
	})
}

// MARK: IsNotEmpty

func TestExpect_IsNotEmpty(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "empty string",
			Act: func() { Expect("").IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty string>")
			},
		},
		{Scenario: "empty array",
			Act: func() { Expect([0]int{}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty array>")
			},
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty slice>")
			},
		},
		{Scenario: "empty channel",
			Act: func() { Expect(make(chan struct{})).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty chan>")
			},
		},
		{Scenario: "empty map",
			Act: func() { Expect(map[string]struct{}{}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty map>")
			},
		},
		{Scenario: "empty []string]",
			Act: func() { Expect([]string{}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: <non-empty []string>")
			},
		},

		{Scenario: "Count() int returns 0",
			Act: func() { Expect(implementsCount[int]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Count() > 0")
			},
		},
		{Scenario: "Count() int64 returns 0",
			Act: func() { Expect(implementsCount[int64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Count() > 0")
			},
		},
		{Scenario: "Count() uint returns 0",
			Act: func() { Expect(implementsCount[uint]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Count() > 0")
			},
		},
		{Scenario: "Count() uint64 returns 0",
			Act: func() { Expect(implementsCount[uint64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Count() > 0")
			},
		},
		{Scenario: "Len() int returns 0",
			Act: func() { Expect(implementsLen[int]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Len() > 0")
			},
		},
		{Scenario: "Len() int64 returns 0",
			Act: func() { Expect(implementsLen[int64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Len() > 0")
			},
		},
		{Scenario: "Len() uint returns 0",
			Act: func() { Expect(implementsLen[uint]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Len() > 0")
			},
		},
		{Scenario: "Len() uint64 returns 0",
			Act: func() { Expect(implementsLen[uint64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Len() > 0")
			},
		},
		{Scenario: "Length() int returns 0",
			Act: func() { Expect(implementsLength[int]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Length() > 0")
			},
		},
		{Scenario: "Length() int64 returns 0",
			Act: func() { Expect(implementsLength[int64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Length() > 0")
			},
		},
		{Scenario: "Length() uint returns 0",
			Act: func() { Expect(implementsLength[uint]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Length() > 0")
			},
		},
		{Scenario: "Length() uint64 returns 0",
			Act: func() { Expect(implementsLength[uint64]{0}).IsNotEmpty() },
			Assert: func(result *R) {
				result.Expect("expected: Length() > 0")
			},
		},

		// non-empty
		{Scenario: "non-empty array",
			Act: func() { Expect([3]int{}).IsNotEmpty() },
		},
		{Scenario: "non-empty channel",
			Act: func() { ch := make(chan struct{}, 1); ch <- struct{}{}; Expect(ch).IsNotEmpty() },
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]struct{}{"set": {}}).IsNotEmpty() },
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{0, 1, 2}).IsNotEmpty() },
		},
		{Scenario: "non-empty string",
			Act: func() { Expect("abc").IsNotEmpty() },
		},
		{Scenario: "non-empty []string",
			Act: func() { Expect([]string{"abc", "defNot"}).IsNotEmpty() },
		},

		// Count/Len/Length int
		{Scenario: "non-empty Count() int",
			Act: func() { Expect(implementsCount[int]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Len() int",
			Act: func() { Expect(implementsLen[int]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Length() int",
			Act: func() { Expect(implementsLength[int]{1}).IsNotEmpty() },
		},

		// Count/Len/Length int64
		{Scenario: "non-empty Count() int64",
			Act: func() { Expect(implementsCount[int64]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Len() int64",
			Act: func() { Expect(implementsLen[int64]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Length() int64",
			Act: func() { Expect(implementsLength[int64]{1}).IsNotEmpty() },
		},

		// Count/Len/Length uint
		{Scenario: "non-empty Count() uint",
			Act: func() { Expect(implementsCount[uint]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Len() uint",
			Act: func() { Expect(implementsLen[uint]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Length() uint",
			Act: func() { Expect(implementsLength[uint]{1}).IsNotEmpty() },
		},

		// Count/Len/Length uint64
		{Scenario: "non-empty Count() uint64",
			Act: func() { Expect(implementsCount[uint64]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Len() uint64",
			Act: func() { Expect(implementsLen[uint64]{1}).IsNotEmpty() },
		},
		{Scenario: "non-empty Length() uint64",
			Act: func() { Expect(implementsLength[uint64]{1}).IsNotEmpty() },
		},

		// invalid type
		{Scenario: "invalid type",
			Act: func() { Expect(1).IsNotEmpty() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"IsNotEmpty: requires a value of type string, array, slice, channel or map,",
					"            or a type that implements a Count(), Len(), or Length() function",
					"            returning int, int64, uint, or uint64.",
					"",
					"            A value of type int does not meet these criteria.",
				)
			},
		},

		// custom error report
		{Scenario: "custom error report",
			Act: func() {
				Expect([]int{}).IsNotEmpty(opt.OnFailure("custom error report"))
			},
			Assert: func(result *R) {
				result.Expect("custom error report")
			},
		},
	})
}
