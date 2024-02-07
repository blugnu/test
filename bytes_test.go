package test

import (
	"testing"
)

func TestBytes(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		act      func(T)
		assert   func(HelperTest)
	}{
		// these tests should pass
		{scenario: "BytesEqual(got === wanted)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				BytesEqual(t, g, g)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "BytesEqual(got == wanted)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				w := []byte{0x01, 0x02, 0xf0}
				BytesEqual(t, g, w)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Bytes().Equals() (got === wanted)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				Bytes(t, g).Equals(g)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},
		{scenario: "Bytes().Equals() (got == wanted)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				w := []byte{0x01, 0x02, 0xf0}
				Bytes(t, g).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidPass()
				test.Report.IsEmpty()
			},
		},

		// these should fail
		{scenario: "BytesEqual(!=)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				w := []byte{0x01, 0x02, 0xff}
				BytesEqual(t, g, w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [01 02 ff]",
					"got   : [01 02 f0]",
				})
			},
		},
		{scenario: "BytesEqual(!=,BytesBinary)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0xf0}
				w := []byte{0x01, 0x02, 0xff}
				BytesEqual(t, g, w, BytesBinary)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=,BytesBinary)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [00000001 00000010 11111111]",
					"got   : [00000001 00000010 11110000]",
				})
			},
		},
		{scenario: "BytesEqual(!=) (BytesString)",
			act: func(t T) {
				g := []byte("abc")
				w := []byte("def")
				BytesEqual(t, g, w, BytesString)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=)_(BytesString)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: \"def\"",
					"got   : \"abc\"",
				})
			},
		},
		{scenario: "BytesEqual(!=) (len > 20)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0xff}
				w := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
				BytesEqual(t, g, w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=)_(len_>_20)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [01 02 03 ... 13 14 15] len == 21",
					"got   : [01 02 03 ... 13 14 ff] len == 21",
				})
			},
		},
		{scenario: "BytesEqual(!=) (wanted len > 20)",
			act: func(t T) {
				g := []byte(nil)
				w := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
				BytesEqual(t, g, w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=)_(wanted_len_>_20)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: [01 02 03 ... 13 14 15] len == 21",
					"got   : []",
				})
			},
		},
		{scenario: "BytesEqual(!=) (got len > 20)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0xff}
				w := []byte(nil)
				BytesEqual(t, g, w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=)_(got_len_>_20)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []",
					"got   : [01 02 03 ... 13 14 ff] len == 21",
				})
			},
		},
		{scenario: "BytesEqual(!=,10) (got len > 20)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0xff}
				w := []byte(nil)
				BytesEqual(t, g, w, 10)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("BytesEqual(!=,10)_(got_len_>_20)/bytes")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []",
					"got   : [01 02 03 ... 13 14 ff] len == 21",
				})
			},
		},
		{scenario: "Bytes(g,10).Equals(w) (got len > 20)",
			act: func(t T) {
				g := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0xff}
				w := []byte(nil)
				Bytes(t, g, 10).Equals(w)
			},
			assert: func(test HelperTest) {
				test.DidFail()
				test.Report.Contains("Bytes(g,10).Equals(w)_(got_len_>_20)/bytes/equals")
				test.Report.Contains([]string{
					currentFilename(),
					"wanted: []",
					"got   : [01 02 03 ... 13 14 ff] len == 21",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.assert(Helper(t, tc.act))
		})
	}
}
