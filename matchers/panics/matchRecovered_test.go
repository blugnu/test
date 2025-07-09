package panics_test

import (
	"errors"
	"strings"
	"testing"

	. "github.com/blugnu/test"
)

func TestPanic_DidOccur(t *testing.T) {
	With(t)

	Run(HelperTests([]HelperScenario{
		{Scenario: "did not panic when panic was expected",
			Act: func() {
				err := errors.New("panic with error")
				defer Expect(Panic(err)).DidOccur()
			},
			Assert: func(result *R) {
				result.Expect(
					"expected panic: *errors.errorString(panic with error)",
					"  recovered   : nil (did not panic)",
				)
			},
		},
		{Scenario: "panicked when panic was implictly unexpected",
			Act: func() {
				panic("panics with string")
			},
			Assert: func(result *R) {
				result.Expect(TestPanicked, "panics with string")
			},
		},
		{Scenario: "panicked with unexpected error",
			Act: func() {
				err := errors.New("panic with error")
				defer Expect(Panic(err)).DidOccur()

				panic(errors.New("panic with different error"))
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  expected : *errors.errorString(panic with error)",
					"  recovered: *errors.errorString(panic with different error)",
				)
			},
		},
		{Scenario: "panicked with expected non-error value",
			Act: func() {
				defer Expect(Panic(1)).DidOccur()
				panic(2)
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  expected : int(1)",
					"  recovered: int(2)",
				)
			},
		},
		{Scenario: "panicked with expected non-error value",
			Act: func() {
				defer Expect(Panic(1)).DidOccur()
				panic(1)
			},
		},
		{Scenario: "custom comparison function",
			Act: func() {
				defer Expect(Panic("x")).DidOccur(func(expected, got any) bool {
					exps, _ := expected.(string)
					gots, _ := got.(string)
					return exps == strings.ToLower(gots)
				})

				panic("X")
			},
		},
		{Scenario: "expecting value of non-assignable type",
			Act: func() {
				defer Expect(Panic(1)).DidOccur()

				panic("1")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  expected : int(1)",
					"  recovered: string(\"1\")",
				)
			},
		},

		// Panic() tests
		{Scenario: "Panic().DidOccur(), panicked",
			Act: func() {
				defer Expect(Panic()).DidOccur()
				panic("foo")
			},
		},
		{Scenario: "Panic().DidOccur(), did not panic",
			Act: func() {
				defer Expect(Panic()).DidOccur()
			},
			Assert: func(result *R) {
				result.Expect(
					"expected panic: <any value recovered>",
					"  recovered   : nil (did not panic)",
				)
			},
		},
		{Scenario: "Panic().DidNotOccur(), did not panic",
			Act: func() {
				defer Expect(Panic()).DidNotOccur()
			},
		},
		{Scenario: "Panic().DidNotOccur(), panicked",
			Act: func() {
				defer Expect(Panic()).DidNotOccur()
				panic("foo")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  recovered: string(\"foo\")",
				)
			},
		},

		// Panic(nil) tests
		{Scenario: "Panic(nil).DidOccur(), panicked",
			Act: func() {
				defer Expect(Panic(nil)).DidOccur()
				panic("foo")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					"  recovered: string(\"foo\")",
				)
			},
		},
		{Scenario: "Panic(nil).DidOccur(), did not panic",
			Act: func() {
				defer Expect(Panic(nil)).DidOccur()
			},
		},
		{Scenario: "Panic(nil).DidNotOccur(), did not panic",
			Act: func() {
				defer Expect(Panic(nil)).DidNotOccur()
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"DidNotOccur: may not be used with Panic(nil); did you mean NilPanic()?",
				)
			},
		},
		{Scenario: "Panic(nil).DidNotOccur(), panicked",
			Act: func() {
				defer Expect(Panic(nil)).DidNotOccur()
				panic("foo")
			},
			Assert: func(result *R) {
				result.ExpectInvalid(
					"DidNotOccur: may not be used with Panic(nil); did you mean NilPanic()?",
				)
			},
		},

		// NilPanic tests
		{Scenario: "NilPanic() panicked",
			// from go1.21 onwards, panic(nil) triggers a runtime.PanicNilError;
			// Expect(Panic(nil)) tests specifically for this case.  To test that
			// no panic occurs (rather than allowing the panic to vomit all over
			// the test output) use Expect(NoPanic()).DidOccur()
			//
			// https://go.dev/blog/compat#expanded-godebug-support-in-go-121
			Act: func() {
				defer Expect(NilPanic()).DidOccur()

				// these contortions are just to side-step the nilpanic lint that an
				// explicit panic(nil) would trigger
				nz := func() any { return nil }()

				panic(nz)
			},
		},
		{Scenario: "NilPanic() does not panic",
			Act: func() {
				defer Expect(NilPanic()).DidOccur()
			},
			Assert: func(result *R) {
				result.Expect(
					"expected panic: *runtime.PanicNilError",
					"  recovered   : nil (did not panic)",
				)
			},
		},

		// DidNotOccur tests
		{Scenario: "did not panic when no panic was expected",
			Act: func() {
				defer Expect(Panic()).DidNotOccur()
			},
		},
		{Scenario: "panicked when no panic was expected",
			Act: func() {
				defer Expect(Panic()).DidNotOccur()

				panic("panics with string")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					`  recovered: string("panics with string")`,
				)
			},
		},
		{Scenario: "panicked with x, not expected to occur",
			Act: func() {
				defer Expect(Panic("x")).DidNotOccur()

				panic("x")
			},
			Assert: func(result *R) {
				result.Expect(
					`expected: panic with string("x"): should not have occurred`,
				)
			},
		},
		{Scenario: "panicked with x, panic with y not expected to occur",
			Act: func() {
				defer Expect(Panic("y")).DidNotOccur()

				panic("x")
			},
			Assert: func(result *R) {
				result.Expect(
					"unexpected panic:",
					`  recovered: string("x")`,
				)
			},
		},
	}...))
}
