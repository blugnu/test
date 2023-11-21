package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorIs(t *testing.T) {
	// ARRANGE
	a := errors.New("error a")
	b := errors.New("error b")
	wa := fmt.Errorf("wrapped: %w", a)

	testcases := []struct {
		name    string
		sut     func(*testing.T)
		outcome HelperResult
		output  any
	}{
		// execpted to pass
		{name: "ErrorIs, got == wanted", sut: func(st *testing.T) { ErrorIs(st, a, a) }, outcome: ShouldPass},
		{name: "ErrorIs, got wraps wanted", sut: func(st *testing.T) { ErrorIs(st, a, wa) }, outcome: ShouldPass},
		{name: "UnexpectedError, nil", sut: func(st *testing.T) { UnexpectedError(st, nil, ErrorDecl) }, outcome: ShouldPass},

		// expected to fail
		{name: "UnexpectedError, non-nil, default format", sut: func(st *testing.T) { UnexpectedError(st, errors.New("unexpected error")) },
			outcome: ShouldFail,
			output:  "unexpected error: unexpected error",
		},
		{name: "got != wanted", sut: func(st *testing.T) { ErrorIs(st, a, b) },
			outcome: ShouldFail,
			output: []string{
				"wanted error: error a",
				"got         : error b",
			},
		},
		{name: "got != wanted (wanted == nil)", sut: func(st *testing.T) { ErrorIs(st, nil, b) },
			outcome: ShouldFail,
			output:  "unexpected error: error b",
		},
		{name: "got != wanted, ErrorDefault", sut: func(st *testing.T) { ErrorIs(st, a, b, ErrorDefault) },
			outcome: ShouldFail,
			output: []string{
				"wanted error: error a",
				"got         : error b",
			},
		},
		{name: "got != wanted, ErrorString", sut: func(st *testing.T) { ErrorIs(st, a, b, ErrorString) },
			outcome: ShouldFail,
			output: []string{
				"wanted error: error a",
				"got         : error b",
			},
		},
		{name: "got != wanted, ErrorDecl", sut: func(st *testing.T) { ErrorIs(st, a, b, ErrorDecl) },
			outcome: ShouldFail,
			output: []string{
				"wanted error: &errors.errorString{s:\"error a\"}",
				"got         : &errors.errorString{s:\"error b\"}"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			stdout, _ := Helper(t, func(st *testing.T) {
				tc.sut(st)
			}, tc.outcome)

			// ASSERT
			stdout.Contains(t, tc.output)
		})
	}
}
