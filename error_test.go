package test

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	// ARRANGE
	g := errors.New("got error")
	wrappedg := fmt.Errorf("wrapped: %w", g)
	oe := errors.New("some other error")

	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// expected to pass
		{scenario: "got == wanted", act: func(t T) {
			Error(t, g).Is(g)
		},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scenario: "got wrapped wanted", act: func(t T) {
			Error(t, wrappedg).Is(g)
		},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},
		{scenario: "got nil , want nil", act: func(t T) {
			Error(t, nil).Is(nil)
		},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			}},

		// expected to fail
		{scenario: "got err, want nil",
			act: func(t T) {
				Error(t, errors.New("unexpected error")).Is(nil)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("got_err,_want_nil")
				test.Report.Contains([]string{
					currentFilename(),
					"unexpected error: unexpected error",
				})
			},
		},
		{scenario: "Error(nil).Is(non-nil)",
			act: func(t T) {
				Error(t, nil).Is(errors.New("desired error"))
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Error(nil).Is(non-nil)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted error: desired error",
					"got         : <nil>",
				})
			},
		},
		{scenario: "Error(err).Is(other)",
			act: func(t T) {
				Error(t, g).Is(oe)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Error(err).Is(other)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted error: some other error",
					"got         : got error",
				})
			},
		},
		{scenario: "Error(err,ErrorDefault).Is(other)",
			act: func(t T) {
				Error(t, g, ErrorDefault).Is(oe)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Error(err,ErrorDefault).Is(other)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted error: some other error",
					"got         : got error",
				})
			},
		},
		{scenario: "Error(err,ErrorString).Is(other)",
			act: func(t T) {
				Error(t, g, ErrorString).Is(oe)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Error(err,ErrorString).Is(other)")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted error: some other error",
					"got         : got error",
				})
			},
		},
		{scenario: "got != wanted, ErrorDecl",
			act: func(t T) {
				Error(t, g, ErrorDecl).Is(oe)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains([]string{
					`wanted error: &errors.errorString{s:"some other error"}`,
					`got         : &errors.errorString{s:"got error"}`,
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

func TestIsError(t *testing.T) {
	// ARRANGE
	var result bool
	var err error
	var thisTestFile = currentFilename()

	t.Run("IsError(nil)", func(t *testing.T) {
		// ACT
		sut := Helper(t, func(t T) {
			err, result = IsError(t, nil)
		})

		// ASSERT
		IsFalse(t, result)
		IsNil(t, err)

		sut.DidFail()
		sut.Report.Contains("IsError(nil)")
		sut.Report.Contains([]string{
			thisTestFile,
			"wanted: error",
			"got   : nil",
		})
	})

	t.Run("IsError(error)", func(t *testing.T) {
		// ACT
		sut := Helper(t, func(t T) {
			err, result = IsError(t, errors.New("some error"))
		})

		// ASSERT
		IsTrue(t, result)
		IsNotNil(t, err)

		sut.DidPass()
		sut.Report.IsEmpty()
	})

	t.Run("IsError(<not an error>)", func(t *testing.T) {
		// ACT
		sut := Helper(t, func(t T) {
			err, result = IsError(t, 42)
		})

		// ASSERT
		IsFalse(t, result)
		IsNil(t, err)

		sut.DidFail()
		sut.Report.Contains([]string{
			thisTestFile,
			"wanted: error",
			"got   : 42 (int)",
		})
	})
}
