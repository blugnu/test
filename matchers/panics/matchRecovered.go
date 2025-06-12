package panics

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/blugnu/test/opt"
)

// MatchRecovered is a struct that implements the Matcher[Expected]
// interface.
//
// NOTE: This one is a bit different from the others since it is not
// intended to be used with To/ToNot.  Instead it is a shared
// implementation detail of the Expect().DidOccur/DidNotOccur methods.
//
// A panics.Expected value is used to capture the value that will be
// recovered from an expected panic:
//
//	Expect(Panic(x))    <-- expectation subject is a panics.Expected{}
//	                         with R = x
//
// A *MatchRecovered is created by the (deferred) Did/DidNotOccur()
// function, which captures the value that was actually recovered:
//
//	defer Expect(Panic(x)).DidOccur()  <-- creates a *MatchRecovered
//	                                        with R = recover()
//
// i.e. the value captured by Expect() is the "expected" value, rather
// than the subject (got). The subject is captured (i.e. recovered) by
// the Did/DidNotOccur function.
//
// This means that when the Match function is called, the expected and
// got values are reversed.
//
// To clarify, and to make both expected and got values available to the
// TestFailure function for reporting any failure, the Match() function
// unpacks the R values from the two Matchers and assigns them to the
// unexported got and expected fields of the receiver.
type MatchRecovered struct {
	R        any
	got      any
	expected any
}

type Expected struct {
	R any
}

func (pm *MatchRecovered) Match(target Expected, opts ...any) bool {
	pm.got = pm.R
	pm.expected = target.R

	switch {
	case pm.expected == opt.NoPanicExpected(true):
		// when expected is NoPanicExpected, we are not expecting any panic to have occurred
		// so we return true if the got value is nil, meaning no panic occurred
		return pm.got == nil

	case pm.expected == nil:
		return pm.got == nil

	default:
		if err, expectedErr := pm.expected.(error); expectedErr {
			if got, gotErr := pm.got.(error); gotErr {
				return errors.Is(got, err)
			}
			return false
		}
		expectType := reflect.TypeOf(pm.expected)
		gotType := reflect.TypeOf(pm.got)
		if expectType != gotType && gotType != nil && !gotType.AssignableTo(expectType) {
			return false
		}
		cmp := reflect.DeepEqual
		if fn, ok := opt.Get[func(any, any) bool](opts); ok {
			cmp = fn
		}

		return cmp(pm.expected, pm.got)
	}
}

// OnTestFailure returns a report of the failure for the matcher.
func (pm *MatchRecovered) OnTestFailure(opts ...any) []string {
	switch {
	case pm.got == nil:
		// when got is nil, some expected panic has failed to occur
		return []string{
			fmt.Sprintf("expected panic: %T(%v)", pm.expected, opt.ValueAsString(pm.expected, opts...)),
			"  recovered   : " + nilRecovered,
		}
	case pm.expected == nil || pm.expected == opt.NoPanicExpected(true):
		// when expected is nil we are not expecting any panic to have occurred
		return []string{
			"unexpected panic:",
			fmt.Sprintf("  recovered: %T(%v)", pm.got, opt.ValueAsString(pm.got, opts...)),
		}
	case opt.IsSet(opts, opt.ToNotMatch(true)):
		// when ToNotMatch is set, we must have recovered from an expected panic that NOT have occurred
		return []string{
			fmt.Sprintf("expected: panic with %T(%v): should not have occurred", pm.expected, opt.ValueAsString(pm.expected, opts...)),
		}
	default:
		// otherwise we were expecting a specific value to be recovered from a panic but
		// we got something else instead
		return []string{
			"unexpected panic:",
			fmt.Sprintf("  expected : %T(%v)", pm.expected, opt.ValueAsString(pm.expected, opts...)),
			fmt.Sprintf("  recovered: %T(%v)", pm.got, opt.ValueAsString(pm.got, opts...)),
		}
	}
}
