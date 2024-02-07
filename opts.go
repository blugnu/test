package test

import (
	"fmt"
	"reflect"
	"testing"
)

// checks that all options are of the types specified.  If any option
// is not of a supported type, the test is failed.  A failed test
// report will include a list of unsupported option types, similar to:
//
//	invalid test:
//	  option #2 is an unsupported type: string
//	  option #3 is an unsupported type: int
func checkOptTypes(t *testing.T, types []reflect.Type, opts ...any) {
	t.Helper()
	i := []reflect.Type{}

	for _, opt := range opts {
		ot := reflect.TypeOf(opt)
		if !sliceContains(types, ot, nil) {
			i = append(i, ot)
		}
	}
	if len(i) > 0 {
		s := ""
		for n, t := range i {
			s += fmt.Sprintf("\n  option #%d is an unsupported type: %s", n+1, t)
		}
		t.Errorf("\ninvalid test: %s", s)
	}
}

// extracts a name (string) and format (value of type T) from a list of
// options.  Only the first string or value of type T are used; any
// additional string or T values are ignored.
//
// The function returns true if an option of the desired type was
// identified, otherwise false.
func getOpt[T any](dest *T, opts ...any) bool {
	for _, opt := range opts {
		switch v := opt.(type) {
		case T:
			*dest = v
			return true
		}
	}
	return false
}

// returns a slice of reflect.Type values representing the types of the
// specified options.
func optTypes(opts ...any) []reflect.Type {
	t := []reflect.Type{}
	for _, opt := range opts {
		t = append(t, reflect.TypeOf(opt))
	}
	return t
}
