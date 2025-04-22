package test

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	With(t)

	// ARRANGE
	type key int
	ctxBg := context.Background()
	ctxWithKey := context.WithValue(ctxBg, key(1), "value-1")

	type testCase struct {
		scenario string
		act      func()
		assert   func(R)
	}
	RunScenarios(
		func(tc testCase) {
			result := Test(tc.act)
			tc.assert(result)
		},
		[]testCase{
			{
				scenario: "key/present",
				act:      func() { Expect(ctxWithKey).To(HaveContextKey(key(1))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{
				scenario: "key/not present",
				act:      func() { Expect(ctxBg).To(HaveContextKey(key(1))) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						currentFilename(),
						"expected key: test.key(1): not present in context",
					}))
				},
			},
			{
				scenario: "value/present",
				act:      func() { Expect(ctxWithKey).To(HaveContextValue(key(1), "value-1")) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestPassed))
					Expect(result.Report).IsEmpty()
				},
			},
			{
				scenario: "value/key not present",
				act:      func() { Expect(ctxBg).To(HaveContextValue(key(2), "value-1")) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						currentFilename(),
						"context value:",
						"  test.key(2): not present in context",
					}))
				},
			},
			{
				scenario: "value/not equal",
				act:      func() { Expect(ctxWithKey).To(HaveContextValue(key(1), "value-2")) },
				assert: func(result R) {
					Expect(result.Outcome).To(Equal(TestFailed))
					Expect(result.Report).To(ContainStrings([]string{
						currentFilename(),
						"context value: test.key(1)",
						"  expected: value-2",
						"  got     : value-1",
					}))
				},
			},
		},
	)
}

// func TestContextIndicator(t *testing.T) {
// 	// ARRANGE
// 	ctx := context.Background()

// 	testcases := []struct {
// 		scenario string
// 		exec     func(t *testing.T)
// 	}{
// 		{scenario: "key is present",
// 			exec: func(t *testing.T) {
// 				// ARRANGE
// 				ctx := context.WithValue(ctx, ctxKey, "value")

// 				// ACT
// 				ContextIndicator(t, ctx, ContextKey).Equals(true)
// 			},
// 		},
// 		{scenario: "key not present",
// 			exec: func(t *testing.T) {
// 				// ACT
// 				ContextIndicator(t, ctx, ContextKey).Equals(false)
// 			},
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.exec(t)
// 		})
// 	}
// }
// func TestContextValue(t *testing.T) {
// 	// ARRANGE
// 	ctx := context.Background()

// 	testcases := []struct {
// 		scenario string
// 		exec     func(t *testing.T)
// 	}{
// 		{scenario: "key is present",
// 			exec: func(t *testing.T) {
// 				// ARRANGE
// 				ctx := context.WithValue(ctx, ctxKey, "value")

// 				// ACT
// 				ContextValue(t, ctx, ContextKey).Equals("value")
// 			},
// 		},
// 		{scenario: "key not present",
// 			exec: func(t *testing.T) {
// 				// ACT
// 				ContextValue(t, ctx, ContextKey).Equals("")
// 			},
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.exec(t)
// 		})
// 	}
// }
