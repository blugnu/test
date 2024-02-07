package test

import (
	"errors"
	"testing"
)

func TestExpectedPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("error")

	// ARRANGE
	testcases := []struct {
		scenario string
		arg      any
		result   *PanicTest
	}{
		{scenario: "nil argument", arg: nil, result: nil},
		{scenario: "error argument", arg: err, result: &PanicTest{err}},
		{scenario: "int argument", arg: 42, result: &PanicTest{42}},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			got := ExpectPanic(tc.arg)

			// ASSERT
			switch {
			case tc.result == nil:
				Equal(t, tc.result, got)
			case got != nil:
				t.Run("Panic with error", func(t *testing.T) {
					Equal(t, got.r, tc.result.r)
				})
			default:
				t.Errorf("\nwanted: %#v\ngot   : nil", tc.result)
			}
		})
	}
}

func TestPanic(t *testing.T) {
	// ARRANGE
	err := errors.New("error")

	testcases := []struct {
		scenario string
		sut      *PanicTest
		act      func(T)
		assert   func(HelperTest)
	}{
		{scenario: "nil receiver, no panic",
			act: func(t T) {
				defer (*PanicTest)(nil).Assert(t)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "nil receiver, panicked",
			act: func(t T) {
				defer (*PanicTest)(nil).Assert(t)
				panic(42)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("nil_receiver,_panicked")
				test.Report.Contains([]string{
					currentFilename(),
					"unexpected panic: int: 42",
				})
			},
		},
		{scenario: "expect panic, does not panic",
			act: func(t T) {
				defer ExpectPanic(errors.New("panicked")).Assert(t)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("expect_panic,_does_not_panic")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: panic: *errors.errorString: panicked",
					"got   : (did not panic)",
				})
			},
		},
		{scenario: "expect panic(err), panicked with err",
			act: func(t T) {
				defer ExpectPanic(err).Assert(t)
				panic(err)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "expect panic(42), panicked with 42",
			act: func(t T) {
				defer ExpectPanic(42).Assert(t)
				panic(42)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "expect panic(\"42\"), panicked with 42",
			act: func(t T) {
				defer ExpectPanic("42").Assert(t)
				panic(42)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("expect_panic(\"42\"),_panicked_with_42")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: panic: string: 42",
					"got   : panic: int: 42",
				})
			},
		},
		{scenario: "expect panic(err), panicked with other err",
			act: func(t T) {
				defer ExpectPanic(err).Assert(t)
				panic(errors.New("other error"))
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("expect_panic(err),_panicked_with_other_err")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: panic: *errors.errorString: error",
					"got   : panic: *errors.errorString: other error",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}
