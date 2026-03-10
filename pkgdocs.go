package test

// AnyMatcher is the interface implemented by matchers that can test
// any type of value.  It is used to apply matchers to expectations
// that are not type-specific.
//
// It is preferable to use the Matcher[T] interface for type-safe
// expectations; AnyMatcher is provided for situations where the
// compatible types for a matcher cannot be enforced at compile-time.
//
// When implementing an AnyMatcher, it is important to ensure that
// the matcher fails a test if it is not used correctly, i.e. if the
// matcher is not compatible with the type of the value being tested.
//
// An AnyMatcher must be used with the Expect().Should() matching
// function; they may also be used with Expect(got).To() where the got
// value is of type `any`, though this is not recommended.
//
// NOTE: This interface is not referenced by the test package; the
// declaration is provided for documentation purposes only.
type AnyMatcher interface {
	Match(got any, opts ...any) bool
}

// Matcher is a generic interface that describes the method set of
// a type-safe matcher for testing values of a specific type, T.
//
// A type-safe matcher is not necessarily generic.  A matcher
// that implements Match(got X, opts ...any) bool, where X is a formal,
// literal type (i.e. not generic) is also a type-safe matcher.
//
// `Matcher[T]` describes the general form of a type-safe matcher.
//
// # Generic Matchers
//
// A type-safe matcher that implements support for a specific type is
// limited to only testing values of that explicit type.
//
// A generic matcher uses type constraints to specify a set of
// compatible types for the matcher, enabling that matcher to be used
// with a variety of types.
//
// For example, the equals.Matcher[T comparable] may be used to test,
// values of any type that is comparable (supports comparison using ==).
//
// NOTE: This interface is not referenced by the test package; the
// declaration is provided for documentation purposes only.
type Matcher[T any] interface {
	Match(got T, opts ...any) bool
}
