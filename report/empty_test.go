package report_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/report"
)

func TestEmpty(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "Empty with no type name",
			Act: func() {
				Expect(report.Empty()).To(Equal("<empty>"))
			},
		},
		{Scenario: "Empty with type name",
			Act: func() {
				Expect(report.Empty("string")).To(Equal("<empty string>"))
			},
		},
	}...))
}
