package ordered

import (
	"cmp"

	"github.com/blugnu/test/opt"
)

type IsLessThan[T cmp.Ordered] struct {
	Expected T
}

func (m IsLessThan[T]) Match(got T, opts ...any) bool {
	if cmp, ok := opt.Get[func(T, T) bool](opts); ok {
		return cmp(m.Expected, got)
	}
	return got < m.Expected
}

func (m IsLessThan[T]) OnTestFailure(got T, opts ...any) []string {
	return failureReport("less than", m.Expected, got, opts...)
}
