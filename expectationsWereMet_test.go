package test

import (
	"errors"
	"testing"
)

type mockExpectationsWereMet struct {
	ExpectationsWereMetWasCalled bool
	ResetWasCalled               bool
	error
}

func (mock *mockExpectationsWereMet) ExpectationsWereMet() error {
	mock.ExpectationsWereMetWasCalled = true
	return mock.error
}

func (mock *mockExpectationsWereMet) Reset() {
	mock.ResetWasCalled = true
}

func TestExpectationsWereMet(t *testing.T) {
	// ARRANGE
	mock1 := &mockExpectationsWereMet{error: errors.New("mock1 error")}
	mock2 := &mockExpectationsWereMet{error: errors.New("mock2 error")}

	// ACT
	test := Helper(t, func(st *testing.T) { ExpectationsWereMet(st, mock1, mock2) })

	// ASSERT
	Bool(t, mock1.ExpectationsWereMetWasCalled).IsTrue()
	Bool(t, mock2.ExpectationsWereMetWasCalled).IsTrue()
	Bool(t, mock1.ResetWasCalled).IsTrue()
	Bool(t, mock2.ResetWasCalled).IsTrue()

	test.DidFail()
	test.Report.Contains([]string{
		"mock1 error",
		"mock2 error",
	})
}
