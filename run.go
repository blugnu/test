package test

import (
	"testing"

	"github.com/blugnu/test/internal/testframe"
)

// Runnable is an interface implemented by types that can be executed
// in a test frame using the [Run] function provided by the interface.
type Runnable interface {
	Run()
}

// Run executes the provided [Runnable] in the current test frame.
//
// It must be called from a Test..(*testing.T) func; it is not
// supported in Example..() funcs.
func Run(r Runnable) {
	t, ok := testframe.Peek[*testing.T]()
	if !ok {
		panic("ERROR: test.Run() must be called from a Test..(*testing.T) func; it is not supported in Example..() funcs")
	}

	t.Helper()
	r.Run()
}
