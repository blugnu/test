package test

import (
	"errors"
	"testing"

	"github.com/blugnu/test/opt"
)

func TestExpect_IsNil(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "nil",
			Act: func() { var subject any; Expect(subject).IsNil() },
		},
		{Scenario: "nil error",
			Act: func() { var err error; Expect(err).IsNil() },
		},

		{Scenario: "non-nil error",
			Act: func() { Expect(errors.New("error")).IsNil() },
			Assert: func(result *R) {
				result.Expect("expected nil, got error: error")
			},
		},
		{Scenario: "non-nilable subject",
			Act: func() { Expect(0).IsNil() },
			Assert: func(result *R) {
				result.ExpectInvalid(
					"nilness.Matcher: values of type 'int' are not nilable",
				)
			},
		},
		{Scenario: "with custom failure report",
			Act: func() {
				Expect(errors.New("not nil")).IsNil(opt.FailureReport(func(a ...any) []string {
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
		{Scenario: "non-nil error",
			Act: func() { Expect(errors.New("error")).IsNotNil() },
		},
		{Scenario: "non-nilable subject",
			Act: func() { Expect(0).IsNotNil() },
		},

		{Scenario: "nil",
			Act: func() { var subject any; Expect(subject).IsNotNil() },
			Assert: func(result *R) {
				result.Expect("expected not nil, got nil")
			},
		},
		{Scenario: "with custom failure report",
			Act: func() {
				var subject any
				Expect(subject).IsNotNil(opt.FailureReport(func(a ...any) []string {
					return []string{"custom failure report"}
				}))
			},
			Assert: func(result *R) {
				result.Expect("custom failure report")
			},
		},
	})
}
