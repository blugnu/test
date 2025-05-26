package contexts_test

import (
	"context"
	"testing"

	. "github.com/blugnu/test"
)

func TestContext(t *testing.T) {
	With(t)

	// ARRANGE
	type key int
	ctxBg := context.Background()
	ctxWithKey := context.WithValue(ctxBg, key(1), "value-1")

	Run("KeyMatcher", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "expect key to be present that is present",
				Act: func() { Expect(ctxWithKey).To(HaveContextKey(key(1))) },
			},
			{Scenario: "expect key to be present that is not present",
				Act: func() { Expect(ctxBg).To(HaveContextKey(key(1))) },
				Assert: func(result *R) {
					result.Expect(
						"expected key: contexts_test.key(1)",
						"  key not present in context",
					)
				},
			},
			{Scenario: "expect key to not be present that is present",
				Act: func() { Expect(ctxWithKey).ToNot(HaveContextKey(key(1))) },
				Assert: func(result *R) {
					result.Expect(
						"unexpected key: contexts_test.key(1)",
						"  key should not be present in context",
					)
				},
			},
		})
	})

	Run("ValueMatcher", func() {
		RunTestScenarios([]TestScenario{
			{Scenario: "expect key to have value that is present and has that value",
				Act: func() {
					Expect(ctxWithKey).To(HaveContextValue(key(1), "value-1"))
				},
			},
			{Scenario: "expect key to have value that is present and has a different value",
				Act: func() { Expect(ctxWithKey).To(HaveContextValue(key(1), "value-2")) },
				Assert: func(result *R) {
					result.Expect(
						"context value: contexts_test.key(1)",
						`  expected: "value-2"`,
						`  got     : "value-1"`,
					)
				},
			},
			{Scenario: "expect key to have value when the key is not present",
				Act: func() { Expect(ctxBg).To(HaveContextValue(key(2), "value-1")) },
				Assert: func(result *R) {
					result.Expect(
						"context value: contexts_test.key(2)",
						"  key not present in context",
					)
				},
			},
			{Scenario: "expect key to not have value that is present and has that value",
				Act: func() { Expect(ctxWithKey).ToNot(HaveContextValue(key(1), "value-1")) },
				Assert: func(result *R) {
					result.Expect(
						"context value: contexts_test.key(1)",
						`  key was not expected to have value: "value-1"`,
					)
				},
			},
			{Scenario: "expected value exists but is of different type than expected",
				Act: func() {
					type custom string
					Expect(ctxWithKey).To(HaveContextValue(key(1), custom("value-1")))
				},
				Assert: func(result *R) {
					result.Expect(
						"context value: contexts_test.key(1)",
						`  expected value of type: contexts_test.custom`,
						`  got: string`,
					)
				},
			},
			{Scenario: "custom value comparison function (type V)",
				Act: func() {
					Expect(ctxWithKey).To(HaveContextValue(key(1), "any"), func(a, b string) bool {
						return true
					})
				},
			},
		})
	})
}
