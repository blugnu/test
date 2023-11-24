package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestExpectedPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("error")

	// ARRANGE
	testcases := []struct {
		name   string
		arg    error
		result *Panic
	}{
		{name: "nil argument", arg: nil, result: nil},
		{name: "non-nil argument", arg: err, result: &Panic{err}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			got := ExpectPanic(tc.arg)

			// ASSERT
			switch {
			case tc.result == nil:
				Equal(t, tc.result, got)
			case got != nil:
				t.Run("Panic with error", func(t *testing.T) {
					Equal(t, tc.result.error, got.error)
				})
			default:
				t.Errorf("\nwanted: %#v\ngot   : nil", tc.result)
			}
		})
	}
}

func TestPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("panic")
	werr := fmt.Errorf("wrapped: %w", err)

	testcases := []struct {
		name    string
		sut     *Panic
		fn      func()
		outcome any
		output  any
	}{
		{name: "nil receiver, no panic", fn: func() {},
			outcome: ShouldPass,
		},
		{name: "nil receiver, panicked", fn: func() { panic(err) },
			outcome: ShouldFail,
			output:  "unexpected panic: panic",
		},
		{name: "expected panic, no panic", sut: ExpectPanic(err), fn: func() {},
			outcome: ShouldFail,
			output: []string{
				"wanted (panic): panic",
				"got           : (did not panic)"},
		},
		{name: "expected panic, panicked", sut: ExpectPanic(err), fn: func() { panic(err) },
			outcome: ShouldPass,
		},
		{name: "expected panic, panicked with wrapped err", sut: ExpectPanic(err), fn: func() { panic(werr) },
			outcome: ShouldPass,
		},
		{name: "expected panic, panicked with other err", sut: ExpectPanic(err), fn: func() { panic(errors.New("other error")) },
			outcome: ShouldFail,
			output: []string{
				"wanted (panic): panic",
				"got    (panic): other error"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			stdout, _ := Helper(t, func(st *testing.T) {
				defer tc.sut.IsRecovered(st)
				tc.fn()
			}, tc.outcome)

			// ASSERT
			stdout.Contains(t, tc.output)
		})
	}
}
