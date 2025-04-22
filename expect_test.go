package test

import (
	"errors"
	"testing"
)

func TestExpect(t *testing.T) {
	// for this test we deliberately fail to initialise a test frame
	// using With(t) so that we can test the panic handling in Expect().
	//
	// This means we must use the standard library testing.T for the
	// test.
	defer func() {
		// TODO: helper function for this, since we do it in multiple places
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok || !errors.Is(err, ErrNoTestFrame) {
				t.Errorf("expected ErrNoTestFrame, got %v", r)
			}
			return
		}
		t.Errorf("expected panic, got nil")
	}()

	Expect("any value will do")
}

func TestExpect_Error(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "no args/no name",
			Act: func() {
				Expect(any(nil)).Error()
			},
			Assert: func(result R) {
				Expect(result.Outcome).To(Equal(TestFailed))
			},
		},
		{Scenario: "no args/named",
			Act: func() {
				Expect(any(nil), "name").Error()
			},
			Assert: func(result R) {
				result.Assert("name failed")
			},
		},
		{Scenario: "string arg/no name",
			Act: func() {
				Expect(any(nil)).Error("failed with message", "other args ignored")
			},
			Assert: func(result R) {
				result.Assert("failed with message")
			},
		},
		{Scenario: "string arg/named",
			Act: func() {
				Expect(any(nil), "name").Error("failed with message", "other args ignored")
			},
			Assert: func(result R) {
				result.Assert("name: failed with message")
			},
		},
		{Scenario: "string slice arg/empty/no name",
			Act: func() {
				Expect(any(nil)).Error([]string{})
			},
			Assert: func(result R) {
				result.Assert("FAIL: TestExpect_Error/string_slice_arg/empty/no_name")
			},
		},
		{Scenario: "string slice arg/empty/named",
			Act: func() {
				Expect(any(nil), "name").Error([]string{})
			},
			Assert: func(result R) {
				result.Assert(
					"FAIL: TestExpect_Error/string_slice_arg/empty/named",
					currentFilename(),
					"name failed",
				)
			},
		},
		{Scenario: "string slice arg/no name",
			Act: func() {
				Expect(any(nil)).Error([]string{
					"failed with message",
					"and additional information",
				})
			},
			Assert: func(result R) {
				result.Assert("failed with message", "and additional information")
			},
		},
		{Scenario: "string slice arg/named",
			Act: func() {
				Expect(any(nil), "name").Error([]string{
					"failed with message",
					"and additional information",
				})
			},
			Assert: func(result R) {
				result.Assert("name: failed with message", "and additional information")
			},
		},
	})
}

func TestExpect_Errorf(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "no args/no name",
			Act: func() {
				Expect(any(nil)).Errorf("failed with message")
			},
			Assert: func(result R) {
				result.Assert("failed with message")
			},
		},
		{Scenario: "no args/named",
			Act: func() {
				Expect(any(nil), "name").Errorf("failed with message")
			},
			Assert: func(result R) {
				result.Assert("name: failed with message")
			},
		},
		{Scenario: "string arg/no name",
			Act: func() {
				Expect(any(nil)).Errorf("failed with message", "other args ignored")
			},
			Assert: func(result R) {
				result.Assert("failed with message")
			},
		},
		{Scenario: "string arg/named",
			Act: func() {
				Expect(any(nil), "name").Errorf("failed with message", "other args ignored")
			},
			Assert: func(result R) {
				result.Assert("name: failed with message")
			},
		},
		{Scenario: "args/no name",
			Act: func() {
				Expect(any(nil)).Errorf("failed with %s", "args")
			},
			Assert: func(result R) {
				result.Assert("failed with args")
			},
		},
		{Scenario: "args/named",
			Act: func() {
				Expect(any(nil), "name").Errorf("failed with %s", "args")
			},
			Assert: func(result R) {
				result.Assert("name: failed with args")
			},
		},
	})
}

func TestExpect_Is(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{
			Scenario: "nil is nil",
			Act:      func() { var a any; Expect(a).Is(nil) },
		},
		{
			Scenario: "non-nil is nil",
			Act:      func() { var a any = 1; Expect(a).Is(nil) },
			Assert: func(result R) {
				result.Assert("expected nil, got 1")
			},
		},
		{
			Scenario: "nil is non-nil",
			Act:      func() { var a any; Expect(a).Is(1) },
			Assert: func(result R) {
				result.Assert("expected 1, got nil")
			},
		},
		{
			Scenario: "nil error is nil",
			Act:      func() { var err error; Expect(err).Is(nil) },
		},
		{
			Scenario: "sentinel is sentinel",
			Act: func() {
				sent := errors.New("sentinel")
				Expect(sent).Is(sent)
			},
		},
		{
			Scenario: "sentinel is other sentinel",
			Act: func() {
				senta := errors.New("sentinel-a")
				sentb := errors.New("sentinel-b")
				Expect(senta).Is(sentb)
			},
			Assert: func(result R) {
				result.Assert([]string{
					"expected error: sentinel-b",
					"got           : sentinel-a",
				})
			},
		},
		{
			Scenario: "nil error nil vs error",
			Act: func() {
				var err error = errors.New("error")
				Expect(err).Is(nil)
			},
			Assert: func(result R) {
				result.Assert("expected nil, got error")
			},
		},
		{
			Scenario: "struct is equal struct",
			Act:      func() { Expect(struct{ a int }{a: 1}).Is(struct{ a int }{a: 1}) },
		},
		{
			Scenario: "struct is inequal struct",
			Act:      func() { Expect(struct{ a int }{a: 1}).Is(struct{ a int }{a: 2}) },
			Assert: func(result R) {
				result.Assert([]string{
					"expected: struct { a int }{a:2}",
					"got     : struct { a int }{a:1}",
				})
			},
		},
	})
}

