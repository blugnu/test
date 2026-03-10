package testframe

import (
	"reflect"
	"testing"
)

// IsParallel checks if the current test frame is running in parallel.
//
// The underlying type of the current test frame must be a *testing.T;
// any other implementation will return false.
//
// For *testing.T values, reflection is used to examine the internal
// state to determine if it is running in parallel.
func IsParallel() bool {
	t, ok := Peek[*testing.T]()
	if !ok {
		// tests cannot be parallel if the TestingT is not a *testing.T
		return false
	}

	c := reflect.Indirect(reflect.ValueOf(t)).FieldByName("common")
	ip := reflect.Indirect(c).FieldByName("isParallel")

	return ip.Bool()
}
