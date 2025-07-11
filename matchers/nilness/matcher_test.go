package nilness_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func byref[T any](v T) *T {
	return &v
}

func TestShould_BeNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "nil is nil",
			Act: func() { var nilAny any; Expect(nilAny).Should(BeNil()) },
		},
		{Scenario: "int is nil",
			Act: func() { var intAny any = 0; Expect(intAny).Should(BeNil()) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'int' are not nilable",
				)
			},
		},
		{Scenario: "*struct nil",
			Act: func() { var nilStructPtr *struct{ a int }; Expect(nilStructPtr).Should(BeNil()) },
		},
		{Scenario: "struct is nil",
			Act: func() { Expect(struct{}{}).Should(BeNil()) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'struct {}' are not nilable",
				)
			},
		},
		{Scenario: "error is nil",
			Act: func() { Expect(errors.New("some error message")).Should(BeNil()) },
			Assert: func(result *R) {
				result.Expect("expected nil, got error: some error message")
			},
		},
		{Scenario: "nil slice",
			Act: func() { Expect([]int(nil)).Should(BeNil()) },
		},
		{Scenario: "non-nil slice",
			Act: func() { Expect([]int{1}).Should(BeNil()) },
			Assert: func(result *R) {
				result.Expect("expected nil, got []int{1}")
			},
		},
		{Scenario: "nil interface",
			Act: func() { var x any; Expect(x).Should(BeNil()) },
		},
		{Scenario: "non-nil interface",
			Act: func() { var x any = byref("any"); Expect(x).Should(BeNil()) },
			Assert: func(result *R) {
				result.Expect(`expected nil, got &("any")`)
			},
		},
		{Scenario: "*string nil",
			Act: func() { var ptr *string; Expect(ptr).Should(BeNil()) },
		},
		{Scenario: "*string non-nil",
			Act: func() { Expect(byref("string")).Should(BeNil()) },
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
			Act: func() { Expect("non-empty string").Should(BeNil()) },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'string' are not nilable",
				)
			},
		},

		{Scenario: "*string not-nil",
			Act: func() { ptr := byref("non-empty string"); Expect(ptr).Should(BeNil()) },
			Assert: func(result *R) {
				result.Expect("expected nil, got &(\"non-empty string\")")
			},
		},
		{Scenario: "*struct not-nil",
			Act: func() { ptr := byref(struct{ a int }{a: 1}); Expect(ptr).Should(BeNil()) },
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
	}...))
}

func TestShouldNot_BeNil(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "nil is not nil",
			Act: func() { Expect((any)(nil)).ShouldNot(BeNil()) },
			Assert: func(result *R) {
				result.Expect("expected not nil, got nil")
			},
		},
		{Scenario: "int is not nil",
			Act: func() { Expect(0).ShouldNot(BeNil()) },
		},
		{Scenario: "struct is not nil",
			Act: func() { Expect(struct{}{}).ShouldNot(BeNil()) },
		},
		{Scenario: "error is not nil",
			Act: func() { Expect(errors.New("some error message")).ShouldNot(BeNil()) },
		},
	}...))
}
