package test

import (
	"github.com/blugnu/test/matchers/contexts"
)

// HaveContextKey returns a matcher that checks if a context contains a specific key.
// The key type must be comparable. The matcher will fail if the key is not
// present in the context.
//
// # Supported Options
//
//	opt.QuotedStrings(bool)    // determines whether string values are quoted in test
//	                           // failure report (quoted by default); the option has
//	                           // has no effect if the key is not a string type
//
//	opt.FailureReport(func)    // a function returning a custom failure report
//	                           // in the event that the test fails
func HaveContextKey[K comparable](k K) *contexts.KeyMatcher[K] {
	return &contexts.KeyMatcher[K]{Expected: k}
}

// HaveContextValue returns a matcher that checks if a context contains a specific
// key-value pair.
//
// The key type (K) must be comparable. The matcher will fail if the key is not
// present in the context or if the value does not match the expected value.
// The value type (V) can be any type.
//
// The matcher uses reflect.DeepEqual to compare the expected value with any value
// in the context for the specified key; this may be overridden by supplying a
// custom comparison function in the options.
//
// # Supported Options
//
//	func(V, V) bool            // a custom comparison function to compare values
//	                           // (overriding the use of reflect.DeepEqual)
//
//	opt.QuotedStrings(bool)    // determines whether string keys or values are quoted
//	                           // in the test failure report (quoted by default);
//	                           // the option has no effect for keys or values that
//	                           // are not string type
//
//	opt.FailureReport(func)    // a function returning a custom failure report
//	                           // in the event that the test fails
func HaveContextValue[K comparable, V any](k K, v V) *contexts.ValueMatcher[K, V] {
	return &contexts.ValueMatcher[K, V]{
		Key:      k,
		Expected: v,
	}
}
