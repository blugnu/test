package typecheck_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/blugnu/test"
)

func TestTo_MatchBeError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "error matches a specified target",
			Act: func() {
				err := errors.New("some other error")
				Expect(err).To(BeError(err))
			},
		},

		{Scenario: "error does not match target of same type",
			Act: func() {
				err := errors.New("error")
				target := errors.New("some other error")
				Expect(err).To(BeError(target))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: \"some other error\" [*errors.errorString]",
					"got           : \"error\"",
				)
			},
		},

		{Scenario: "errors of different types; expected Error shorter",
			Act: func() {
				err := errors.New("error")
				target := &json.SyntaxError{}
				Expect(err).To(BeError(target))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: \"\"      [*json.SyntaxError]",
					"got           : \"error\" [*errors.errorString]",
				)
			},
		},

		{Scenario: "errors of different types; expected Error longer",
			Act: func() {
				var (
					err    error = &json.SyntaxError{}
					target error = errors.New("error")
				)
				Expect(err).To(BeError(target))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: \"error\" [*errors.errorString]",
					"got           : \"\"      [*json.SyntaxError]",
				)
			},
		},

		{Scenario: "error is nil",
			Act: func() {
				target := errors.New("error")
				err := error(nil)
				Expect(err).To(BeError(target))
			},
			Assert: func(result *R) {
				result.Expect(
					"expected error: \"error\" [*errors.errorString]",
					"got           : nil",
				)
			},
		},

		{Scenario: "error and target are both nil",
			Act: func() {
				err := error(nil)
				target := error(nil)
				Expect(err).To(BeError(target))
			},
		},

		{Scenario: "target is nil",
			Act: func() {
				err := errors.New("error")
				target := error(nil)
				Expect(err).To(BeError(target))
			},
			Assert: func(result *R) {
				result.ExpectInvalid("target error must be non-nil")
			},
		},
	}...))
}

func TestToNot_MatchBeError(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "error matches a specified target",
			Act: func() {
				err := errors.New("some other error")
				Expect(err).ToNot(BeError(err))
			},
			Assert: func(result *R) {
				result.Expect(
					"should not be: \"some other error\" [*errors.errorString]",
				)
			},
		},
		{Scenario: "error does not match a specified target",
			Act: func() {
				err := errors.New("error")
				other := errors.New("some other error")
				Expect(err).ToNot(BeError(other))
			},
		},
	}...))
}
