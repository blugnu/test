package maps

import (
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

type EqualMatcher[K comparable, V any] struct {
	Expected map[K]V
}

func (m EqualMatcher[K, V]) Match(got map[K]V, opts ...any) bool {
	if len(m.Expected) != len(got) {
		return false
	}

	return containsMap(got, m.Expected, opts...)
}

func (m EqualMatcher[K, V]) OnTestFailure(got map[K]V, opts ...any) []string {
	result := make([]string, 0, 2+len(got)+len(m.Expected))

	inv := opt.IsSet(opts, opt.ToNotMatch(true))
	switch {
	case len(m.Expected) == 0 && inv:
		result = append(result, "unexpected: <empty map>")

	case len(m.Expected) == 0:
		result = report.AppendMap(result, m.Expected, opt.WithName(opts, "expected:")...)
		result = report.AppendMap(result, got, opt.Force(opts, opt.Name("got:"))...)

	case inv:
		result = report.AppendMap(result, m.Expected, opt.WithNamef(opts, "expected %T not equal to:", m.Expected)...)

	default:
		result = report.AppendMap(result, m.Expected, opt.WithNamef(opts, "expected %T:", m.Expected)...)
		result = report.AppendMap(result, got, opt.Force(opts, opt.Name("got:"))...)
	}
	return result
}
