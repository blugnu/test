package test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/blugnu/test/opt"
)

type match[T any] struct{ result bool }

func (m match[T]) Match(got T, _ ...any) bool { return m.result }

type matcherWithExpectedField[T comparable] struct {
	match[T]
	Expected T
}

type matcherImplementingExpectedT[T any] struct {
	match[T]
	expected T
}

func (m matcherImplementingExpectedT[T]) Expected() T { return m.expected }

type matcherImplementingExpectedAny[T any] struct {
	match[T]
	expected T
}

func (m matcherImplementingExpectedAny[T]) Expected() any { return m.expected }

type matcherImplementingFormatValue[T any] struct {
	match[T]
	Expected T
	prefix   string
}

func (m matcherImplementingFormatValue[T]) FormatValue(v any, _ ...any) string {
	if m.prefix != "" {
		return fmt.Sprintf("%s %v", m.prefix, v)
	}
	return fmt.Sprintf("%v", v)
}

type matcherImplementingTestFailureOptsReturningString[T any] struct {
	match[T]
	report string
}

func (m matcherImplementingTestFailureOptsReturningString[T]) OnTestFailure(_ ...any) string {
	return m.report
}

// MARK: FailureReports

func TestExpect_TestFailureReporting(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "expectation fails with a nil report",
			Act: func() {
				Expect(any(nil)).err(nil)
			},
			Assert: func(result *R) {
				// TODO: support regex matching for failure report entries to avoid having to test separately for different elements on the same line of a report
				result.Expect(
					"test failed",
				)
			},
		},
		{Scenario: "named expectation fails with a nil report",
			Act: func() {
				Expect(any(nil), "expectation name").err(nil)
			},
			Assert: func(result *R) {
				result.Expect(
					"test failed (expectation name)",
				)
			},
		},
		{Scenario: "expectation fails with an empty string report",
			Act: func() {
				Expect(any(nil)).err("")
			},
			Assert: func(result *R) {
				result.Expect(
					"test failed",
				)
			},
		},
		{Scenario: "expectation fails with an empty string slice report",
			Act: func() {
				Expect(any(nil)).err([]string{})
			},
			Assert: func(result *R) {
				result.Expect(
					"test failed",
				)
			},
		},
		{Scenario: "named expectation fails with an empty string slice report",
			Act: func() {
				Expect(any(nil), "name").err([]string{})
			},
			Assert: func(result *R) {
				result.Expect("test failed (name)")
			},
		},

		{Scenario: "expectation fails with a string report",
			Act: func() {
				Expect(any(nil)).err("failed with message")
			},
			Assert: func(result *R) {
				result.Expect("failed with message")
			},
		},
		{Scenario: "named expectation fails with a string report",
			Act: func() {
				Expect(any(nil), "name").err("failed with message")
			},
			Assert: func(result *R) {
				result.Expect("name: failed with message")
			},
		},
		{Scenario: "expectation fails with a formatted string report",
			Act: func() {
				Expect(any(nil)).errf("failed with %s", "message")
			},
			Assert: func(result *R) {
				result.Expect("failed with message")
			},
		},
		{Scenario: "named expectation fails with a formatted string report",
			Act: func() {
				Expect(any(nil), "named").errf("failed with %s", "message")
			},
			Assert: func(result *R) {
				result.Expect("named: failed with message")
			},
		},

		{Scenario: "expectation fails with a string slice report",
			Act: func() {
				Expect(any(nil)).err([]string{
					"failed with message",
					"and additional information",
				})
			},
			Assert: func(result *R) {
				result.Expect(
					"failed with message",
					"and additional information",
				)
			},
		},
		{Scenario: "named expectation fails with a string slice report",
			Act: func() {
				Expect(any(nil), "name").err([]string{
					"failed with message",
					"and additional information",
				})
			},
			Assert: func(result *R) {
				result.Expect(
					"name:",
					"  failed with message",
					"  and additional information")
			},
		},

		{Scenario: "expectation fails with a non-string report",
			Act: func() {
				Expect(any(nil)).err(42)
			},
			Assert: func(result *R) {
				result.Expect("test failed with: 42")
			},
		},

		{Scenario: "expectation fails with short expected and got string values",
			Act: func() {
				Expect("short 1").To(Equal("short 2"))
			},
			Assert: func(result *R) {
				result.Expect(`expected "short 2", got "short 1"`)
			},
		},
		{Scenario: "expectation fails with long expected and short string values",
			Act: func() {
				Expect("short").To(Equal("long string"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: "long string"`,
					`got     : "short"`,
				)
			},
		},
		{Scenario: "expectation fails with short expected and long string values",
			Act: func() {
				Expect("long string").To(Equal("short"))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: "short"`,
					`got     : "long string"`,
				)
			},
		},
		{Scenario: "expectation fails with short values",
			Act: func() {
				Expect(42).To(Equal(99))
			},
			Assert: func(result *R) {
				result.Expect(`expected 99, got 42`)
			},
		},
		{Scenario: "expectation fails with long values",
			Act: func() {
				Expect(1111111111).To(Equal(1111111112))
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: 1111111112`,
					`got     : 1111111111`,
				)
			},
		},
		{Scenario: "matcher with no Expected field or method",
			Act: func() {
				Expect(12).To(match[int]{false})
			},
			Assert: func(result *R) {
				result.Expect(`got 12`)
			},
		},
		{Scenario: "matcher with Expected field",
			Act: func() {
				Expect(12).To(matcherWithExpectedField[int]{match[int]{false}, 10})
			},
			Assert: func(result *R) {
				result.Expect(`expected 10, got 12`)
			},
		},
		{Scenario: "matcher with Expected() returning T",
			Act: func() {
				Expect(12).To(matcherImplementingExpectedT[int]{match[int]{false}, 10})
			},
			Assert: func(result *R) {
				result.Expect(`expected 10, got 12`)
			},
		},
		{Scenario: "matcher with Expected() returning any",
			Act: func() {
				Expect(12).To(matcherImplementingExpectedAny[int]{match[int]{false}, 10})
			},
			Assert: func(result *R) {
				result.Expect(`expected 10, got 12`)
			},
		},
		{Scenario: "matcher with failure report accepting options only",
			Act: func() {
				Expect(12).To(matcherImplementingTestFailureOptsReturningString[int]{match[int]{}, "custom test failure report"})
			},
			Assert: func(result *R) {
				result.Expect("custom test failure report")
			},
		},
		{Scenario: "matcher implementing FormatValue",
			Act: func() {
				Expect(12).To(matcherImplementingFormatValue[int]{Expected: 10, prefix: "formatted"})
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: formatted 10",
					"got     : formatted 12",
				)
			},
		},
	})
}

func TestExpect_Is(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{
			Scenario: "expecting nil and got nil",
			Act:      func() { var a any; Expect(a).Is(nil) },
		},
		{
			Scenario: "expecting nil when got is of not nilable type",
			Act:      func() { var a any = 1; Expect(a).Is(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.IsNil: values of type 'int' are not nilable",
				)
			},
		},
		{
			Scenario: "expected and got are equal ints",
			Act:      func() { var a any; Expect(a).Is(1) },
			Assert: func(result *R) {
				result.Expect("expected 1, got nil")
			},
		},

		{
			Scenario: "expecting nil and got nil error",
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
			Assert: func(result *R) {
				result.Expect(
					"expected error: sentinel-b",
					"got           : sentinel-a",
				)
			},
		},
		{
			Scenario: "nil error nil vs error",
			Act: func() {
				var err error = errors.New("error")
				Expect(err).Is(nil)
			},
			Assert: func(result *R) {
				result.Expect("expected nil, got error")
			},
		},
		{
			Scenario: "struct is equal struct",
			Act:      func() { Expect(struct{ a int }{a: 1}).Is(struct{ a int }{a: 1}) },
		},
		{
			Scenario: "struct is inequal struct",
			Act:      func() { Expect(struct{ a int }{a: 1}).Is(struct{ a int }{a: 2}) },
			Assert: func(result *R) {
				result.Expect(
					"expected: struct { a int }{a:2}",
					"got     : struct { a int }{a:1}",
				)
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
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.IsNil: values of type 'int' are not nilable",
				)
			},
		},
		{Scenario: "struct is nil",
			Act: func() { Expect(struct{}{}).IsNil() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.IsNil: values of type 'struct {}' are not nilable",
				)
			},
		},
		{Scenario: "error is nil",
			Act: func() { Expect(errors.New("some error message")).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got error: some error message")
			},
		},
		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).IsNil() },
		},
		{Scenario: "non-nil slice",
			Act: func() { Expect([]int{1}).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got []int{1}")
			},
		},
		{Scenario: "nil interface",
			Act: func() { var x any; Expect(x).IsNil() },
		},
		{Scenario: "non-nil interface",
			Act: func() { var x any = byref("any"); Expect(x).IsNil() },
			Assert: func(result *R) {
				result.Expect(`expected nil, got &("any")`)
			},
		},
		{Scenario: "*string nil",
			Act: func() { var ptr *string; Expect(ptr).IsNil() },
		},
		{Scenario: "*string non-nil",
			Act: func() { Expect(byref("string")).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got &(\"string\")")
			},
		},
		{Scenario: "*string non-nil with unquoted strings",
			Act: func() { Expect(byref("string")).IsNil(opt.UnquotedStrings()) },
			Assert: func(result *R) {
				result.Expect("expected nil, got &(string)")
			},
		},
		{Scenario: "string not-nil",
			Act: func() { Expect("non-empty string").IsNil() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.IsNil: values of type 'string' are not nilable",
				)
			},
		},
		{Scenario: "*string not-nil",
			Act: func() { ptr := byref("non-empty string"); Expect(ptr).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got &(\"non-empty string\")")
			},
		},
		{Scenario: "*struct not-nil",
			Act: func() { ptr := byref(struct{ a int }{a: 1}); Expect(ptr).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got &(struct { a int }{a:1})")
			},
		},

		{Scenario: "with custom failure report",
			Act: func() {
				Expect(byref(42)).IsNil(opt.FailureReport(func(a ...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect("custom failure report")
			},
		},
	})
}

func TestExpect_IsNotNil(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "got non-nil nillable type",
			Act: func() { Expect(errors.New("error")).IsNotNil() },
		},
		{Scenario: "got non-nillable type",
			Act: func() { Expect(0).IsNotNil() },
		},
		{Scenario: "got nil",
			Act: func() { Expect(any(nil)).IsNotNil() },
			Assert: func(result *R) {
				result.Expect("expected not nil, got nil")
			},
		},
		{Scenario: "with custom failure report",
			Act: func() {
				Expect(any(nil)).IsNotNil(opt.FailureReport(func(a ...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect("custom failure report")
			},
		},
	})
}

func TestExpect_To(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "no matcher",
			Act: func() { Expect(42).To(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.To: a matcher must be specified",
				)
			},
		},
	})
}

func TestExpect_ToNot(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "a heterogenous matcher does not match and was not expected to",
			Act: func() { Expect(context.Background()).ToNot(HaveContextKey(99)) },
		},
		{Scenario: "a heterogenous matcher matches and was not expected to",
			Act: func() {
				type k int
				ctx := context.WithValue(context.Background(), k(99), "value")
				Expect(ctx).ToNot(HaveContextKey(k(99)))
			},
			Assert: func(result *R) {
				result.Expect(TestFailed, opt.IgnoreReport(true)) // testing the behaviour of ToNot(), not the output of the matcher used
			},
		},
	})
}
