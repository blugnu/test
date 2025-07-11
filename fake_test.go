package test_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
)

func TestFakeResult(t *testing.T) {
	With(t)

	// ARRANGE
	sut := FakeResult[int]{Result: 42, Err: errors.New("faked error")}

	// ACT
	sut.Reset()

	// ASSERT
	Expect(sut).Is(FakeResult[int]{})
}

func TestFakeResult_Returns(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "returns value",
			Act: func() {
				// ARRANGE
				sut := FakeResult[int]{}

				// ACT
				sut.Returns(42)

				// ASSERT
				Expect(sut).To(Equal(FakeResult[int]{Result: 42}))
			},
		},
		{Scenario: "returns error",
			Act: func() {
				// ARRANGE
				sut := FakeResult[int]{}
				err := errors.New("fake error")

				// ACT
				sut.Returns(err)

				// ASSERT
				Expect(sut).To(Equal(FakeResult[int]{Err: err}))
			},
		},
		{Scenario: "returns result and nil",
			Act: func() {
				// ARRANGE
				sut := FakeResult[int]{}

				// ACT
				sut.Returns(42, nil)

				// ASSERT
				Expect(sut).To(Equal(FakeResult[int]{Result: 42}))
			},
		},
		{Scenario: "returns result and error",
			Act: func() {
				// ARRANGE
				sut := FakeResult[int]{}
				err := errors.New("fake error")

				// ACT
				sut.Returns(42, err)

				// ASSERT
				Expect(sut).To(Equal(FakeResult[int]{Result: 42, Err: err}))
			},
		},
		{Scenario: "multiple result values",
			Act: func() {
				// ARRANGE + ASSERT
				sut := FakeResult[int]{}

				// ACT
				sut.Returns(1, 2)
			},
			Assert: func(result *R) {
				result.ExpectInvalid(ErrInvalidOperation)
			},
		},
		{Scenario: "multiple error values",
			Act: func() {
				// ARRANGE + ASSERT
				sut := FakeResult[int]{}
				err := errors.New("fake error")

				// ACT
				sut.Returns(err, err)
			},
			Assert: func(result *R) {
				result.ExpectInvalid(ErrInvalidOperation)
			},
		},
		{Scenario: "invalid type",
			Act: func() {
				// ARRANGE + ASSERT
				sut := FakeResult[int]{}

				// ACT
				sut.Returns("invalid type")
			},
			Assert: func(result *R) {
				result.ExpectInvalid(ErrInvalidOperation)
			},
		},
	}...))
}
