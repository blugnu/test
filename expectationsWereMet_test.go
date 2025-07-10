package test_test

import (
	"errors"
	"testing"

	. "github.com/blugnu/test"
)

type mock struct {
	err                           error
	expectationsWereMet_wasCalled bool
	reset_wasCalled               bool
}

func (mock *mock) ExpectationsWereMet() error {
	mock.expectationsWereMet_wasCalled = true
	return mock.err
}

func (mock *mock) Reset() {
	mock.reset_wasCalled = true
}

func TestExpectationsWereMet(t *testing.T) {
	With(t)

	// ARRANGE
	mock := &mock{
		err: errors.New("mock expectation error"),
	}

	// ACT
	result := TestHelper(func() { ExpectationsWereMet(mock) })

	// ASSERT
	Expect(mock.expectationsWereMet_wasCalled, "expectationsWereMet was called").Is(true)
	Expect(mock.reset_wasCalled, "reset was called").Is(true)

	result.Expect(
		"mock expectation error",
	)
}

func TestMeetsExpectations(t *testing.T) {
	With(t)

	// ARRANGE
	mock := &mock{
		err: errors.New("mock expectation error\non multiple lines"),
	}

	// ACT
	result := TestHelper(func() { Expect(mock, "my mock").Should(MeetExpectations()) })

	// ASSERT
	Expect(mock.expectationsWereMet_wasCalled, "expectationsWereMet was called").Is(true)
	Expect(mock.reset_wasCalled, "reset was called").Is(true)

	result.Expect(
		"my mock:",
		"mock expectation error",
		"on multiple lines",
	)
}
