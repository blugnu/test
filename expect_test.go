package test //nolint: testpackage // tests private types and functions

import (
	"errors"
	"fmt"
	"testing"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type match[T any] struct{ result bool }

func (m match[T]) Match(T, ...any) bool { return m.result }

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

type matcherImplementingTestFailureAnyOptsReturningString[T any] struct {
	match[T]
	report string
}

func (m matcherImplementingTestFailureAnyOptsReturningString[T]) OnTestFailure(any, ...any) string {
	return m.report
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

	Run(HelperTests([]HelperScenario{
		{Scenario: "expectation fails with a nil report",
			Act: func() {
				Expect(any(nil)).err(nil)
			},
			Assert: func(result *R) {
				result.Expect("test failed")
			},
		},
		{Scenario: "named expectation fails with a nil report",
			Act: func() {
				Expect(any(nil), "expectation name").err(nil)
			},
			Assert: func(result *R) {
				result.Expect(
					"expectation name: test failed",
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
				result.Expect("name: test failed")
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
				result.Expect("test failed: 42")
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
		{Scenario: "matcher (pointer) with Expected field",
			Act: func() {
				Expect(12).To(&matcherWithExpectedField[int]{match[int]{false}, 10})
			},
			Assert: func(result *R) {
				result.Expect(`expected 10, got 12`)
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
		{Scenario: "matcher with failure report accepting any",
			Act: func() {
				Expect(12).To(matcherImplementingTestFailureAnyOptsReturningString[int]{match[int]{}, "custom test failure report"})
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
	}...))
}

func TestExpect_Is(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{
			Scenario: "expecting nil and got nil",
			Act:      func() { var a any; Expect(a).Is(nil) },
		},
		{
			Scenario: "expecting nil when got is of not nilable type",
			Act:      func() { var a any = 1; Expect(a).Is(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'int' are not nilable",
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
				err := errors.New("error")
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
	}...))
}

func TestExpect_Should(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no matcher",
			Act: func() { Expect(42).Should(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid("test.Should: a matcher must be specified")
			},
		},
		{Scenario: "matcher passes",
			Act: func() { Expect([]int{}).Should(BeEmpty()) },
		},
		{Scenario: "matcher fails",
			Act: func() {
				Expect([]int{1}).Should(BeEmpty())
			},
			Assert: func(result *R) {
				result.Expect(TestFailed, opt.IgnoreReport(true)) // testing the behaviour of Should(); matcher report is not significant
			},
		},
	}...))
}

func TestExpect_ShouldNot(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no matcher",
			Act: func() { Expect(true).ShouldNot(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid("test.ShouldNot: a matcher must be specified")
			},
		},
		{Scenario: "matcher fails",
			Act: func() { Expect([1]int{}).ShouldNot(BeEmpty()) },
		},
		{Scenario: "matcher passes",
			Act: func() {
				Expect([]int{}).ShouldNot(BeEmpty())
			},
			Assert: func(result *R) {
				result.Expect(TestFailed, opt.IgnoreReport(true)) // testing the behaviour of ToNot(); matcher report is not significant
			},
		},
	}...))
}

func TestExpect_To(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "no matcher",
			Act: func() { Expect(42).To(nil) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"test.To: a matcher must be specified",
				)
			},
		},
	}...))
}

func TestExpect_ToNot(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "matcher fails, as expected",
			Act: func() { Expect(true).ToNot(Equal(false)) },
		},
		{Scenario: "matcher does not fail but should have",
			Act: func() {
				Expect(true).ToNot(Equal(true))
			},
			Assert: func(result *R) {
				result.Expect(TestFailed, opt.IgnoreReport(true)) // testing the behaviour of ToNot(), not the output of the matcher used
			},
		},
	}...))
}

func TestRequire(t *testing.T) {
	With(t)

	result := TestHelper(func() {
		Require(true).To(Equal(false))
		Expect("black").To(Equal("white"))
	})

	result.Expect("expected false, got true")
}

func ExampleRequire() {
	test.Example()

	// this test will fail
	Require(true).To(Equal(false))

	// this will not be executed because the previous expectation was
	// required to pass and did not
	Expect("apples").To(Equal("oranges"))

	// Output:
	// expected false, got true
}
