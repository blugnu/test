package report_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/expect"
	"github.com/blugnu/test/report"
)

func TestExpected(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Expected with a simple string",
			Act: func() {
				result := report.Expected("42")
				expect.Slice(result).Equals([]string{"expected: 42"})
			},
		},
	}...))
}

func TestExpectedGot(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "ExpectedGot with simple strings",
			Act: func() {
				result := report.ExpectedGot("42", "43")
				expect.Slice(result).Equals([]string{
					"expected: 42",
					"got     : 43",
				})
			},
		},
	}...))
}
