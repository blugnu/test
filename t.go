package test

import (
	"fmt"
	"runtime"
	"testing"
)

type TestRunner interface {
	Run(string, func(t *testing.T)) bool
	Error(...any)
	Errorf(string, ...any)
	Fail()
	FailNow()
	Fatal(...any)
	Fatalf(string, ...any)
	Helper()
}

var _ TestRunner = (*ExampleTestRunner)(nil)
var _ TestRunner = (*MockTestRunner)(nil)

// MARK: ExampleTestRunner

type ExampleTestRunner struct{}

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
func (m ExampleTestRunner) Fail()    {}
func (m ExampleTestRunner) FailNow() { runtime.Goexit() }
func (m ExampleTestRunner) Fatal(args ...any) {
	m.Error(args...)
	m.FailNow()
}
func (m ExampleTestRunner) Fatalf(format string, args ...any) {
	m.Errorf(format, args...)
	m.FailNow()
}
func (m ExampleTestRunner) Helper() {}

// MARK: MockTestRunner

type MockTestRunner struct {
	failed bool
	output []string
}

func (m *MockTestRunner) Error(args ...any) {
	m.output = append(m.output, fmt.Sprint(args...))
	m.failed = true
}

func (m MockTestRunner) Errorf(format string, args ...any) {
	m.Error(fmt.Sprintf(format, args...))
}

func (m *MockTestRunner) Fail() {
	m.failed = true
}
func (m MockTestRunner) FailNow() {
	m.Fail()
	runtime.Goexit()
}
func (m MockTestRunner) Fatal(args ...any) {
	m.Error(args...)
	m.FailNow()
}
func (m MockTestRunner) Fatalf(format string, args ...any) {
	m.Errorf(format, args...)
	m.FailNow()
}
func (m MockTestRunner) Helper() {}

func (m MockTestRunner) Run(name string, f func(t *testing.T)) bool {
	f(nil)
	return false
}

func (m MockTestRunner) Output() []string {
	return m.output
}
func (m MockTestRunner) Failed() bool {
	return m.failed
}

// MARK: T()

// T retrieves the *testing.T for the calling test frame.
func T() TestRunner {
	t, _ := TestFrame()
	if t == nil {
		panic("T: no test frame found; did you call With(t)?")
	}
	return t
}

// GetT retrieves the *testing.T for the calling test frame,
// by calling T().
//
// This is used internally where calling the T() function directly
// is not possible e.g. due to a name collision with a generic type
// parameter.
func GetT() TestRunner {
	return T()
}
