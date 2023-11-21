package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestExpectPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("error")

	// ACT
	got := ExpectPanic(err)

	// ASSERT
	t.Run("returns", func(t *testing.T) {
		if got, ok := Type[*ExpectedPanic](t, got); ok {
			t.Run("with error", func(t *testing.T) {
				Equal(t, err, got.error)
			})
		}
	})
}

func TestExpectedPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("panic")
	werr := fmt.Errorf("wrapped: %w", err)

	testcases := []struct {
		name    string
		sut     *ExpectedPanic
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
				defer tc.sut.Assert(st)
				tc.fn()
			}, tc.outcome)

			// ASSERT
			stdout.Contains(t, tc.output)
		})
	}
}
