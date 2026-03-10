package require_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/require"
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
		{Scenario: "error is of expected type",
			Act: func() {
				// arrange
				err := fmt.Errorf("an error occurred: %w", MyError{"invalid argument"})

				// act
				e := require.Error[MyError](err)

				// assert
				Expect(e).To(Equal(MyError{"invalid argument"}))
				Error("subsequent failures should appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"subsequent failures should appear in report",
				)
			},
		},

		{Scenario: "error is not of expected type",
			Act: func() {
				// arrange
				err := errors.New("some other error")

				// act
				_ = require.Error[MyError](err)
				Error("subsequent failures should not appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error to match or wrap: require_test.MyError",
					"got: \"some other error\" [*errors.errorString]",
				)
			},
		},

		{Scenario: "nil should not be identified as any error",
			Act: func() {
				// act
				_ = require.Error[error](error(nil))
				Error("subsequent failures should not appear in report")
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: error",
					"got           : nil",
				)
			},
		},
	}...))
}
