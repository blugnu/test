package test

import "fmt"

// Format values are used to specify the format of
// values displayed in test failure reports.
//
// Supported formats are:
//
//   - FormatBin		(%.8b)
//   - FormatHex		(%x)
//   - FormatDecl		(%#v)
//   - FormatDefault	(%v)

type Format string

const (
	formatNone    Format = ""     // format is undefined
	FormatBin     Format = "%.8b" // formats as 8-bit binary, e.g.: [00000001 00000010 11110000]
	FormatHex     Format = "%02x" // formats as hexadecimal, e.g.: 0102f0
	FormatDecl    Format = "%#v"  // formats using the go declaration for its value, e.g.: []byte{0x01 0x02 0xf0}
	FormatString  Format = "%s"   // formats using the %s representation of the value
	FormatDefault Format = "%v"   // formats using the default representation for its value, e.g.: [1 2 240]
)

// BytesFormat values are used to specify the format of
// bytes displayed in test failure reports produced by
// test methods of a BytesTest.
//
// Supported formats are:
//
//   - BytesBinary		(%.8b)
//   - BytesDecimal		(%d)
//   - BytesHex			(%x)
//   - BytesString		(%x)
type BytesFormat string

const (
	bytesFormatNotSet BytesFormat = ""    // zero value for BytesFormat, indicates not set
	BytesBinary       BytesFormat = ".8b" // formats as 8-bit binary, e.g.: [00000001 00000010 11110000]
	BytesDecimal      BytesFormat = "d"   // formats as decimal, e.g.: 1 2 240
	BytesHex          BytesFormat = "x"   // formats as hexadecimal, e.g.: 0102f0
	BytesString       BytesFormat = "s"   // formats as string
)

// format is a function used to format the output of a value
// according to a specified format.
func format(v any, f Format) string {
	if f == formatNone {
		f = FormatDefault
	}
	return fmt.Sprintf(string(f), v)
}
