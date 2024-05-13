package test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// provides methods for testing []string values.
type StringsTest struct {
	testable[[]string]
}

// creates a testable []string value.  Options of the following types are accepted:
//
//	string              // a name for the test; if not specified, "strings" is used
func Strings(t *testing.T, got any, opts ...any) StringsTest {
	n := "strings"
	checkOptTypes(t, optTypes(n), opts...)
	getOpt(&n, opts...)

	var sut []string
	switch got := got.(type) {
	case []string:
		sut = got
	case string:
		sut = []string{got}
	case []byte:
		sut = strings.Split(string(got), "\n")
	default:
		panic(ErrInvalidArgument)
	}

	return StringsTest{newTestable(t, sut, n)}
}

// fails the test if the []string does not contain the wanted string
// or strings.
//
// If a []string is provided then the strings must appear in the
// []string contiguously in the same order as they are provided.  Any
// additional strings preceding or following the wanted []string are
// ignored.  If any elements in a wanted []string are empty, the
// corresponding string in the []string being tested must also be empty
// or consist entirely of whitespace.
//
// Whether specifying a string or []string, leading and trailing
// whitespace is ignored.
//
// An attempt to test either for nil or an empty slice will fail with
// an error identifying the test as invalid.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		var err error
//
//		// ACT
//		stdout, _ :=test.CaptureOutput(t, func (t *testing.T) {
//		   err = doSomething()
//		})
//
//		// ASSERT
//		test.UnexpectedError(t, err)
//		stdout.Contains("expected output")
//	  }
func (c StringsTest) Contains(want any) {
	c.Helper()
	c.Run("contains", func(t *testing.T) {
		t.Helper()

		switch want := want.(type) {
		case nil:
			t.Errorf("\nContains(nil) is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")

		case []string:
			switch len(want) {
			case 0:
				t.Errorf("\nContains(<empty slice>) is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
			default:
				c.containsStrings(t, want)
			}

		case string:
			if len(want) == 0 {
				t.Errorf("\nContains(\"\") is invalid: Contains() accepts string or []string: did you mean IsEmpty()?")
				return
			}
			c.containsString(t, want)

		default:
			t.Errorf("\nContains(%T) is invalid: Contains() accepts string or []string", want)
		}
	})
}

// fails the test if the []string does not contain a string that
// matches the specified regular expression.
func (st StringsTest) ContainsMatch(re string) {
	st.Helper()
	st.Run("contains match", func(t *testing.T) {
		t.Helper()

		// compile the reg ex
		cre, err := regexp.Compile(re)
		if err != nil {
			t.Errorf("\ninvalid test: ContainsMatch(%s): %s", re, err)
			return
		}

		for _, s := range st.got {
			if cre.FindIndex([]byte(s)) != nil {
				return
			}
		}
		t.Errorf("\nwanted: contains match: %s\ngot   : %s", re, st.format(st.got))
	})
}

// fails the test if the []string contains the specified string.
//
// The specified string must not be empty or consist entirely of
// whitespace and is deemed to exist in the []string being tested
// if it appears in any element of the []string as a string or
// substring.
func (st StringsTest) DoesNotContain(want string) {
	st.Helper()
	st.Run("does not contain", func(t *testing.T) {
		t.Helper()

		if len(strings.TrimSpace(want)) == 0 {
			t.Errorf("\nDoesNotContain() invalid test: specified string is empty or consists entirely of whitespace")
			return
		}

		for _, s := range st.got {
			if strings.Contains(s, want) {
				t.Errorf("\nwanted: does not contain: %q\ngot   : %s", want, st.format(st.got))
				return
			}
		}
	})
}

// fails if the []string being tested is not exactly equal to some
// specified []string; they slices must have the same number of
// elements in the same order with the same content.
func (st StringsTest) Equals(want []string) {
	st.Helper()
	st.Run("equals", func(t *testing.T) {
		t.Helper()

		if slicesEqual(want, st.got, nil) {
			return
		}

		diff := []string{}
		if len(want) != len(st.got) {
			diff = append(diff, fmt.Sprintf("wanted: %d elements\ngot   : %d", len(want), len(st.got)))
		}

		ne := len(want)
		if len(st.got) < ne {
			ne = len(st.got)
		}
		ne--

		dg := len(strconv.Itoa(ne))
		ef := fmt.Sprintf("[%%%dd]", dg)
		for i := 0; i <= ne; i++ {
			if want[i] != st.got[i] {
				diff = append(diff, fmt.Sprintf(ef+" wanted: %q\n%sgot: %q", i, want[i], strings.Repeat(" ", 6+dg), st.got[i]))
			}
		}

		got := []string{"got:"}
		for i, s := range st.got {
			got = append(got, fmt.Sprintf("  "+ef+": %q", i, s))
		}

		got = append(got, "--")
		got = append(got, diff...)
		st.errorf(t, strings.Join(got, "\n"))
	})
}

// fails if the []string being tested is not empty.
//
// Note that this is not the same as testing for nil.
//
// Example:
//
//	func TestSomething(t *testing.T) {
//	  // ARRANGE
//	  var err error
//
//	  // ACT
//	  stdout, _ :=test.CaptureOutput(t, func (t *testing.T) {
//	     err = doSomething()
//	  })
//
//	  // ASSERT
//	  test.UnexpectedError(t, err)
//	  stdout.IsEmpty()
//	}
func (st StringsTest) IsEmpty() {
	if len(st.got) == 0 {
		return
	}
	st.Helper()
	st.Run("is empty", func(t *testing.T) {
		t.Helper()
		t.Errorf("\nwanted: <empty slice>\ngot   : %s", st.format(st.got))
	})
}

// fails if the []string being tested is not nil.
//
// Note that this is not the same as testing for an empty slice.
//
// Example:
//
//	func TestSomething(t *testing.T) {
//	  // ARRANGE
//	  var err error
//
//	  // ACT
//	  stdout, _ :=test.CaptureOutput(t, func (t *testing.T) {
//	     err = doSomething()
//	  })
//
//	  // ASSERT
//	  test.UnexpectedError(t, err)
//	  stdout.IsNil()
//	}
func (st StringsTest) IsNil() {
	st.Helper()
	st.Run("is nil", func(t *testing.T) {
		t.Helper()
		t.Errorf("\nwanted: nil\ngot   : %s", st.format(st.got))
	})
}

// returns a new StringsTest with all strings in the []string trimmed
// of leading and trailing whitespace.
func (st StringsTest) Trimmed() StringsTest {
	tg := []string{}
	for _, s := range st.got {
		tg = append(tg, strings.TrimSpace(s))
	}
	return Strings(st.T, tg, fmt.Sprintf("%s/trimmed", st.name))
}

// containsString is a function used to test if the []string
// contains the wanted string.  It is not exported, being
// called from Contains when the wanted value is a string or
// a slice of strings with a single entry.
func (st StringsTest) containsString(t *testing.T, want string) {
	t.Helper()

	for _, s := range st.got {
		if strings.Contains(s, want) {
			return
		}
	}
	t.Errorf("\nwanted: %q\ngot   : %s", want, st.format(st.got))
}

// containsStrings is a function used to test if the []string
// contains the wanted strings.  It is not exported, being
// called from Contains when the wanted value is a []string with
// at least 2 entries.
//
// The wanted strings must appear in the []string contiguously
// in the same order as they are provided.  If any of the wanted
// strings are empty the corresponding string in the []string
// must also be empty or consist entirely of whitespace.
func (st StringsTest) containsStrings(t *testing.T, want []string) {
	t.Helper()

	match := func(s, w string) bool {
		if w == "" {
			return strings.TrimSpace(s) == ""
		}
		return strings.Contains(s, w)
	}

	ix := []int{}
	for i, s := range st.got {
		if match(s, want[0]) {
			ix = append(ix, i)
		}
	}

	for _, is := range ix {
		found := is+len(want) <= len(st.got)
		for iw := 1; found && iw < len(want); iw++ {
			found = match(st.got[is+iw], want[iw])
		}
		if found {
			return
		}
	}

	t.Run("contains strings", func(t *testing.T) {
		t.Helper()
		t.Errorf("\nwanted: %s\ngot: %s", st.format(want), st.format(st.got))
	})
}

// formats a specified slice of strings for display in a test
// failure report.
func (st StringsTest) format(v []string) string {
	if len(v) == 0 {
		return "<empty slice>"
	}

	s := "["
	for _, v := range v {
		s += fmt.Sprintf("\n  %s", v)
	}
	s += "\n]"
	return s
}
