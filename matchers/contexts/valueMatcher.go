package contexts

import (
	"context"
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

type ValueMatcher[K comparable, V any] struct {
	Key      K
	Expected V
}

func (vm *ValueMatcher[K, V]) Match(ctx context.Context, opts ...any) bool {
	v := ctx.Value(vm.Key)
	if v == nil {
		return false
	}

	if _, ok := v.(V); !ok {
		return false
	}

	if cmp, ok := opt.Get[func(V, V) bool](opts); ok {
		return cmp(v.(V), vm.Expected)
	}

	return reflect.DeepEqual(v, vm.Expected)
}

func (km ValueMatcher[K, V]) OnTestFailure(ctx context.Context, opts ...any) []string {
	got := ctx.Value(km.Key)
	if got == nil {
		return []string{
			fmt.Sprintf("context value: %[1]T(%[1]v):", km.Key),
			"  key not present in context",
		}
	}

	result := []string{
		fmt.Sprintf("context value: %[1]T(%[1]v)", km.Key),
	}
	switch opt.IsSet(opts, opt.ToNotMatch(true)) {
	case true:
		result = append(result, fmt.Sprintf("  key was not expected to have value: %v", opt.ValueAsString(km.Expected, opts...)))
	default:
		gotType := fmt.Sprintf("%T", got)
		expType := fmt.Sprintf("%T", km.Expected)
		if gotType != expType {
			result = append(result, "  expected value of type: "+expType)
			result = append(result, "  got: "+gotType)
			return result
		}

		result = append(result,
			fmt.Sprintf("  expected: %v", opt.ValueAsString(km.Expected, opts...)),
			fmt.Sprintf("  got     : %v", opt.ValueAsString(got, opts...)),
		)
	}
	return result
}
