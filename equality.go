package test

// this file contains the implementation of the test helper functions
// (test.Equal, test.IsNil etc)
//
// the implementation of the test.Helper test factory may be found in
// the file helper.go

import (
	"fmt"
)

// identifies the type of equality test to perform.
type Equality int

const (
	customEquality  Equality = iota // equality is determined using a custom comparison function
	ShallowEquality                 // equality is determined using the == operator
	DeepEquality                    // equality is determined using reflect.DeepEqual
)

// returns the name of an Equality value.
func (e Equality) String() string {
	s := map[Equality]string{
		customEquality:  "func(got, wanted) bool {...}",
		ShallowEquality: "test.ShallowEquality",
		DeepEquality:    "test.DeepEquality",
	}
	if s, ok := s[e]; ok {
		return s
	}
	return fmt.Sprintf("test.Equality(%d)", int(e))
}
