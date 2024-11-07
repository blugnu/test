package test

import (
	"errors"
	"testing"
)

func TestFake(t *testing.T) {
	// ARRANGE
	fake := Fake[int]{Result: 42}

	// ACT
	fake.Reset()

	// ASSERT
	That(t, fake.Result).Equals(0)
}

func TestFakeReturns(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "returns value",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := Fake[int]{}

				// ACT
				sut.Returns(42)

				// ASSERT
				That(t, sut).Equals(Fake[int]{Result: 42})
			},
		},
		{scenario: "returns error",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := Fake[int]{}
				err := errors.New("fake error")

				// ACT
				sut.Returns(err)

				// ASSERT
				That(t, sut).Equals(Fake[int]{Err: err})
			},
		},
		{scenario: "returns result and error",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := Fake[int]{}
				err := errors.New("fake error")

				// ACT
				sut.Returns(42, err)

				// ASSERT
				That(t, sut).Equals(Fake[int]{Result: 42, Err: err})
			},
		},
		{scenario: "multiple result values",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := Fake[int]{}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.Returns(1, 2)
			},
		},
		{scenario: "multiple error values",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := Fake[int]{}
				err := errors.New("fake error")
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.Returns(err, err)
			},
		},
		{scenario: "invalid type",
			exec: func(t *testing.T) {
				// ARRANGE + ASSERT
				sut := Fake[int]{}
				defer ExpectPanic(ErrInvalidOperation).Assert(t)

				// ACT
				sut.Returns("invalid type")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
