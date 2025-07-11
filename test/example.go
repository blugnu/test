package test

import (
	"fmt"
	"testing"

	"github.com/blugnu/test/internal/testframe"
)

// Example is used to initialise a testframe for an Example function
// (that is, a runnable example on pkg.go.dev)
func Example() {
	// NOTE: at time of writing, Example functions appear to be run in a single
	// goroutine, so we don't need to worry about concurrency issues here.
	//
	// However, since we always push a new example runner each time this function
	// is called (and this should be called from every example) it does mean
	// that multiple example runners will be pushed onto the testframe stack
	// when running examples.
	//
	// We will live with this for now at least, since:
	//
	// - running examples is a specific and short-lived process
	// - the amount of memory involved is negligible
	// - trying to manage example testframes in a more complex way would
	//   introduce greater complexity and potential for bugs

	testframe.Push(&example{})
}

// example is a dummy implementation of the TestRunner interface.
//
// It is used to provide a TestRunner for Example...() functions in *_test files
// as these do not have a *testing.T available.  The implementation is a no-op for
// most operations except those that produce output or fail the test.
type example struct {
	// failed is used to indicate if the test has failed
	failed bool

	// used to indicate if the test has exited; this causes any calls to
	// Error etc to be ignored.  This is necessary because FailNow() calling
	// runtimeGoexit() would cause a panic if executed in an Example.
	exited bool
}

func (m *example) Name() string {
	return "ExampleTestRunner"
}

func (m *example) Error(args ...any) {
	if m.exited {
		return // NO-OP if the test has exited
	}

	m.failed = true
	fmt.Println(args...)
}

func (m *example) Errorf(format string, args ...any) {
	// NO-OP if the test has exited
	if !m.exited {
		m.Error(fmt.Sprintf(format, args...))
	}
}

func (m *example) Fatal(args ...any) {
	if m.exited {
		return // NO-OP if the test has exited
	}
	m.Error(args...)
	m.FailNow()
}
func (m *example) Fatalf(format string, args ...any) {
	m.Fatal(fmt.Sprintf(format, args...))
}

func (m *example) Run(name string, f func(t *testing.T)) bool {
	if f == nil {
		return true // NO-OP (successful) if no function is provided
	}

	st := &example{}
	testframe.Push(st)
	defer testframe.Pop()

	f(nil) // run the function with a nil *testing.T

	return st.Failed()
}

func (m *example) Fail()        { m.failed = m.failed || !m.exited }
func (m *example) FailNow()     { m.failed = m.failed || !m.exited; m.exited = true }
func (m *example) SkipNow()     { m.exited = true }
func (m *example) Failed() bool { return m.failed }

func (m *example) Cleanup(fn func())     { /* no-op */ }
func (m *example) Helper()               { /* no-op */ }
func (m *example) Parallel()             { /* no-op */ }
func (m *example) Setenv(string, string) { /* no-op */ }
