package test_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/test"
)

func TestOutcome(t *testing.T) {
	With(t)

	type testcase struct {
		test.Outcome
		result string
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			// ACT
			result := tc.String()

			// ASSERT
			Expect(result).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{Outcome: test.Passed, result: "test.Passed"},
			{Outcome: test.Failed, result: "test.Failed"},
			{Outcome: test.Panicked, result: "test.Panicked"},
			{Outcome: 99, result: "test.Outcome(99)"},
		}),
	))
}
