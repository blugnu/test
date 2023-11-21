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
	formatNotSet  Format = ""     // format is undefined
	FormatBin     Format = "%.8b" // formats as 8-bit binary, e.g.: [00000001 00000010 11110000]
	FormatHex     Format = "%x"   // formats as hexadecimal, e.g.: 0102f0
	FormatDecl    Format = "%#v"  // formats using the go declaration for its value, e.g.: []byte{0x01 0x02 0xf0}
	FormatDefault Format = "%v"   // formats using the default representation for its value, e.g.: [1 2 240]
)

// format is a function used to format the output of a value
// according to a specified format.
func format(v any, f Format) string {
	if f == formatNotSet {
		f = FormatDefault
	}
	return fmt.Sprintf(string(f), v)
}
