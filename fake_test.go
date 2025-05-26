package test

import (
	"errors"
	"testing"
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

	type testcase struct {
		scenario string
		exec     func()
	}
	RunScenarios(
		func(tc *testcase, _ int) {
			tc.exec()
		},
		[]testcase{
			{scenario: "returns value",
				exec: func() {
					// ARRANGE
					sut := FakeResult[int]{}

					// ACT
					sut.Returns(42)

					// ASSERT
					Expect(sut).To(Equal(FakeResult[int]{Result: 42}))
				},
			},
			{scenario: "returns error",
				exec: func() {
					// ARRANGE
					sut := FakeResult[int]{}
					err := errors.New("fake error")

					// ACT
					sut.Returns(err)

					// ASSERT
					Expect(sut).To(Equal(FakeResult[int]{Err: err}))
				},
			},
			{scenario: "returns result and error",
				exec: func() {
					// ARRANGE
					sut := FakeResult[int]{}
					err := errors.New("fake error")

					// ACT
					sut.Returns(42, err)

					// ASSERT
					Expect(sut).To(Equal(FakeResult[int]{Result: 42, Err: err}))
				},
			},
			{scenario: "multiple result values",
				exec: func() {
					// ARRANGE + ASSERT
					sut := FakeResult[int]{}
					defer Expect(Panic(ErrInvalidOperation)).DidOccur()

					// ACT
					sut.Returns(1, 2)
				},
			},
			{scenario: "multiple error values",
				exec: func() {
					// ARRANGE + ASSERT
					sut := FakeResult[int]{}
					err := errors.New("fake error")
					defer Expect(Panic(ErrInvalidOperation)).DidOccur()

					// ACT
					sut.Returns(err, err)
				},
			},
			{scenario: "invalid type",
				exec: func() {
					// ARRANGE + ASSERT
					sut := FakeResult[int]{}
					defer Expect(Panic(ErrInvalidOperation)).DidOccur()

					// ACT
					sut.Returns("invalid type")
				},
			},
		})
}
