package test

import "testing"

func TestTestable(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "errorf/no args",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := testable[string]{}

				// ACT
				result := Helper(t, func(t *testing.T) {
					sut.errorf(t, "error message")
				})

				// ASSERT
				result.Report.ContainsMatch("^[\\s]+error message$")
			},
		},
		{scenario: "errorf/with args",
			exec: func(t *testing.T) {
				// ARRANGE
				sut := testable[string]{}

				// ACT
				result := Helper(t, func(t *testing.T) {
					sut.errorf(t, "error message %s", "with arg")
				})

				// ASSERT
				result.Report.ContainsMatch("^[\\s]+error message with arg$")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
