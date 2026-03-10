package typecheck_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"
)

type MyError struct {
	msg string
}

func (e MyError) Error() string {
	return e.msg
}

func TestTo_BeErrorType(t *testing.T) {
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
					"expected: typecheck_test.MyError",
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
	}...))
}

func TestToNot_BeError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "error is not of the expected type",
			Act: func() {
				err := errors.New("some other error")
				Expect(err).ToNot(BeError[MyError]())
			},
		},

		{Scenario: "error is of the expected type",
			Act: func() {
				err := fmt.Errorf("an error occurred: %w", MyError{"invalid argument"})
				Expect(err).ToNot(BeError[MyError]())
			},
			Assert: func(result *R) {
				result.Expect(
					"should not be of type: typecheck_test.MyError",
				)
			},
		},

		{Scenario: "nil is not of any error type",
			Act: func() { Expect(error(nil)).ToNot(BeError[error]()) },
		},
	}...))
}
