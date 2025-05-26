package test

import (
	"reflect"
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

	result = getScenarioName(reflect.ValueOf(structWith_Scenario{"name"}))
	Expect(result, "Scenario").To(Equal("name"))

	result = getScenarioName(reflect.ValueOf(structWith_scenario{"name"}))
	Expect(result, "scenario").To(Equal("name"))

	result = getScenarioName(reflect.ValueOf(structWith_Name{"name"}))
	Expect(result, "Name").To(Equal("name"))

	result = getScenarioName(reflect.ValueOf(structWith_name{"name"}))
	Expect(result, "name").To(Equal("name"))

	result = getScenarioName(reflect.ValueOf(structWith_scenario_and_name{scenario: "scenario", name: "name"}))
	Expect(result, "scenario (not name)").To(Equal("scenario"))
}
