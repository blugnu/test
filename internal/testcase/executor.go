package testcase

import (
	"fmt"

	"github.com/blugnu/test/test"
)

type Executor[T any] struct {
	anon  AnonCaseExecution[T]
	named NamedCaseExecution[T]
}

func NewExecutor[T any](fn any) Executor[T] {
	switch fn := fn.(type) {
	case func(string, T):
		return Executor[T]{named: fn}

	case func(T):
		return Executor[T]{anon: fn}

	default:
		GetT().Helper()
		test.Warning(fmt.Sprintf("%T is not a valid test function", fn))
		return Executor[T]{}
	}
}

func (exec Executor[T]) Execute(name string, tc T) {
	GetT().Helper()

	switch {
	case exec.anon != nil:
		exec.anon(tc)

	case exec.named != nil:
		exec.named(name, tc)

	default:
		test.Invalid("a test function must be provided")
	}
}
