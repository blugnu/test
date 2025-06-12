package test

import "github.com/blugnu/test/internal/testframe"

// NilFrame returns a sentinel implementation of test.TestingT.  It is used to
// simulate a nil test frame in tests, allowing for testing of behavior when a
// test helper or matcher is used without a valid test frame.
//
// The NilFrame captures the current test frame in order to ensure correct
// cleanup; if there is no valid test frame, the NilFrame function itself
// will panic with testframe.ErrNoTestFrame.
//
// # Notes On Usage
//
// The NilFrame() sentinel mechanism is primarily intended for internal use,
// to verify the behaviour of the test package itself.  It provides a mechanism
// to differentiate between an invalid call to With(nil) versus a deliberate
// intention to introduce a (simulated) nil frame to verify the behaviour of the
// package under those conditions.
//
// It should not be necessary to use this in your own tests.
//
// There are some important considerations to keep in mind if you do choose to
// test with a nil test frame:
//
// Pushing a test.NilFrame() onto the stack obscurs any current test frame
// which is required for any tests to be evaluated.  This means that expectations
// must be established *before* pushing the nil frame and evaluation of those
// expectations deferred, to execute after the code under test has completed.
//
// For example, the internal usage is primarily to ensure that panics are
// raised by the test package when a nil test frame is encountered:
//
//	func TestSomething(t *testing.T) {
//	   test.With(t)
//
//	   Run("some subtest", func() {
//	      defer Expect(Panic(test.ErrNoTestFrame)).DidOccur()
//
//	      test.With(test.NilFrame())
//
//	      // this should panic, meeting the expectation established
//	      // and deferred (above)
//	      T()
//	   })
//	}
func NilFrame() testframe.Nil {
	if tf, ok := T().(testframe.Cleanup); ok {
		return testframe.Nil{T: tf}
	}
	panic(ErrNoTestFrame)
}
