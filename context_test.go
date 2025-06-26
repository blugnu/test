package test_test

import (
	"context"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/opt"
)

func TestContext(t *testing.T) {
	With(t)

	type key string
	ctx := context.WithValue(context.Background(), key("key"), "value")

	Run("ContextKey", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "expected key present",
				Act: func() {
					Expect(ctx).To(HaveContextKey(key("key")))
				},
			},
			{Scenario: "expected key not present",
				Act: func() {
					Expect(ctx).To(HaveContextKey(key("other-key")))
				},
				Assert: func(result *R) {
					result.Expect(TestFailed, opt.IgnoreReport(true))
				},
			},
		})
	})

	Run("ContextValue", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "expected value present",
				Act: func() {
					Expect(ctx).To(HaveContextValue(key("key"), "value"))
				},
			},
			{Scenario: "expected value not present",
				Act: func() {
					Expect(ctx).To(HaveContextValue(key("other-key"), "value"))
				},
				Assert: func(result *R) {
					result.Expect(TestFailed, opt.IgnoreReport(true))
				},
			},
			{Scenario: "expected value present but different",
				Act: func() {
					Expect(ctx).To(HaveContextValue(key("key"), "other value"))
				},
				Assert: func(result *R) {
					result.Expect(TestFailed, opt.IgnoreReport(true))
				},
			},
		})
	})
}
