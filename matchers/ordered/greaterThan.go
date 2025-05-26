package ordered

import (
	"cmp"

	"github.com/blugnu/test/opt"
)

type IsGreaterThan[T cmp.Ordered] struct {
	Expected T
}

func (m IsGreaterThan[T]) Match(got T, opts ...any) bool {
	if cmp, ok := opt.Get[func(T, T) bool](opts); ok {
		return cmp(got, m.Expected)
	}
	return got > m.Expected
}

func (m IsGreaterThan[T]) OnTestFailure(got T, opts ...any) []string {
	return failureReport("greater than", m.Expected, got, opts...)
}
