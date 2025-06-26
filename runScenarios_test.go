package test //nolint: testpackage // tests private functions

import (
	"testing"
)

func Test_getScenarioName(t *testing.T) {
	With(t)

	type structWith_Name struct {
		Name string
	}

	type structWith_name struct {
		name string
	}
	type structWith_Scenario struct {
		Scenario string
	}
	type structWith_scenario struct {
		scenario string
	}
	type structWith_scenario_and_name struct {
		scenario string
		name     string
	}

	var result string

	result = getScenarioName(nil, 1)
	Expect(result, "nil").To(Equal("testcase-001"))

	result = getScenarioName("no name", 1)
	Expect(result, "default").To(Equal("testcase-001"))

	result = getScenarioName(&structWith_Scenario{"name"}, 1)
	Expect(result, "Scenario").To(Equal("name"))

	result = getScenarioName(structWith_Scenario{"name"}, 1)
	Expect(result, "Scenario").To(Equal("name"))

	result = getScenarioName(structWith_scenario{"name"}, 1)
	Expect(result, "scenario").To(Equal("name"))

	result = getScenarioName(structWith_Name{"name"}, 1)
	Expect(result, "Name").To(Equal("name"))

	result = getScenarioName(structWith_name{"name"}, 1)
	Expect(result, "name").To(Equal("name"))

	result = getScenarioName(structWith_scenario_and_name{scenario: "scenario", name: "name"}, 1)
	Expect(result, "scenario (not name)").To(Equal("scenario"))
}
