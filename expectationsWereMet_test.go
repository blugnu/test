package test

import (
	"errors"
	"testing"
)

type mock struct {
	expectationsWereMet_wasCalled bool
	reset_wasCalled               bool
}

func (mock *mock) ExpectationsWereMet() error {
	mock.expectationsWereMet_wasCalled = true
	return errors.New("mock expectation error")
}

func (mock *mock) Reset() {
	mock.reset_wasCalled = true
}

func TestExpectationsWereMet(t *testing.T) {
	With(t)

	// ARRANGE
	mock := &mock{}

	// ACT
	result := Test(func() { ExpectationsWereMet(mock) })

	// ASSERT
	Expect(mock.expectationsWereMet_wasCalled, "expectationsWereMet was called").Is(true)
	Expect(mock.reset_wasCalled, "reset was called").Is(true)

	result.Expect(
		"mock expectation error",
	)
}
