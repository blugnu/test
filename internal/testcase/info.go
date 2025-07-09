package testcase

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldValue retrieves a field value of a specified type (reflect.Kind) from
// a struct or pointer to a struct.
//
// It checks for the specified field names and returns the value if found
// using a supplied function to obtain the value from the field ref.
func FieldValue[T any](
	tc any,
	k reflect.Kind,
	result func(v reflect.Value) T,
	fields ...string,
) (T, bool) {
	// if the test case is not a struct (or pointer to a struct), return false
	ref := reflect.Indirect(reflect.ValueOf(tc))
	if ref.Kind() != reflect.Struct {
		return *new(T), false
	}

	// otherwise, check for supported fields that may contain a usable name
	val := reflect.Indirect(ref)
	for _, c := range fields {
		if f := val.FieldByName(c); f.IsValid() && f.Kind() == k {
			return result(f), true
		}
	}

	return *new(T), false
}

// NameOrDefault derives a name for a given test case.  If a name (n) is given,
// it is used without checking the test case.
//
// If no name is given and the test case is a struct with a string field named
// "Scenario", "scenario", "NameOrDefault", or "name", then this field is used as the
// name (if not empty or whitespace).
//
// If the test case is not a struct, has no name field or has a name which is
// empty, a default name in the format "testcase-NNN" where NNN is the index
// (i) of the scenario.
func NameOrDefault(tc any, n string, i int) string {
	if n = strings.TrimSpace(n); n == "" {
		extractFn := func(v reflect.Value) string { return strings.TrimSpace(v.String()) }
		fields := []string{"name", "scenario", "Name", "Scenario"}

		if s, ok := FieldValue(tc, reflect.String, extractFn, fields...); ok && s != "" {
			n = s
		}
	}

	switch {
	case n == "":
		return fmt.Sprintf("testcase-%.3d", i)
	default:
		return n
	}
}

// IsDebugging returns true if the specified test case is a struct with
// a bool debug/Debug field that is set true, otherwise it returns false.
func IsDebugging(tc any) bool {
	extractFn := func(v reflect.Value) bool { return v.Bool() }
	fields := []string{"debug", "Debug"}

	result, _ := FieldValue(tc, reflect.Bool, extractFn, fields...)

	return result
}

// IsSkipping returns true if the specified test case is a struct with
// a bool skip/Skip field that is set true, otherwise it returns false.
func IsSkipping(tc any) bool {
	extractFn := func(v reflect.Value) bool { return v.Bool() }
	fields := []string{"skip", "Skip"}

	result, _ := FieldValue(tc, reflect.Bool, extractFn, fields...)

	return result
}
