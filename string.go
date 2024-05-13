package test

import (
	"fmt"
	"strings"
	"testing"
)

// provides methods for testing []string values.
type StringTest struct {
	testable[string]
}

// creates a testable string value.  The got value being tested may be of any type
// and will be interpreted as follows:
//
//	string       // a string to be tested
//
//	[]byte       // a byte slice containing a string to be tested
//
//	any          // any other type will be converted using fmt.Sprintf(%v, got)
//	             // to the string to be tested
//
// In addition, options of the following types are accepted:
//
//	string     // a name for the test; if not specified, "string" is used
func String(t *testing.T, got any, opts ...any) StringTest {
	n := "string"
	checkOptTypes(t, optTypes(n), opts...)
	getOpt(&n, opts...)

	var sut string
	switch got := got.(type) {
	case string:
		sut = got
	case []byte:
		sut = string(got)
	default:
		sut = fmt.Sprintf("%v", got)
	}

	return StringTest{newTestable(t, sut, n)}
}

// Contains asserts that the string being tested contains the specified substring.
func (c StringTest) Contains(want string) {
	c.Helper()
	c.Run("contains", func(t *testing.T) {
		t.Helper()

		if len(want) == 0 {
			t.Errorf("\nContains(<empty string>) is invalid: did you mean IsEmpty()?")
		}
		if !strings.Contains(c.got, want) {
			t.Errorf("\nwanted: string containing: %q\ngot   : %q", want, c.got)
		}
	})
}

// DoesNotContain asserts that the string being tested does not contain the
// specified substring.
func (c StringTest) DoesNotContain(substring string) {
	c.Helper()
	c.Run("does_not_contain", func(t *testing.T) {
		t.Helper()

		if len(substring) == 0 {
			t.Errorf("\nDoesNotContain(<empty string>) is invalid test: did you mean IsNotEmpty()?")
		}
		if x := strings.Index(c.got, substring); x > -1 {
			t.Errorf("\nfound: %q\nat   : %d\ngot  : %q\n        %s%s",
				substring,
				x,
				c.got,
				strings.Repeat(" ", x),
				strings.Repeat("^", len(substring)),
			)
		}
	})
}

// Equals asserts that the string being tested is equal to the specified string.
func (c StringTest) Equals(want string) {
	c.Helper()
	c.Run("equals", func(t *testing.T) {
		t.Helper()

		if c.got != want {
			t.Errorf("\nwanted: %q\ngot   : %q", want, c.got)
		}
	})
}
