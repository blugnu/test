package test_test

import (
	"errors"
	"fmt"
	"io/fs"
	"runtime"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/matchers/panics"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

// MARK: errors

type MyError struct {
	msg string
}

func (e MyError) Error() string {
	return e.msg
}

func TestBeError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "error is of expected type",
			Act: func() {
				err := fmt.Errorf("an error occurred: %w", MyError{"invalid argument"})
				Expect(err).To(BeError[MyError]())
			},
		},

		{Scenario: "error is not of expected type",
			Act: func() {
				err := errors.New("some other error")
				Expect(err).To(BeError[MyError]())
			},
			Assert: func(result *R) {
				result.Expect(
					"expected: test_test.MyError",
					"got     : *errors.errorString",
				)
			},
		},

		{Scenario: "nil is not error",
			Act: func() { Expect(error(nil)).To(BeError[error]()) },
			Assert: func(result *R) {
				result.Expect(
					"expected: error",
					"got     : nil",
				)
			},
		},

		{Scenario: "too many errors",
			Act: func() {
				err1 := errors.New("error one")
				err2 := errors.New("error two")
				Expect(err1).To(BeError(err1, err2))
			},
			Assert: func(result *R) {
				result.ExpectInvalid("BeError: at most one error argument is supported, got 2")
			},
		},
	}...))
}

// MARK: panics

func TestNilPanic(t *testing.T) {
	With(t)

	result := NilPanic()
	Expect(result).To(Equal(panics.Expected{R: &runtime.PanicNilError{}}))
}

func TestPanic(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "with no args",
			Act: func() {
				result := Panic()
				Expect(result).To(Equal(panics.Expected{}))
			},
		},
		{Scenario: "with nil arg",
			Act: func() {
				result := Panic(nil)
				Expect(result).To(Equal(panics.Expected{R: opt.NoPanicExpected(true)}))
			},
		},
		{Scenario: "with non-nil arg",
			Act: func() {
				result := Panic("some string")
				Expect(result).To(Equal(panics.Expected{R: "some string"}))
			},
		},
		{Scenario: "with multiple args",
			Act: func() {
				result := Panic("one", "two")
				Expect(result).To(Equal(panics.Expected{}))
			},
			Assert: func(result *R) {
				result.ExpectInvalid("Panic: expected at most one argument, got 2")
			},
		},
	}...))
}

// =================================================================
// MARK: examples
// =================================================================

func ExampleBeError() {
	test.Example()

	err := fmt.Errorf("an error occurred: %w", MyError{"invalid argument"})

	// this test will pass:
	Expect(err).To(BeError[MyError]())
	Expect(err).ToNot(BeError[*fs.PathError]())

	// these tests will fail:
	Expect(error(nil)).To(BeError[error]())
	Expect(err).To(BeError[*fs.PathError]())

	// Output:
	// expected: error
	// got     : nil
	//
	// expected: *fs.PathError
	// got     : *fmt.wrapError
}

func ExamplePanic() {
	test.Example()

	// panic tests must be deferred so that the panic can be recovered;
	//
	// a stack trace is included by default but is disabled here to
	// avoid breaking the example output
	defer Expect(Panic("some string")).DidOccur(opt.NoStackTrace())

	// cause a panic; note that the value is different to the value
	// expected to be recovered (established above), so the test will fail
	panic("some other string")

	// Output:
	// unexpected panic:
	//   expected : "some string"
	//   recovered: "some other string"
}
