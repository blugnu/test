package test

import (
	"testing"
)

func TestBytes(t *testing.T) {
	// ARRANGE
	a := []byte{0x01, 0x02, 0xf0}
	b := []byte{0x01, 0x02, 0xff}
	l := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}

	type args struct {
		wanted []byte
		got    []byte
		format BytesFormat
	}
	type result struct {
		outcome HelperResult
		output  []string
	}
	testcases := []struct {
		name string
		args
		result
	}{
		{name: "wanted == got", args: args{wanted: a, got: a}, result: result{outcome: ShouldPass}},
		{name: "no format", args: args{wanted: a, got: b}, result: result{outcome: ShouldFail, output: []string{"wanted: [01 02 f0]", "got   : [01 02 ff]"}}},
		{name: "HexBytes", args: args{a, b, BytesHex}, result: result{outcome: ShouldFail, output: []string{"wanted: [01 02 f0]", "got   : [01 02 ff]"}}},
		{name: "BinBytes", args: args{a, b, BytesBinary}, result: result{outcome: ShouldFail, output: []string{"wanted: [00000001 00000010 11110000]", "got   : [00000001 00000010 11111111]"}}},
		{name: "len(wanted) > 20", args: args{l, b, BytesHex}, result: result{outcome: ShouldFail, output: []string{"wanted: [01 02 03 ... 13 14 15] len == 21", "got   : [01 02 ff]"}}},
		{name: "len(wanted) > 20", args: args{a, l, BytesHex}, result: result{outcome: ShouldFail, output: []string{"wanted: [01 02 f0]", "got   : [01 02 03 ... 13 14 15] len == 21"}}},
		{name: "len(both) > 20", args: args{append(l, 0xff), append(l, 0x0), BytesHex}, result: result{outcome: ShouldFail, output: []string{"wanted: [01 02 03 ... 14 15 ff] len == 22", "got   : [01 02 03 ... 14 15 00] len == 22"}}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, _ := Helper(t, func(st *testing.T) {
				if tc.args.format == bytesFormatNotSet {
					Bytes(st, tc.args.wanted, tc.args.got)
				} else {
					Bytes(st, tc.args.wanted, tc.args.got, tc.args.format)
				}
			}, tc.result.outcome)

			// ASSERT
			stdout.Contains(t, tc.result.output)
		})
	}
}
