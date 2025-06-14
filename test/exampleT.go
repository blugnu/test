package test

import (
	"fmt"
	"testing"
)

// ExampleT() provides a dummy implementation of the TestRunner interface
// for use in Example...() functions in *_test files where a *testing.T is not
// available.
//
//	func ExampleFunctionName() {
//	    With(test.ExampleT())
//
//	    // ... example code ...
//
//	    // Output:
//	    // the expected output
//	}
func ExampleT() *exampleT { return &exampleT{} }

// ExampleT() is a dummy implementation of the TestRunner interface.
//
// It is used to provide a TestRunner for Example...() functions in *_test files
// as these do not have a *testing.T available.  The implementation is a no-op for
// most operations except those that produce output or fail the test.
type exampleT struct {
	// failed is used to indicate if the test has failed
	failed bool

	// used to indicate if the test has exited; this causes any calls to
	// Error etc to be ignored.  This is necessary because FailNow() calling
	// runtimeGoexit() would cause a panic if executed in an Example.
	exited bool
}

func (m exampleT) Name() string {
	return "ExampleTestRunner"
}

func (m *exampleT) Error(args ...any) {
	if m.exited {
		return // NO-OP if the test has exited
	}

	m.failed = true
	fmt.Println(args...)
}

func (m *exampleT) Errorf(format string, args ...any) {
	// NO-OP if the test has exited
	if !m.exited {
		m.Error(fmt.Sprintf(format, args...))
	}
}

func (m *exampleT) Fatal(args ...any) {
	if m.exited {
		return // NO-OP if the test has exited
	}
	m.Error(args...)
	m.FailNow()
}
func (m *exampleT) Fatalf(format string, args ...any) {
	m.Fatal(fmt.Sprintf(format, args...))
}

func (m exampleT) Run(name string, f func(t *testing.T)) bool {
	return false
}

func (m *exampleT) Fail()        { m.failed = m.failed || !m.exited }
func (m *exampleT) FailNow()     { m.failed = m.failed || !m.exited; m.exited = true }
func (m *exampleT) SkipNow()     { m.exited = true }
func (m *exampleT) Failed() bool { return m.failed }

func (m exampleT) Cleanup(func())        { /* no-op */ }
func (m exampleT) Helper()               { /* no-op */ }
func (m exampleT) Parallel()             { /* no-op */ }
func (m exampleT) Setenv(string, string) { /* no-op */ }
