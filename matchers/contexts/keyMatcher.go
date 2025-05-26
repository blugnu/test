package contexts

import (
	"context"
	"fmt"

	"github.com/blugnu/test/opt"
)

type KeyMatcher[K comparable] struct {
	Expected K
	// got      context.Context
}

// Match ensures the matcher is compatible with Expect(context.Context)
func (km *KeyMatcher[K]) Match(ctx context.Context, opts ...any) bool {
	return ctx.Value(km.Expected) != nil
}

func (km KeyMatcher[K]) OnTestFailure(opts ...any) []string {
	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{
			fmt.Sprintf("unexpected key: %[1]T(%[1]v)", km.Expected),
			"  key should not be present in context",
		}
	}
	return []string{
		fmt.Sprintf("expected key: %[1]T(%[1]v)", km.Expected),
		"  key not present in context",
	}
}
