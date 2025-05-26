package test

import "github.com/blugnu/test/matchers/maps"

// KeysOfMap returns the keys of a map as a slice, provided to enable
// map keys to be tests using slice matchers, written expressively as:
//
//	Expect(KeysOfMap(someMap)).To(ContainItem(expectedKey))
//
// The order of the keys in the returned slice is not guaranteed.
func KeysOfMap[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// ValuesOfMap returns the values of a map as a slice, provided to enable
// map values to be tests using slice matchers, written expressively as:
//
//	Expect(ValuesOfMap(someMap)).To(ContainItem(expectedValue))
//
// The order of the values in the returned slice is not guaranteed.
func ValuesOfMap[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// ContainMap returns a matcher that will match if a map contains an expected
// compatible map.  The matcher will pass if the map contains all of the
// key-value pairs in the expected map.
//
// Values are compared as follows, in order of preference:
//
//  1. V.Equal(V) when V implements the Equal method
//  1. a comparison function option of the form func(V, V) bool
//  2. a comparison function option of the form func(any, any) bool
//  3. reflect.DeepEqual
//
// # Supported Options
//
//	func(V, V) bool             // a custom comparison function to compare
//	                            // map values.
//
//	func(any, any) bool         // a custom comparison function to compare
//	                            // map values.
//
//	opt.CaseSensitive(bool)     // used to indicate whether the comparison of
//	                            // string values should be case-sensitive
//	                            // (default is true).
//	                            //
//	                            // This option is also applied to values that
//	                            // are slices or arrays of a string type.
//	                            //
//	                            // The option does NOT apply to string fields
//	                            // in structs or values that are maps.
//	                            //
//	                            // NOTE: string keys are ALWAYS case-sensitive.
//
//	opt.ExactOrder(bool)        // used to indicate whether the order of
//	                            // elements in values that are slices or arrays
//	                            // is significant (default is true).
//	                            //
//	                            // NOTE: order is never significant for map keys.
//
//	opt.QuotedStrings(bool)     // determines whether string keys or values are
//	                            // quoted in any test failure report.  Default
//	                            // is true.
//	                            //
//	                            // This option has no effect for keys or values
//	                            // that are not strings.
//
//	opt.FailureReport(func)     // a function that returns a custom test
//	                            // failure report if the test fails.
//
//	opt.OnFailure(string)       // a string to output as the failure
//	                            // report if the test fails.
func ContainMap[K comparable, V any](want map[K]V) maps.ContainsMatcher[K, V] {
	T().Helper()
	switch {
	case want == nil:
		invalidTest(
			"ContainMap() called with nil map.",
			"Did you mean Expect(map).IsNil() or Expect(map).IsEmpty()?",
		)
	case len(want) == 0:
		invalidTest(
			"ContainMap() called with empty map.",
			"Did you mean Expect(map).To(EqualMap(<empty map>)) or Expect(map).IsEmpty()?",
		)
	default:
		return maps.ContainsMatcher[K, V]{Expected: want}
	}
	return maps.ContainsMatcher[K, V]{}
}

// ContainMapEntry is a convenience function to test that a map contains
// a specific key-value pair.
//
// i.e. the following are equivalent:
//
//	Expect(someMap).To(ContainMapEntry(key, value))
//	Expect(someMap).To(ContainMap(map[K]V{key: value}))
//
// Values are compared as follows, in order of preference:
//
//  1. V.Equal(V) when V implements the Equal method
//  1. a comparison function option of the form func(V, V) bool
//  2. a comparison function option of the form func(any, any) bool
//  3. reflect.DeepEqual
//
// # Supported Options
//
//	// supported options are the same as for ContainMap()
func ContainMapEntry[K comparable, V any](key K, value V) maps.ContainsMatcher[K, V] {
	return maps.ContainsMatcher[K, V]{Expected: map[K]V{key: value}}
}

// EqualMap returns a matcher that will match if the wanted map is
// equal to the actual map.
//
// To be equal, the maps must have the same number of keys and each key
// must have a value that is equal to the value of the same key in both maps.
//
// The order of the keys in the maps is not significant.
//
// Values are compared as follows, in order of preference:
//
//  1. V.Equal(V) when V implements the Equal method
//  1. a comparison function option of the form func(V, V) bool
//  2. a comparison function option of the form func(any, any) bool
//  3. reflect.DeepEqual
//
// This test is a convenience for separately testing length and content of
// two maps.  i.e. the following are equivalent:
//
//	Expect(got).To(EqualMap(want))
//
// and
//
//	Expect(len(got)).To(Equal(len(want)))
//	Expect(got).To(ContainMap(want))
//
// However, although the above tests are equivalent in terms of test outcome,
// the test failure report for the EqualMap() test is more informative than
// the two tests combined, expressing the intent that the maps be equal.
//
// # Supported Options
//
//	// supported options are the same as for ContainMap()
func EqualMap[K comparable, V any](want map[K]V) maps.EqualMatcher[K, V] {
	return maps.EqualMatcher[K, V]{Expected: want}
}
