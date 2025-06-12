package test

import (
	"context"
	"testing"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
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

func ExampleHaveContextKey() {
	test.Example()

	type key int
	ctx := context.WithValue(context.Background(), key(57), "varieties")

	// these tests will pass
	Expect(ctx).To(HaveContextKey(key(57)))
	Expect(ctx).ToNot(HaveContextKey(key(58)))

	// this test will fail
	Expect(ctx).To(HaveContextKey(key(58)))

	// Output:
	// expected key: test.key(58)
	//   key not present in context
}

func ExampleHaveContextValue() {
	// this is needed to make the example work; this would be usually
	// be `With(t)` where `t` is the *testing.T
	test.Example()

	type key int
	ctx := context.WithValue(context.Background(), key(57), "varieties")

	// these tests will pass
	Expect(ctx).To(HaveContextValue(key(57), "varieties"))
	Expect(ctx).ToNot(HaveContextValue(key(56), "varieties"))
	Expect(ctx).ToNot(HaveContextValue(key(57), "flavours"))

	// this test will fail
	Expect(ctx).To(HaveContextValue(key(57), "flavours"))

	// Output:
	// context value: test.key(57)
	//   expected: "flavours"
	//   got     : "varieties"
}
