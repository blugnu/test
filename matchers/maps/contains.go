package maps

import (
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type ContainsMatcher[K comparable, V any] struct {
	Expected map[K]V
}

func (m ContainsMatcher[K, V]) Match(got map[K]V, opts ...any) bool {
	if len(m.Expected) > len(got) {
		return false
	}

	return containsMap(got, m.Expected, opts...)
}

func (m ContainsMatcher[K, V]) OnTestFailure(got map[K]V, opts ...any) []string {
	result := make([]string, 0, 2+len(got)+len(m.Expected))

	inv := opt.IsSet(opts, opt.ToNotMatch(true))
	switch {
	case inv:
		result = report.AppendMap(result, m.Expected, opt.WithName(opts, "expected map not containing:")...)

	default:
		result = report.AppendMap(result, m.Expected, opt.WithName(opts, "expected map containing:")...)
		result = report.AppendMap(result, got, opt.Force(opts, opt.Name("got:"))...)
	}
	return result
}
