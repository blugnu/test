package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/blugnu/capture"
)

// CapturedOutput is a type used to capture the output of a test or
// other function
type CapturedOutput struct {
	src string
	s   []string
}

// Contains fails the test if the captured output does not contain
// the wanted string or strings.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		var err error
//
//		// ACT
//		stdout, _ := test.CaptureOutput(t, func (t *testing.T) {
//		   err = doSomething()
//		})
//
//		// ASSERT
//		test.UnexpectedError(t, err)
//		stdout.Contains(t, "expected output")
//	  }
func (c CapturedOutput) Contains(t *testing.T, want any) {
	t.Helper()

	switch want := want.(type) {
	case nil:
		c.IsEmpty(t)

	case []string:
		switch len(want) {
		case 0:
			c.IsEmpty(t)
		case 1:
			c.containsString(t, want[0])

		default:
			c.containsStrings(t, want)
		}

	case string:
		c.containsString(t, want)

	default:
		panic(fmt.Errorf("%w: accepts string, []string or nil", ErrInvalidArgument))
	}
}

// containsString is a function used to test if the captured
// output contains the wanted string.  It is not exported, being
// called from Contains when the wanted value is a string or
// a slice of strings with a single entry.
func (c CapturedOutput) containsString(t *testing.T, want string) {
	t.Helper()

	switch {
	case len(want) == 0:
		c.IsEmpty(t)
		return

	default:
		for _, s := range c.s {
			if strings.Contains(s, want) {
				return
			}
		}

		t.Run("output contains", func(t *testing.T) {
			t.Helper()
			t.Errorf("\nwanted: %q\ngot   : %s", want, c.dump())
		})
	}
}

// containsStrings is a function used to test if the captured
// output contains the wanted strings.  It is not exported, being
// called from Contains when the wanted value is a []string with
// at least 2 entries.
//
// The wanted strings must appear in the captured output contiguously
// in the same order as they are provided.  If any of the wanted
// strings are empty the corresponding string in the captured output
// must also be empty or consist entirely of whitespace.
func (c CapturedOutput) containsStrings(t *testing.T, want []string) {
	match := func(s, w string) bool {
		if w == "" {
			return strings.TrimSpace(s) == ""
		}
		return strings.Contains(s, w)
	}

	ix := []int{}
	for i, s := range c.s {
		if match(s, want[0]) {
			ix = append(ix, i)
		}
	}

	for i, is := range ix {
		found := true
		for iw := 1; found && i < len(want); i++ {
			found = match(c.s[is+iw], want[iw])
		}
		if found {
			return
		}
	}

	w := "\n  " + strings.Join(want, "\n  ")

	t.Run("output contains", func(t *testing.T) {
		t.Helper()
		t.Errorf("\nwanted:%s\ngot: %s", w, c.dump())
	})
}

// dump is a helper function used to format the captured output
// for display in a test failure message.
func (c CapturedOutput) dump() string {
	if len(c.s) == 0 {
		return "(no output)"
	}
	s := "\n  " + strings.Join(c.s, "\n  ")
	return fmt.Sprintf("BEGIN CAPTURE [%s]%s\n--END CAPTURE", c.src, s)
}

// IsEmpty fails the test if the captured output is not empty.
func (c CapturedOutput) IsEmpty(t *testing.T) {
	t.Helper()
	if len(c.s) == 0 {
		return
	}
	t.Run("output is empty", func(t *testing.T) {
		t.Helper()
		t.Errorf("\nwanted: (no output)\ngot   : %s", c.dump())
	})
}

// CaptureOutput captures the stdout and stderr output of a test helper
// or other function.
//
// CaptureOutput does not perform any tests on the captured output but
// returns a CapturedOutput for both stdout and stderr.  Either of these
// values may be discarded or used to test the output, as required.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		var err error
//
//		// ACT
//		stdout, _ := test.CaptureOutput(t, func (t *testing.T) {
//		   err = doSomething()
//		})
//
//		// ASSERT
//		test.UnexpectedError(t, err)
//		stdout.Contains(t, "expected output")
//	  }
func CaptureOutput(t *testing.T, fn func(*testing.T)) (CapturedOutput, CapturedOutput) {
	stdout, stderr, _ := capture.Output(func() error {
		fn(t)
		return nil
	})
	return CapturedOutput{
			src: "stdout",
			s:   stdout,
		}, CapturedOutput{
			src: "stderr",
			s:   stderr,
		}
}
