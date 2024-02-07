package test

import (
	"bytes"
	"fmt"
	"testing"
)

// provides methods for testing a []byte value.
type BytesTest struct {
	testable[[]byte]
}

// returns a new BytesTest with methods for testing a specified
// []byte using a specified testing.T.
//
// Additional options of the following types may be specified:
//
//   - string : a name for the value being tested; if not specified "bytes" is used
//
//   - BytesFormat : a format verb for the value being tested; if not specified
//     BytesHex is used.  This option is ignored if a format function is also
//     specified.  If BytesString is specified, slices are reported as quoted strings.
//
//   - int : the maximum number of bytes to display in a test failure report when
//     formatting values using a specified BytesFormat other than BytesString.  If
//     not specified or a value less than 20 is specified, 20 is used. If a slice
//     has more elements than the specified maximum then only the first and last 3
//     bytes are output together with the total length of the slice. This option is
//     ignored if BytesString is specified (the entire slice is output as a quoted
//     string) or if a function is specified
//
//   - func([]byte) string : a function that returns a string representation of the
//     slice type being tested.  If not specified, values are formatted using the
//     BytesFormat and int options.
//
// If more than one option of any of the above types is specified then only the first
// is applied; additional values of that option type are ignored.
func Bytes(t *testing.T, got []byte, opts ...any) BytesTest {
	t.Helper()

	n := "bytes"
	f := BytesHex
	m := 20
	ffn := *new(func([]byte) string)
	checkOptTypes(t, optTypes(n, f, m, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	getOpt(&m, opts...)
	if m < 20 {
		m = 20
	}

	getOpt(&ffn, append(opts, func(v []byte) string {
		if f == BytesString {
			return fmt.Sprintf("%q", v)
		}

		if len(v) > m {
			return fmt.Sprintf(fmt.Sprintf("[%% %[1]s ... %% %[1]s]", f)+" len == %d", v[:3], v[len(v)-3:], len(v))
		} else if f != BytesBinary {
			return format(v, Format(fmt.Sprintf("[%% %s]", f)))
		}

		return format(v, Format(fmt.Sprintf("%%%s", f)))
	})...)

	return BytesTest{testable: newTestable(t, got, n, Format(f), ffn)}
}

// fails the test if the []byte being tested does not equal the wanted
// value.
//
// If the length of either the wanted slice or the value in the BytesTest
// (got) exceeds the maximum report bytes for the test, the output in a
// test failure report will be truncated to the first 3 and last 3 bytes,
// with the length of each slice reported.
//
// Example:
//
//	test.Bytes(t, result.Bytes(), "result").Equals(expected)
func (bt BytesTest) Equals(wanted []byte) {
	if bytes.Equal(wanted, bt.got) {
		return
	}

	bt.Helper()
	bt.Run("equals", func(t *testing.T) {
		t.Helper()
		bt.fail(t, wanted)
	})
}

// fails the test if got is not equal to wanted.  An optional name and
// BytesFormat may be specified; if no name is specified then a name of
// "bytes" is used.  If no BytesFormat is specified, BytesHex is used.
//
// Example:
//
//	test.BytesEqual(t, result.Bytes(), Equals(expected), "result buffer", BytesString)
//
// BytesEqual(t, got, wanted, name, format) is a convenience short-hand for:
//
//	test.Bytes(t, name, got, format).Equals(wanted)
//
// The following tests are exactly equivalent:
//
//	test.Bytes(t, result.Bytes(), "result", BytesString).Equals(expected)
//	test.BytesEqual(t, result.Bytes(), expected, "result", BytesString)
//
// and:
//
//	test.Bytes(t, result.Bytes()).Equals(expected)
//	test.BytesEqual(t, result.Bytes(), expected)
func BytesEqual(t *testing.T, got, wanted []byte, opts ...any) {
	t.Helper()

	if bytes.Equal(wanted, got) {
		return
	}

	// to produce the test report we will initialise a BytesTest which
	// will provide defaults for the max length and format function if
	// these options are not specified for this test.

	n := "bytes"
	f := BytesHex
	var m int
	var ffn func([]byte) string
	checkOptTypes(t, optTypes(n, f, m, ffn), opts...)
	getOpt(&n, opts...)
	getOpt(&f, opts...)
	opts = append(opts, n, f)

	t.Run(n, func(t *testing.T) {
		t.Helper()
		Bytes(t, got, opts...).fail(t, wanted)
	})
}
