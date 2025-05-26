package bytes

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/blugnu/test/opt"
)

type EqualMatcher[T ~byte] struct {
	Expected []T
}

func (bm *EqualMatcher[T]) getDiffs(want, got []T) []int {
	rlen := len(want)
	if len(got) < rlen {
		rlen = len(got)
	}
	result := make([]int, 0, rlen+1)

	for iw, want := range want {
		if iw >= len(got) {
			return append(result, iw)
		}

		if want != got[iw] {
			result = append(result, iw)
		}
	}
	if len(got) > len(want) {
		return append(result, len(want))
	}
	return result
}

// Match returns true if the got value is equal to the wanted value.
func (bm *EqualMatcher[T]) Match(got []T, _ ...any) bool {
	return reflect.DeepEqual(bm.Expected, got)
}

func (bm *EqualMatcher[T]) OnTestFailure(got []T, opts ...any) []string {
	want := bm.Expected

	if opt.IsSet(opts, opt.ToNotMatch(true)) {
		return []string{"unexpected: []byte should not be equal"}
	}

	diff, diffs, out := bm.reportInit(want, got)
	ss, we, ge, pfx, wsfx, gsfx := bm.reportCalc(want, got, diff, diffs)

	expectedBytes := strings.Builder{}
	gotBytes := strings.Builder{}
	markers := strings.Builder{}
	expectedBytes.WriteString(pfx)
	gotBytes.WriteString(pfx)
	markers.WriteString(strings.Repeat(" ", len(pfx)))
	for i := ss; i < max(we, ge); i++ {
		if i < len(want) {
			diffs = bm.reportExpectedByte(i, ss, want, got, &expectedBytes, &markers, diffs)
		}
		if i < len(got) {
			diffs = bm.reportGotByte(i, ss, want, got, &gotBytes, &markers, diffs)
		}
	}
	expectedBytes.WriteString(wsfx)
	gotBytes.WriteString(gsfx)
	if len(want) == 0 {
		expectedBytes.WriteString("<empty>")
	}
	if len(got) == 0 {
		gotBytes.WriteString("<empty>")
	}

	out = append(out, "expected: "+expectedBytes.String())
	out = append(out, "        | "+markers.String())
	out = append(out, "got     : "+gotBytes.String())
	return out
}

func (bm EqualMatcher[T]) reportInit(want, got []T) (int, []int, []string) {
	diffs := bm.getDiffs(want, got)

	diff := -1
	out := make([]string, 0, 4)
	out = append(out, "bytes not equal:")
	if len(want) != len(got) {
		out = append(out, "  different lengths: expected "+strconv.Itoa(len(want))+", got "+strconv.Itoa(len(got)))

		// if there is a different in length we will show the last
		// diff to highlight the overrun/underrun
		diff = len(diffs) - 1
	}
	if (diff == -1 && len(diffs) > 0) || len(diffs) > 1 {
		out = append(out, "  differences at: "+strings.Join(strings.Fields(fmt.Sprint(diffs)), ", "))

		// if we aren't showing an overrun/underrun diff, we will show the first
		if diff == -1 {
			diff = 0
		}
	}
	return diff, diffs, out
}

func (bm EqualMatcher[T]) reportCalc(want, got []T, diff int, diffs []int) (int, int, int, string, string, string) {
	fd := diffs[diff]
	ss := max(0, min(fd-2, min(len(want)-2, len(got)-2)))

	we := min(fd+3, len(want))
	wsfx := "..."
	if we == len(want) {
		wsfx = ""
	}

	ge := min(fd+3, len(got))
	gsfx := "..."
	if ge == len(got) {
		gsfx = ""
	}

	pfx := "..."
	if ss == 0 {
		pfx = ""
	}
	return ss, we, ge, pfx, wsfx, gsfx
}

func (bm *EqualMatcher[T]) reportExpectedByte(i, ss int, want, got []T, expectedBytes, markers *strings.Builder, diffs []int) []int {
	if i > ss {
		expectedBytes.WriteString(" ")
		markers.WriteString(" ")
	}
	expectedBytes.WriteString(fmt.Sprintf("%02x", want[i]))
	switch {
	case i < len(got):
		if want[i] != got[i] {
			markers.WriteString("**")
			diffs = remove(diffs, i)
		} else {
			markers.WriteString("  ")
		}
	case i == len(got):
		markers.WriteString("--")
		diffs = remove(diffs, i)
	case i >= len(got):
		markers.WriteString("--")
	}
	return diffs
}

func (bm *EqualMatcher[T]) reportGotByte(i, ss int, want, got []T, gotBytes, markers *strings.Builder, diffs []int) []int {
	if i > ss {
		gotBytes.WriteString(" ")
		if i >= len(want) {
			markers.WriteString(" ")
		}
	}
	gotBytes.WriteString(fmt.Sprintf("%02x", got[i]))
	switch {
	case i == len(want):
		markers.WriteString("++")
		diffs = remove(diffs, i)
	case i > len(want):
		markers.WriteString("++")
	}
	return diffs
}
