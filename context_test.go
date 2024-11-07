package test

import (
	"context"
	"testing"
)

type key int

var ctxKey key = 1

func ContextKey(ctx context.Context) (string, bool) {
	if v := ctx.Value(ctxKey); v != nil {
		return v.(string), true
	}
	return "", false
}

func TestContextIndicator(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "key is present",
			exec: func(t *testing.T) {
				// ARRANGE
				ctx := context.WithValue(ctx, ctxKey, "value")

				// ACT
				ContextIndicator(t, ctx, ContextKey).Equals(true)
			},
		},
		{scenario: "key not present",
			exec: func(t *testing.T) {
				// ACT
				ContextIndicator(t, ctx, ContextKey).Equals(false)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
func TestContextValue(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "key is present",
			exec: func(t *testing.T) {
				// ARRANGE
				ctx := context.WithValue(ctx, ctxKey, "value")

				// ACT
				ContextValue(t, ctx, ContextKey).Equals("value")
			},
		},
		{scenario: "key not present",
			exec: func(t *testing.T) {
				// ACT
				ContextValue(t, ctx, ContextKey).Equals("")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
