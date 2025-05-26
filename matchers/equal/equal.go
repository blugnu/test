package equal

import (
	"github.com/blugnu/test/opt"
)

// MARK: EqualMatcher

type Matcher[T comparable] struct {
	Expected T
}

func (m Matcher[T]) Match(got T, opts ...any) bool {
	if equable, ok := any(m.Expected).(interface{ Equal(T) bool }); ok {
		return equable.Equal(got)
	}

	if cmp, ok := opt.Get[func(T, T) bool](opts); ok {
		return cmp(m.Expected, got)
	}

	return m.Expected == got
}

func (m Matcher[T]) OnTestFailure(got T, opts ...any) []string {
	return DeepMatcher[T](m).OnTestFailure(got, opts...)
}
