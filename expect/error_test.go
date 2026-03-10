package expect_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/expect"
)

type MyError struct {
	msg string
}

func (e MyError) Error() string {
	return e.msg
}

func TestError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "non-nil error without target",
			Act: func() {
				err := errors.New("an error occurred")
				expect.Error(err)
			},
		},

		{Scenario: "nil error without target",
			Act: func() {
				var err error
				expect.Error(err)

				Error("subsequent failures appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error",
				)
				result.Expect(
					"subsequent failures appear in report",
				)
			},
		},

		{Scenario: "error matches target",
			Act: func() {
				target := errors.New("target error")
				err := target

				expect.Error(err, target)
			},
		},
		{Scenario: "error does not match target",
			Act: func() {
				target := errors.New("target error")
				err := errors.New("some other error")

				expect.Error(err, target)
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: \"target error\" [*errors.errorString]",
					"got           : \"some other error\"",
				)
			},
		},
	}...))
}

func TestErrorAs(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "error is of expected type",
			Act: func() {
				err := fmt.Errorf("an error occurred: %w", MyError{"invalid argument"})
				e, ok := expect.ErrorAs[MyError](err)

				Expect(ok).To(BeTrue())
				Expect(e).To(Equal(MyError{"invalid argument"}))
			},
		},

		{Scenario: "error is not of expected type",
			Act: func() {
				err := errors.New("some other error")
				e, ok := expect.ErrorAs[MyError](err)
				Expect(ok).To(BeFalse())
				Expect(e).To(Equal(MyError{})) // zero value of MyError

				Error("subsequent failures should appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error to match or wrap: expect_test.MyError",
					"got: \"some other error\" [*errors.errorString]",
				)
				result.Expect(
					"subsequent failures should appear in report",
				)
			},
		},

		{Scenario: "nil should not be identified as any error",
			Act: func() {
				e, ok := expect.ErrorAs[error](error(nil))
				Expect(ok).To(BeFalse())
				Expect(e).IsNil()

				Error("subsequent failures should appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: error",
					"got           : nil",
				)
				result.Expect(
					"subsequent failures should appear in report",
				)
			},
		},
	}...))
}

func TestNoError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "nil error",
			Act: func() {
				var err error
				expect.NoError(err)
			},
		},

		{Scenario: "non-nil error",
			Act: func() {
				err := errors.New("an error occurred")
				expect.NoError(err)

				Error("subsequent failures appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected error",
					"got: \"an error occurred\" [*errors.errorString]",
				)
				result.Expect(
					"subsequent failures appear in report",
				)
			},
		},
	}...))
}
