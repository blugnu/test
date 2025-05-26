package maps

import (
	"github.com/blugnu/test/opt"
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
		result = appendToReport(result, "expected: map not containing:", m.Expected, opts...)

	default:
		result = appendToReport(result, "expected: map containing:", m.Expected, opts...)
		result = appendToReport(result, "got:", got, opts...)
	}
	return result
}
