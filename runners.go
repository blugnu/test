package test

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: revisit the naming here
// TODO: option or function to run scenarios in parallel

func Run(name string, fn func()) {
	t, _ := TestFrame()
	if t == nil {
		panic("Run: no test frame; did you forget to call With(t)?")
	}

	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()

		frameId := With(t)
		defer testFrames.Remove(frameId)

		fn()
	})
}

type Scenario struct {
	Scenario string
	Run      func()
}

func RunScenario(s Scenario) {
	s.Run()
}

func RunScenarios[T any](f func(T), scns []T) {
	t, _ := TestFrame()
	if t == nil {
		panic("RunScenarios: no test frame; did you forget to call With(t)?")
	}

	t.Helper()
	for num, arg := range scns {
		// use reflection to determine if arg is a struct with a Name field
		name := ""
		if ref := reflect.ValueOf(arg); ref.Kind() == reflect.Struct {
			if nameField := reflect.Indirect(ref).FieldByName("Scenario"); nameField.IsValid() {
				name = nameField.String()
			}
			if name == "" {
				if nameField := reflect.Indirect(ref).FieldByName("scenario"); nameField.IsValid() {
					name = nameField.String()
				}
			}
		}
		if name == "" {
			name = fmt.Sprintf("testcase/%.3d", num)
		}

		t.Run(name, func(t *testing.T) {
			t.Helper()

			frameId := With(t)
			defer testFrames.Remove(frameId)

			f(arg)
		})
	}
}
