package test

import (
	"testing"

	"github.com/blugnu/test/internal/testframe"
)

type TestingT interface {
	Cleanup(fn func())
	Name() string
	Run(name string, fn func(t *testing.T)) bool
	Error(args ...any)
	Errorf(s string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(s string, args ...any)
	Helper()
	Parallel()
	Setenv(name string, value string)
	SkipNow()
}

// MARK: T()

// GetT retrieves the *testing.T for the calling test frame, by calling T().
//
// GetT is provided for use where calling the T() function directly is not
// possible, e.g. due to a name collision with a generic type parameter.
func GetT() TestingT {
	return T()
}

// T retrieves the TestRunner for the calling test frame. When running in a test
// frame, this will return the *testing.T for the test.
//
// The T is returned as a TestingT interface; this provides all of the
// functionality of the *testing.T type, but allows for more flexibility
// in the test package.
func T() TestingT {
	return testframe.MustPeek[TestingT]()
}

// MARK: With()

// With pushes the given TestingT onto the test frame stack; if the TestingT is
// not nil it will be popped from the stack when the test has completed.
//
// This is used to set the current test frame for the test package, typically
// called as the first line of a Test...() function:
//
//	func TestSomething(t *testing.T) {
//	    With(t)
//
//	    // ... rest of the test code ...
//	}
//
// If `blugnu/test` functions are used to run subtests etc, no further calls to With()
// are required in a test function; the test frame will be automatically managed
// by the test package.
//
// If a new test frame is explicitly created, e.g. by calling t.Run(string, func(t *testing.T)),
// then With(t) must be called to push the new test frame onto the stack.  Again, this
// will be automatically popped from the stack when the new test frame completes:
//
//	func TestSomething(t *testing.T) {
//	    With(t)
//
//	    // when using test package functions to run subtests you do not
//	    // need to call With() again
//
//	    Run(Test("subtest", func() {
//	        // ... rest of the subtest code ...
//	    })
//
//	    // but With() must be called if a new test frame is explicitly created
//
//	    t.Run("subtest", func(t *testing.T) {
//	        With(t)
//	        // ... rest of the subtest code ...
//	    })
//
//	    // ... rest of the test code ...
//	}
//
// To simultaneously push a test frame and mark it for parallel execution,
// you can use the Parallel() function:
//
//	func TestSomething(t *testing.T) {
//	    Parallel(t)
//
//	    // ... rest of the test code ...
//	}
func With(t TestingT) {
	if t == nil {
		panic(testframe.ErrNoTestFrame)
	}

	testframe.Push(t)

	t.Cleanup(func() {
		testframe.Pop()
	})
}
