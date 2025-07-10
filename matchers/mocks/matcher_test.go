package mocks_test

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

func TestMatcher(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "mock supports ExpectationsWereMet and Reset",
			Act: func() {
				// arrange
				mock := &mock{
					err: errors.New("mock expectations were not met\nwith several lines of error"),
				}

				// act
				Expect(mock, "my mock").Should(MeetExpectations())

				// assert
				Expect(mock.expectationsWereMet_wasCalled).Is(true)
				Expect(mock.reset_wasCalled).Is(true)
			},
			Assert: func(result *R) {
				result.Expect(
					"my mock:",
					"mock expectations were not met",
					"with several lines of error",
				)
			},
		},
		{Scenario: "mock is not supported",
			Act: func() {
				// arrange
				type mock struct{}

				// act
				Expect(mock{}, "my mock").Should(MeetExpectations())
			},
			Assert: func(result *R) {
				result.ExpectInvalid("mocks.Matcher: mocks_test.mock does not implement a supported mock interface")
			},
		},
	}...))
}
