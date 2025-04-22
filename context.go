package test

import (
	"context"
	"fmt"
)

type ContextKeyMatcher[K comparable] struct {
	ctx  context.Context
	want K
}

func (km ContextKeyMatcher[K]) Expected() any {
	return km.want
}

func (km ContextKeyMatcher[K]) Format(expected any, got any) []string {
	return []string{
		fmt.Sprintf("expected key: %[1]T(%[1]v): not present in context", expected),
	}
}

// Match ensures the matcher is compatible with Expect(context.Context)
func (km ContextKeyMatcher[K]) Match(ctx context.Context, _ ...any) bool { return true }

func (km ContextKeyMatcher[K]) MatchAny(got any, _ ...any) bool {
	km.ctx = got.(context.Context)
	v := km.ctx.Value(km.want)
	return v != nil
}

func HaveContextKey[K comparable](k K) ContextKeyMatcher[K] {
	return ContextKeyMatcher[K]{want: k}
}

type ContextValueMatcher[K comparable, V any] struct {
	ctx  context.Context
	key  K
	want V
}

func (km ContextValueMatcher[K, V]) Format() []string {
	v := km.ctx.Value(km.key)
	if v == nil {
		return []string{
			"context value:",
			fmt.Sprintf("  %[1]T(%[1]v): not present in context", km.key),
		}
	}

	return []string{
		fmt.Sprintf("context value: %[1]T(%[1]v)", km.key),
		fmt.Sprintf("  expected: %v", km.want),
		fmt.Sprintf("  got     : %v", v),
	}
}

func (vm *ContextValueMatcher[K, V]) Match(ctx context.Context, _ ...any) bool {
	vm.ctx = ctx
	return true
}

func (vm ContextValueMatcher[K, V]) MatchAny(got any, _ ...any) bool {
	v := vm.ctx.Value(vm.key)
	if v == nil {
		return false
	}
	return any(v) == any(vm.want)
}

func HaveContextValue[K comparable, V any](k K, v V) *ContextValueMatcher[K, V] {
	return &ContextValueMatcher[K, V]{key: k, want: v}
}
