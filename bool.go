package test

// implements the Matcher interface for testing expected bool values.
type BooleanMatcher struct {
	Expecting[bool]
}

func (m BooleanMatcher) OneLineError() {
	/* implements OneLineExpected for test failure reporting */
}

// Match returns true if the actual value matches the expected value.
func (m BooleanMatcher) Match(actual bool, _ ...any) bool {
	return m.value == actual
}

// BeFalse returns a matcher that will fail if the actual value is not false.
func BeFalse() BooleanMatcher {
	return BooleanMatcher{Expecting[bool]{false}}
}

// BeTrue returns a matcher that will fail if the actual value is not true.
func BeTrue() BooleanMatcher {
	return BooleanMatcher{Expecting[bool]{true}}
}

// ExpectFalse fails a test if a specified bool is not false.  An optional
// name may be specified to be included in the test report in the event of
// failure.
//
// Example:
//
//	ExpectFalse(got)
//
// This function is a convenience for these alternatives:
//
//	Expect(got).To(Equal(false))
//	Expect(got).To(BeFalse())
func ExpectFalse[T ~bool](got T, opts ...any) {
	GetT().Helper()
	Expect(bool(got), opts...).To(BeFalse())
}

// ExpectTrue fails a test if a specified bool is not true.  An optional
// name may be specified to be included in the test report in the event of
// failure.
//
// Example:
//
//	ExpectTrue(got)
//
// This function is a convenience for these alternatives:
//
//	Expect(got).To(Equal(true))
//	Expect(got).To(BeTrue())
func ExpectTrue[T ~bool](got T, opts ...any) {
	GetT().Helper()
	Expect(bool(got), opts...).To(BeTrue())
}
