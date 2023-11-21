package test

import (
	"bytes"
	"fmt"
	"testing"
)

// BytesFormat values are used to specify the format of
// bytes displayed in test failure reports.
//
// Supported formats are:
//
//   - BinBytes		(%.8b)
//   - DecBytes		(%d)
//   - HexBytes		(%x)
type BytesFormat string

const (
	bytesFormatNotSet BytesFormat = ""    // zero value for BytesFormat, indicates not set
	BytesBinary       BytesFormat = ".8b" // formats as 8-bit binary, e.g.: [00000001 00000010 11110000]
	BytesDecimal      BytesFormat = "d"   // formats as decimal, e.g.: 1 2 240
	BytesHex          BytesFormat = "x"   // formats as hexadecimal, e.g.: 0102f0
)

// Bytes fails the test if got is not equal to wanted.
//
// The test output in the event of a failure will be formatted
// according to the first BytesFormat argument. If no BytesFormat
// is provided, BytesHex is used.
//
// If the length of either want or got exceeds 20 bytes, the
// output in any test failure report will be truncated to the
// first 3 and last 3 bytes with the length of the slice.
//
// Example:
//
//	  func TestSomething(t *testing.T) {
//		// ARRANGE
//		a := []byte{0x01, 0x02, 0xf0}
//
//		// ACT
//		b := somethingReturningBytes()
//
//		// ASSERT
//		test.Bytes(t, a, b, BytesBin)
//	  }
func Bytes(t *testing.T, want, got []byte, opt ...BytesFormat) {
	t.Helper()

	if !bytes.Equal(want, got) {
		var d = BytesHex
		if len(opt) > 0 {
			d = opt[0]
		}

		w, g := "", ""
		if len(want) > 20 {
			w = fmt.Sprintf(fmt.Sprintf("[%% %[1]s ... %% %[1]s]", d)+" len == %d", want[:3], want[len(want)-3:], len(want))
		} else if d != BytesBinary {
			w = format(want, Format(fmt.Sprintf("[%% %s]", d)))
		} else {
			w = format(want, Format(fmt.Sprintf("%%%s", d)))
		}

		if len(got) > 20 {
			g = fmt.Sprintf("[% x ... % x] len == %d", got[:3], got[len(got)-3:], len(got))
		} else if d != BytesBinary {
			g = format(got, Format(fmt.Sprintf("[%% %s]", d)))
		} else {
			g = format(got, Format(fmt.Sprintf("%%%s", d)))
		}

		t.Errorf("\nwanted: %s\ngot   : %s", w, g)
	}
}
