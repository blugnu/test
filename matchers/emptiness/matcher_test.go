package emptiness_test

import (
	"testing"

	. "github.com/blugnu/test"
)

type implementsCount[T int | int64 | uint | uint64] struct{ n T }

func (e implementsCount[T]) Count() T { return e.n }

type implementsLen[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLen[T]) Len() T { return e.n }

type implementsLength[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLength[T]) Length() T { return e.n }

func TestBeEmpty(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "empty string",
			Act: func() { Expect("").Should(BeEmpty()) },
		},
		{Scenario: "empty array",
			Act: func() { Expect([0]int{}).Should(BeEmpty()) },
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).Should(BeEmpty()) },
		},
		{Scenario: "empty map",
			Act: func() { Expect(map[string]int{}).Should(BeEmpty()) },
		},
		{Scenario: "empty channel",
			Act: func() { Expect(make(chan int)).Should(BeEmpty()) },
		},
		{Scenario: "with Count() int == 0",
			Act: func() { Expect(implementsCount[int]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Count() int64 == 0",
			Act: func() { Expect(implementsCount[int64]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Count() uint == 0",
			Act: func() { Expect(implementsCount[uint]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Count() uint64 == 0",
			Act: func() { Expect(implementsCount[uint64]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Len() int == 0",
			Act: func() { Expect(implementsLen[int]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Len() int64 == 0",
			Act: func() { Expect(implementsLen[int64]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Len() uint == 0",
			Act: func() { Expect(implementsLen[uint]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Len() uint64 == 0",
			Act: func() { Expect(implementsLen[uint64]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Length() int == 0",
			Act: func() { Expect(implementsLength[int]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Length() int64 == 0",
			Act: func() { Expect(implementsLength[int64]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Length() uint == 0",
			Act: func() { Expect(implementsLength[uint]{n: 0}).Should(BeEmpty()) },
		},
		{Scenario: "with Length() uint64 == 0",
			Act: func() { Expect(implementsLength[uint64]{n: 0}).Should(BeEmpty()) },
		},

		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty slice>",
					"got     : <nil>",
				)
			},
		},
		{Scenario: "nil map",
			Act: func() { Expect(map[string]int(nil)).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got     : <nil>",
				)
			},
		},
		{Scenario: "nil channel",
			Act: func() { var ch chan struct{}; Expect(ch).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty chan>",
					"got     : <nil>",
				)
			},
		},

		{Scenario: "non-empty string",
			Act: func() { Expect("foo").Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty string>",
					"got     : \"foo\"",
				)
			},
		},
		{Scenario: "non-empty array",
			Act: func() { Expect([2]int{1, 2}).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty array>",
					"got     : [ 1",
					"          [ 2",
				)
			},
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{1, 2}).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty slice>",
					"got     : [ 1",
					"          [ 2",
				)
			},
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]int{"foo": 1}).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got     : len() == 1",
				)
			},
		},
		{Scenario: "with Count() > 0",
			Act: func() { Expect(implementsCount[int]{n: 1}).Should(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty emptiness_test.implementsCount[int]>",
					"got     : Count() == 1",
				)
			},
		},

		{Scenario: "nil",
			Act: func() { var a any; Expect(a).To(BeEmpty()) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"emptiness.Matcher: requires a value that is a slice, channel, or map, or of a type",
					"                   that implements a Count(), Len(), or Length() function returning",
					"                   int, int64, uint, or uint64.",
					"",
					"                   A value of type <nil> does not meet these criteria.",
				)
				Expect(result.Report).ToNot(ContainItem("expected"))
			},
		},
	}...))
}

func TestBeEmptyOrNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "empty string",
			Act: func() { Expect("").Should(BeEmptyOrNil()) },
		},
		{Scenario: "empty array",
			Act: func() { Expect([0]int{}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "empty map",
			Act: func() { Expect(map[string]int{}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "empty channel",
			Act: func() { Expect(make(chan int)).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Count() int == 0",
			Act: func() { Expect(implementsCount[int]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Count() int64 == 0",
			Act: func() { Expect(implementsCount[int64]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Count() uint == 0",
			Act: func() { Expect(implementsCount[uint]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Count() uint64 == 0",
			Act: func() { Expect(implementsCount[uint64]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Len() int == 0",
			Act: func() { Expect(implementsLen[int]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Len() int64 == 0",
			Act: func() { Expect(implementsLen[int64]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Len() uint == 0",
			Act: func() { Expect(implementsLen[uint]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Len() uint64 == 0",
			Act: func() { Expect(implementsLen[uint64]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Length() int == 0",
			Act: func() { Expect(implementsLength[int]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Length() int64 == 0",
			Act: func() { Expect(implementsLength[int64]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Length() uint == 0",
			Act: func() { Expect(implementsLength[uint]{n: 0}).Should(BeEmptyOrNil()) },
		},
		{Scenario: "with Length() uint64 == 0",
			Act: func() { Expect(implementsLength[uint64]{n: 0}).Should(BeEmptyOrNil()) },
		},

		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).Should(BeEmptyOrNil()) },
		},
		{Scenario: "nil map",
			Act: func() { Expect(map[string]int(nil)).Should(BeEmptyOrNil()) },
		},
		{Scenario: "nil channel",
			Act: func() { var ch chan struct{}; Expect(ch).Should(BeEmptyOrNil()) },
		},

		{Scenario: "non-empty string",
			Act: func() { Expect("foo").Should(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty string>",
					"got     : \"foo\"",
				)
			},
		},
		{Scenario: "non-empty array",
			Act: func() { Expect([2]int{1, 2}).Should(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty array>",
					"got     : [ 1",
					"          [ 2",
				)
			},
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{1, 2}).Should(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty slice>",
					"got     : [ 1",
					"          [ 2",
				)
			},
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]int{"foo": 1}).Should(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty map>",
					"got     : len() == 1",
				)
			},
		},
		{Scenario: "with Count() > 0",
			Act: func() { Expect(implementsCount[int]{n: 1}).Should(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <empty emptiness_test.implementsCount[int]>",
					"got     : Count() == 1",
				)
			},
		},

		{Scenario: "nil",
			Act: func() { var a any; Expect(a).To(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"emptiness.Matcher: requires a value that is a slice, channel, or map, or of a type",
					"                   that implements a Count(), Len(), or Length() function returning",
					"                   int, int64, uint, or uint64.",
					"",
					"                   A value of type <nil> does not meet these criteria.",
				)
			},
		},
	}...))
}

func TestNotBeEmpty(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "non-empty string",
			Act: func() { Expect("foo").ShouldNot(BeEmpty()) },
		},
		{Scenario: "non-empty array",
			Act: func() { Expect([2]int{1, 2}).ShouldNot(BeEmpty()) },
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{1, 2}).ShouldNot(BeEmpty()) },
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]int{"foo": 1}).ShouldNot(BeEmpty()) },
		},
		{Scenario: "with Count() > 0",
			Act: func() { Expect(implementsCount[int]{n: 1}).ShouldNot(BeEmpty()) },
		},

		{Scenario: "empty string",
			Act: func() { Expect("").ShouldNot(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "empty array",
			Act: func() { Expect([0]int{}).ShouldNot(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).ShouldNot(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "empty map",
			Act: func() { Expect(map[string]int{}).ShouldNot(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "with Count() == 0",
			Act: func() { Expect(implementsCount[int]{n: 0}).ShouldNot(BeEmpty()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty> (using Count() > 0)",
					"got     : Count() == 0",
				)
			},
		},
	}...))
}

func TestNotBeEmptyOrNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "non-empty string",
			Act: func() { Expect("foo").ShouldNot(BeEmptyOrNil()) },
		},
		{Scenario: "non-empty array",
			Act: func() { Expect([2]int{1, 2}).ShouldNot(BeEmptyOrNil()) },
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{1, 2}).ShouldNot(BeEmptyOrNil()) },
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]int{"foo": 1}).ShouldNot(BeEmptyOrNil()) },
		},
		{Scenario: "with Count() > 0",
			Act: func() { Expect(implementsCount[int]{n: 1}).ShouldNot(BeEmptyOrNil()) },
		},

		{Scenario: "empty string",
			Act: func() { Expect("").ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "empty array",
			Act: func() { Expect([0]int{}).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty>",
				)
			},
		},
		{Scenario: "empty slice",
			Act: func() { Expect([]int{}).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : <empty>",
				)
			},
		},
		{Scenario: "empty map",
			Act: func() { Expect(map[string]int{}).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : <empty>",
				)
			},
		},
		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : <nil>",
				)
			},
		},
		{Scenario: "nil map",
			Act: func() { Expect(map[string]int(nil)).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : <nil>",
				)
			},
		},
		{Scenario: "nil channel",
			Act: func() { var ch chan struct{}; Expect(ch).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : <nil>",
				)
			},
		},
		{Scenario: "with Count() == 0",
			Act: func() { Expect(implementsCount[int]{n: 0}).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil>",
					"got     : Count() == 0",
				)
			},
		},
		{Scenario: "with (nil).Count()",
			Act: func() { Expect((*implementsCount[int])(nil)).ShouldNot(BeEmptyOrNil()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: <not empty or nil> (using Count() > 0)",
					"got     : <nil>",
				)
			},
		},
	}...))
}
