package contexts

import (
	"context"
	"reflect"
	"slices"

	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type ValueMatcher[K comparable, V any] struct {
	Key      K
	Expected V
}

func (vm *ValueMatcher[K, V]) Match(ctx context.Context, opts ...any) bool {
	cv := ctx.Value(vm.Key)

	v, ok := cv.(V)
	if !ok || cv == nil {
		return false
	}

	if cmp, ok := opt.Get[func(V, V) bool](opts); ok {
		return cmp(v, vm.Expected)
	}

	return reflect.DeepEqual(v, vm.Expected)
}

func (km ValueMatcher[K, V]) OnTestFailure(ctx context.Context, opts ...any) []string {
	contextSummary := "context value: " + report.TypeName(km.Key) + "(" + report.Value(km.Key, append(slices.Clone(opts), opt.NoTypeNames)...) + ")"

	got := ctx.Value(km.Key)
	if got == nil {
		return []string{
			contextSummary,
			"  key not present in context",
		}
	}

	switch opt.IsSet(opts, opt.ToNotMatch(true)) {
	case true:
		return []string{
			contextSummary,
			"  key was not expected to have value: " + report.Value(km.Expected, opts...),
		}

	default:
		gotType := report.TypeName(got)
		expType := report.TypeName(km.Expected)
		if gotType != expType {
			return []string{
				contextSummary,
				"  expected value of type: " + expType,
				"  got: " + gotType,
			}
		}

		return []string{
			contextSummary,
			"  expected: " + report.Value(km.Expected, opts...),
			"  got     : " + report.Value(got, opts...),
		}
	}
}
