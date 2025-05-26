package test

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: revisit the naming here

func Run(name string, fn func()) {
	t := T()
	t.Helper()

	t.Run(name, func(st *testing.T) {
		With(st)

		st.Helper()
		fn()
	})
}

func getScenarioName(ref reflect.Value) string {
	getField := func(fname string) string {
		if f := reflect.Indirect(ref).FieldByName(fname); f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
		return ""
	}

	candidate := []string{"Scenario", "scenario", "Name", "name"}
	for _, c := range candidate {
		if name := getField(c); name != "" {
			return name
		}
	}

	return ""
}

// RunScenarios is used for running table-driven tests using a function and
// a slice of scenarios.
//
// Each scenario is run in its own test runner; the name for each test is
// determined as follows:
//
//  1. If the scenario is a struct and has a field named "Scenario" or "Name"
//     (or an unexported equivalent), that field is used as the name for the test.
//
//  2. If the scenario is not a struct or does not have a supported name field,
//     the test is named "testcase/000" where 000 is the index of the scenario
//     in the slice.
//
// The function f is called with a pointer to a scenario as its argument.
func RunScenarios[T any](f func(*T, int), scns []T) {
	t := GetT()
	t.Helper()

	for num := range scns {
		scn := scns[num]

		// use reflection to determine if arg is a struct with a name field;
		// name fields used in order of preference:
		// 1. Scenario
		// 2. scenario
		// 3. Name
		// 4. name
		name := ""
		if ref := reflect.ValueOf(scn); ref.Kind() == reflect.Struct {
			name = getScenarioName(ref)
		}
		if name == "" {
			name = fmt.Sprintf("testcase/%.3d", num)
		}

		t.Run(name, func(t *testing.T) {
			With(t)
			t.Helper()
			f(&scn, num)
		})
	}
}
