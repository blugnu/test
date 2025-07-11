package test

import (
	"testing"

	"github.com/blugnu/test/internal/testframe"
)

type Runnable interface {
	Run()
}

func Run(r Runnable) {
	t, ok := testframe.Peek[*testing.T]()
	if !ok {
		panic("ERROR: test.Run() must be called from a Test..(*testing.T) func; it is not supported in Example..() funcs")
	}

	t.Helper()
	r.Run()
}
