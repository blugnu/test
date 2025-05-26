package test

import (
	"fmt"
	"runtime"
	"testing"
)

type TestingT interface {
	Cleanup(func())
	Name() string
	Run(string, func(t *testing.T)) bool
	Error(...any)
	Errorf(string, ...any)
	Fail()
	FailNow()
	Fatal(...any)
	Fatalf(string, ...any)
	Helper()
	Parallel()
	Setenv(string, string)
	SkipNow()
}

var _ TestingT = (*ExampleTestRunner)(nil)

// MARK: ExampleTestRunner

// ExampleTestRunner is a mock implementation of the TestRunner interface.
//
// It is used to provide a TestRunner for Example...() functions since these
// do not have a *testing.T available.  The implementation is a no-op for
// most operations except those that produce output or fail the test.
type ExampleTestRunner struct{}

func (m ExampleTestRunner) Name() string {
	return "ExampleTestRunner"
}

func (m ExampleTestRunner) Run(name string, f func(t *testing.T)) bool {
	return false
}
func (m ExampleTestRunner) Error(args ...any) {
	fmt.Println(args...)
}
func (m ExampleTestRunner) Errorf(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println()
}
func (m ExampleTestRunner) Fail()    { /* no-op */ }
func (m ExampleTestRunner) FailNow() { runtime.Goexit() }
func (m ExampleTestRunner) Fatal(args ...any) {
	m.Error(args...)
	m.FailNow()
}
func (m ExampleTestRunner) Fatalf(format string, args ...any) {
	m.Errorf(format, args...)
	m.FailNow()
}
func (m ExampleTestRunner) Cleanup(func())        { /* no-op */ }
func (m ExampleTestRunner) Helper()               { /* no-op */ }
func (m ExampleTestRunner) Parallel()             { /* no-op */ }
func (m ExampleTestRunner) Setenv(string, string) { /* no-op */ }
func (m ExampleTestRunner) SkipNow()              { /* no-op */ }

// MARK: T()

// GetT retrieves the *testing.T for the calling test frame, by calling T().
//
// This is used where calling the T() function directly is not possible
// e.g. due to a name collision with a generic type parameter.
func GetT() TestingT {
	t := TestFrame(3)
	if t == nil {
		panic(fmt.Errorf("GetT: %w", ErrNoTestFrame))
	}
	return t
}

// hasT retrieves the *testing.T for the calling test frame, by calling T().
//
// This is used where calling the T() function directly is not possible
// e.g. due to a name collision with a generic type parameter.
func hasT() TestingT {
	return TestFrame(3)
}

// T retrieves the TestRunner for the calling test frame. When running in a test
// frame, this will return the *testing.T for the test.  When running in an
// example, this will return an ExampleTestRunner.
func T() TestingT {
	t := TestFrame(3)
	if t == nil {
		panic(fmt.Errorf("GetT: %w", ErrNoTestFrame))
	}
	return t
}
