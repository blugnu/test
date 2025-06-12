package test

import (
	"fmt"
	"reflect"
	"testing"
)

// FUTURE: review and finalise naming

// Run is used to run a subtest with a given name and function.
func Run(name string, fn func()) {
	t := T()
	t.Helper()

	t.Run(name, func(st *testing.T) {
		With(st)

		st.Helper()
		fn()
	})
}

// RunScenarios is used for running table-driven tests using a function and
// a slice of scenarios.  The function accepts a pointer to a scenario and
// an index, allowing it to modify the scenario or use the index for reporting.
//
// Each scenario is run in its own test runner; the name for each test is
// determined as follows:
//
//  1. If the scenario is a struct and has a field named "Scenario" or "Name"
//     (or an unexported equivalent), that field is used as the name for the test.
//
//  2. If the scenario is not a struct or does not have a supported name field,
//     the test is named "000" where 000 is the 0-based index of the scenario
//     in the slice, left-padded with zeroes.
func RunScenarios[T any](f func(*T, int), scns []T) {
	t := GetT()
	t.Helper()

	for num := range scns {
		scn := scns[num]

		name := getScenarioName(scn, num)

		t.Run(name, func(t *testing.T) {
			With(t)
			t.Helper()
			f(&scn, num)
		})
	}
}

// getScenarioName derives a name for a given scenario.
//
// It checks for a field named "Scenario", "scenario", "Name", or "name"
// in the struct and returns its value as a string. If no such field is found,
// it returns a default name in the format "testcase/000" where 000 is the
// index of the scenario in the slice.
func getScenarioName(scn any, idx int) string {
	result := fmt.Sprintf("testcase-%.3d", idx)

	// if the scenario is not a struct (or pointer to a struct), return
	// the default result
	ref := reflect.Indirect(reflect.ValueOf(scn))
	if ref.Kind() != reflect.Struct {
		return result
	}

	// otherwise, check for supported fields that may contain a usable name
	candidate := []string{"Scenario", "scenario", "Name", "name"}
	for _, c := range candidate {
		if f := reflect.Indirect(ref).FieldByName(c); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}

	return result
}
