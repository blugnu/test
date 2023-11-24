package test

import (
	"errors"
	"os"
	"testing"
)

func TestHelper(t *testing.T) {
	t.Run("matchAll", func(t *testing.T) {
		ok, err := matchAll("any", "none")
		if !ok || err != nil {
			t.Errorf("matchAll() = %v, %v, want %v, %v", ok, err, true, nil)
		}
	})

	// ARRANGE
	panerr := errors.New("panic err")

	type args struct {
		fn      func(*testing.T)
		outcome any
	}
	testcases := []struct {
		name    string
		args           // args passed to the Helper() function under test
		outcome any    // expected outcome of the Helper() used to test Helper()
		panic   *Panic // expected panic, if any
		output  any    // expected output
	}{
		{name: "invalid outcome", args: args{func(st *testing.T) {}, 42}, outcome: ExpectPanic(ErrInvalidArgument)},
		{name: "fail, no outcome", args: args{fn: func(st *testing.T) { st.Fail() }}, outcome: ShouldPass},
		{name: "pass, no outcome", args: args{fn: func(st *testing.T) {}}, outcome: ShouldPass},
		{name: "panic, no outcome", args: args{fn: func(st *testing.T) { panic(42) }}, outcome: ShouldPass},
		{name: "should fail", args: args{func(st *testing.T) { st.Fail() }, ShouldFail}, outcome: ShouldPass},
		{name: "should fail (bool)", args: args{func(st *testing.T) { st.Fail() }, false}, outcome: ShouldPass},
		{name: "should pass", args: args{func(st *testing.T) {}, ShouldPass}, outcome: ShouldPass},
		{name: "should pass (bool)", args: args{func(st *testing.T) {}, true}, outcome: ShouldPass},
		{name: "helper panics (expected err)", args: args{func(st *testing.T) { panic(panerr) }, ExpectPanic(panerr)}, outcome: ShouldPass},
		{name: "helper panics (unexpected err)", args: args{func(st *testing.T) { panic(errors.New("other error")) }, ExpectPanic(panerr)},
			outcome: ShouldFail,
			output: []string{
				"wanted (panic): panic err",
				"got    (panic): other error",
			},
		},
		{name: "should pass, panics", args: args{func(st *testing.T) { panic(panerr) }, ShouldPass},
			outcome: ShouldFail,
			output: []string{
				"wanted     : PASS",
				"got (panic): panic err",
			},
		},
		{name: "should fail, panics", args: args{func(st *testing.T) { panic(panerr) }, ShouldFail},
			outcome: ShouldFail,
			output: []string{
				"wanted     : FAIL",
				"got (panic): panic err",
			},
		},
		{name: "passed, should panic", args: args{func(st *testing.T) {}, ExpectPanic(panerr)},
			outcome: ShouldFail,
			output: []string{
				"wanted (panic): panic err",
				"got           : PASS",
			},
		},
		{name: "failed, should panic", args: args{func(st *testing.T) { st.Fail() }, ExpectPanic(panerr)},
			outcome: ShouldFail,
			output: []string{
				"wanted (panic): panic err",
				"got           : FAIL",
			},
		},
		{name: "verbose output from passing internal test is removed", args: args{func(st *testing.T) {
			// _simulates_ the additional output from a test that passes
			// when run with go test -v
			os.Stdout.WriteString("=== RUN   internal\n")
			os.Stdout.WriteString("--- PASS: internal (0.00s)\n")
		}, ShouldPass}, outcome: ShouldPass, output: ""},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			defer tc.panic.IsRecovered(t)

			// HERE BE DRAGONS: the Helper test is used to test the Helper test
			// (and capture the output from it)

			// ACT
			stdout, _ := Helper(t, func(st *testing.T) {
				Helper(st, tc.args.fn, tc.args.outcome)
			}, tc.outcome)

			// ASSERT
			stdout.Contains(t, tc.output)
		})
	}
}
