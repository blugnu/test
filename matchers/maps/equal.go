package maps

import (
	"github.com/blugnu/test/opt"
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
		result = AppendToReport(result, "expected:", m.Expected, opts...)
		result = AppendToReport(result, "got:", got, opts...)

	case inv:
		result = AppendToReport(result, "expected: map not equal to:", m.Expected, opts...)

	default:
		result = AppendToReport(result, "expected map:", m.Expected, opts...)
		result = AppendToReport(result, "got:", got, opts...)
	}
	return result
}
