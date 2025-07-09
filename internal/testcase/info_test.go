package testcase_test

import (
	"testing"

	. "github.com/blugnu/test"
	"github.com/blugnu/test/internal/testcase"
)

func Test_NameOrDefault(t *testing.T) {
	With(t)

	type withname struct{ name string }
	type withName struct{ Name string }
	type withscenario struct{ scenario string }
	type withScenario struct{ Scenario string }

	type tc struct {
		inst   any
		name   string
		index  int
		result string
	}
	Run(Testcases(
		ForEach(func(tc tc) {
			result := testcase.NameOrDefault(tc.inst, tc.name, tc.index)
			Expect(result).To(Equal(tc.result))
		}),
		Case("not a struct", tc{inst: "foo", result: "testcase-000"}),
		Case("with name (empty)", tc{inst: withname{}, index: 1, result: "testcase-001"}),
		Case("with name (+specified)", tc{inst: withname{"name"}, name: "specified", result: "specified"}),
		Case("with name (ptr)", tc{inst: &withname{"name"}, result: "name"}),
		Case("with name", tc{inst: withname{"name"}, result: "name"}),
		Case("with Name", tc{inst: withName{"Name"}, result: "Name"}),
		Case("with scenario", tc{inst: withscenario{"scenario"}, result: "scenario"}),
		Case("with Scenario", tc{inst: withScenario{"Scenario"}, result: "Scenario"}),
		Case("with whitespace", tc{inst: withName{" padded "}, result: "padded"}),
		Case("with whitespace", tc{inst: withName{" padded "}, result: "padded"}),
	))
}

func Test_IsDebugging(t *testing.T) {
	With(t)

	type withdebug struct{ debug bool }
	type withDebug struct{ Debug bool }

	type tc struct {
		inst   any
		result bool
	}

	Run(Testcases(
		ForEach(func(tc tc) {
			result := testcase.IsDebugging(tc.inst)
			Expect(result).To(Equal(tc.result))
		}),
		Case("not a struct", tc{inst: true, result: false}),
		Case("with debug (ptr)", tc{inst: &withdebug{true}, result: true}),
		Case("with debug", tc{inst: withdebug{true}, result: true}),
		Case("with Debug", tc{inst: withDebug{true}, result: true}),
	))
}

func Test_IsSkipping(t *testing.T) {
	With(t)

	type withskip struct{ skip bool }
	type withSkip struct{ Skip bool }

	type tc struct {
		inst   any
		result bool
	}

	Run(Testcases(
		ForEach(func(tc tc) {
			result := testcase.IsSkipping(tc.inst)
			Expect(result).To(Equal(tc.result))
		}),
		Case("not a struct", tc{inst: true, result: false}),
		Case("with skip (ptr)", tc{inst: &withskip{true}, result: true}),
		Case("with skip", tc{inst: withskip{true}, result: true}),
		Case("with Skip", tc{inst: withSkip{true}, result: true}),
	))
}
