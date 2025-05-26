package test

import (
	"strings"
	"testing"
)

func TestRunTestScenarios(t *testing.T) {
	With(t)

	RunTestScenarios([]TestScenario{
		{Scenario: "no scenarios",
			Act: func() {
				RunTestScenarios([]TestScenario{})
			},
		},
		{Scenario: "scenario is skipped",
			Act: func() {
				RunTestScenarios([]TestScenario{
					{Skip: true},
				})
			},
		},
		{Scenario: "scenario with no Act function",
			Act: func() {
				RunTestScenarios([]TestScenario{
					{Scenario: "no Act function"},
				})
			},
			Assert: func(result *R) {
				Expect(result.FailedTests).To(ContainItem("/no_Act_function"), strings.Contains)
				result.ExpectInvalid(
					"test.RunTestScenarios: no Act function defined",
				)
			},
		},
		{Scenario: "scenario does not test the test result",
			Act: func() {
				RunTestScenarios([]TestScenario{
					{Scenario: "result not tested",
						Act:    func() {},
						Assert: func(result *R) {},
					},
				})
			},
			Assert: func(result *R) {
				Expect(result.FailedTests).To(ContainItem("/result_not_tested"), strings.Contains)
				result.ExpectInvalid(
					"test.RunTestScenarios: result not tested; *R.Expect(...) or *R.ExpectInvalid(...) must be called",
				)
			},
		},
		{Scenario: "when debugging a scenario, other scenarios are not run",
			Act: func() {
				RunTestScenarios([]TestScenario{
					{Scenario: "not debugging",
						Act: func() {
							Expect(true).To(BeFalse())
						},
					},
					{Scenario: "debugging",
						Debug: true,
						Act: func() {
							Expect(true).To(BeTrue())
						},
					},
				})
			},
		},
	})
}