type implementsCount[T int | int64 | uint | uint64] struct{ n T }

func (e implementsCount[T]) Count() T { return e.n }

type implementsLen[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLen[T]) Len() T { return e.n }

type implementsLength[T int | int64 | uint | uint64] struct{ n T }

func (e implementsLength[T]) Length() T { return e.n }

func TestExpect_IsEmpty(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
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

		// invalid type
		{Scenario: "invalid type",
			Act: func() { Expect(1).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"IsEmpty: invalid test for type int:",
					"  tested value must be string, array, slice, channel, map",
					"  or a type that implements Count(), Len() or Length()",
					"  returning int, int64, uint, or uint64",
				})
			},
		},

		// non-empty
		{Scenario: "non-empty array",
			Act: func() { Expect([3]int{}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty [3]int, got length 3",
				})
			},
		},
		{Scenario: "non-empty channel",
			Act: func() { ch := make(chan struct{}, 1); ch <- struct{}{}; Expect(ch).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty chan struct {}, got length 1",
				})
			},
		},
		{Scenario: "non-empty map",
			Act: func() { Expect(map[string]struct{}{"set": {}}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty map[string]struct {}, got length 1",
				})
			},
		},
		{Scenario: "non-empty slice",
			Act: func() { Expect([]int{0, 1, 2}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty []int, got length 3",
				})
			},
		},
		{Scenario: "non-empty string",
			Act: func() { Expect("abc").IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty string, got length 3",
					"  value: \"abc\"",
				})
			},
		},

		// Count/Len/Length int
		{Scenario: "non-empty Count() int",
			Act: func() { Expect(implementsCount[int]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsCount[int].Count() int, got length 1",
				})
			},
		},
		{Scenario: "non-empty Len() int",
			Act: func() { Expect(implementsLen[int]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLen[int].Len() int, got length 1",
				})
			},
		},
		{Scenario: "non-empty Length() int",
			Act: func() { Expect(implementsLength[int]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLength[int].Length() int, got length 1",
				})
			},
		},

		// Count/Len/Length int64
		{Scenario: "non-empty Count() int64",
			Act: func() { Expect(implementsCount[int64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsCount[int64].Count() int64, got length 1",
				})
			},
		},
		{Scenario: "non-empty Len() int64",
			Act: func() { Expect(implementsLen[int64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLen[int64].Len() int64, got length 1",
				})
			},
		},
		{Scenario: "non-empty Length() int64",
			Act: func() { Expect(implementsLength[int64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLength[int64].Length() int64, got length 1",
				})
			},
		},

		// Count/Len/Length uint
		{Scenario: "non-empty Count() uint",
			Act: func() { Expect(implementsCount[uint]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsCount[uint].Count() uint, got length 1",
				})
			},
		},
		{Scenario: "non-empty Len() uint",
			Act: func() { Expect(implementsLen[uint]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLen[uint].Len() uint, got length 1",
				})
			},
		},
		{Scenario: "non-empty Length() uint",
			Act: func() { Expect(implementsLength[uint]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLength[uint].Length() uint, got length 1",
				})
			},
		},

		// Count/Len/Length uint64
		{Scenario: "non-empty Count() uint64",
			Act: func() { Expect(implementsCount[uint64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsCount[uint64].Count() uint64, got length 1",
				})
			},
		},
		{Scenario: "non-empty Len() uint64",
			Act: func() { Expect(implementsLen[uint64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLen[uint64].Len() uint64, got length 1",
				})
			},
		},
		{Scenario: "non-empty Length() uint64",
			Act: func() { Expect(implementsLength[uint64]{1}).IsEmpty() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected empty test.implementsLength[uint64].Length() uint64, got length 1",
				})
			},
		},
	})
}

func TestExpect_IsNil(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "nil is nil",
			Act: func() { Expect((any)(nil)).IsNil() },
		},
		{Scenario: "int is nil",
			Act: func() { Expect(0).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"test.IsNil: invalid test:",
					"  values of type 'int' are not nilable",
				})
			},
		},
		{Scenario: "struct is nil",
			Act: func() { Expect(struct{}{}).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"test.IsNil: invalid test:",
					"  values of type 'struct {}' are not nilable",
				})
			},
		},
		{Scenario: "error is nil",
			Act: func() { Expect(errors.New("some error message")).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected nil, got error: some error message",
				})
			},
		},
		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).IsNil() },
		},
		{Scenario: "non-nil slice",
			Act: func() { Expect([]int{1}).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected nil, got []int{1}",
				})
			},
		},
		{Scenario: "nil interface",
			Act: func() { var x Expected; Expect(x).IsNil() },
		},
		{Scenario: "non-nil interface",
			Act: func() { var x Expected = &Expecting[int]{0}; Expect(x).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected nil, got &(test.Expecting[int]{value:0})",
				})
			},
		},
		{Scenario: "*string nil",
			Act: func() { var ptr *string; Expect(ptr).IsNil() },
		},
		{Scenario: "string not-nil",
			Act: func() { Expect("non-empty string").IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"test.IsNil: invalid test:",
					"  values of type 'string' are not nilable",
				})
			},
		},
		{Scenario: "*string not-nil",
			Act: func() { ptr := ByRef("non-empty string"); Expect(ptr).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected nil, got &(\"non-empty string\")",
				})
			},
		},
		{Scenario: "*struct not-nil",
			Act: func() { ptr := ByRef(struct{ a int }{a: 1}); Expect(ptr).IsNil() },
			Assert: func(result R) {
				result.Assert([]string{
					"expected nil, got &(struct { a int }{a:1})",
				})
			},
		},
	})
}
