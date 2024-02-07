package test

import (
	"testing"
)

func TestEquality(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		sut      Equality
		result   string
	}{
		{scenario: "test.ShallowEquality", sut: ShallowEquality, result: "test.ShallowEquality"},
		{scenario: "test.DeepEquality", sut: DeepEquality, result: "test.DeepEquality"},
		{scenario: "invalid (-1)", sut: Equality(-1), result: "test.Equality(-1)"},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			got := tc.sut.String()

			// ASSERT
			Equal(t, got, tc.result)
		})
	}
}
