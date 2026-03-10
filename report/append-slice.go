package report

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/blugnu/test/internal/testframe"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/test"
)

type JustifyPrefix bool

// AppendArray appends a string representation of the provided array a to the
// provided result slice of strings.
//
// Since generic constraints cannot be used to ensure that only arrays are
// passed to this function, a runtime check is performed.  If the provided
// argument is not an array, the current test is failed as invalid.
//
// # Options
//
// The following options are supported for formatting the array and any element:
//
//   - `opt.Name`        A prefix to prepend to the first line of the array representation.
//   - `opt.MaxItems`    The maximum number of items to include in the representation;
//     if the array has more items than this limit, a message indicating
//     how many additional items are not shown is appended to the result.
func AppendArray(r []string, a any, opts ...any) []string {
	if t, ok := testframe.Peek[*testing.T](); ok {
		t.Helper()
	}

	var (
		name    opt.Name
		namePad string
		hasName bool
		val     = reflect.ValueOf(a)
		kind    = val.Kind()
	)

	if kind != reflect.Array {
		test.Invalid("report.AppendArray requires an array argument")
	}

	if name, opts, hasName = opt.Extract[opt.Name](opts); hasName {
		namePad = strings.Repeat(" ", len(name))
	}

	switch {
	case val.Len() == 0:
		return append(r, string(name)+" <empty array>")

	case len(name) > 0:
		name += " "
		namePad += " "
	}

	return appendItems(r, val, name, namePad, opts...)
}

// AppendSlice appends a string representation of the provided slice s to the
// provided result slice of strings.
//
// The function is implemented as a generic function to ensure that only slices
// may be passed to it.
func AppendSlice[T any](r []string, s []T, opts ...any) []string {
	var (
		name    opt.Name
		namePad string
		hasName bool
	)

	if name, opts, hasName = opt.Extract[opt.Name](opts); hasName {
		namePad = strings.Repeat(" ", len(name))
	}

	switch {
	case s == nil:
		return append(r, string(name)+" "+Nil)

	case len(s) == 0:
		return append(r, string(name)+" <empty slice>")

	case len(name) > 0:
		name += " "
		namePad += " "
	}

	return appendItems(r, reflect.ValueOf(s), name, namePad, opts...)
}

func AppendSliceOrArray(r []string, a any, opts ...any) []string {
	if t, ok := testframe.Peek[*testing.T](); ok {
		t.Helper()
	}

	var (
		name     opt.Name
		namePad  string
		hasName  bool
		typeName string
		val      = reflect.ValueOf(a)
		kind     = val.Kind()
	)

	switch kind { //nolint:exhaustive // unsupported kinds handled by default case
	case reflect.Array:
		typeName = "array"

	case reflect.Slice:
		typeName = "slice"

	default:
		test.Invalid("report.AppendSliceOrArray requires a slice or array argument")
	}

	if name, opts, hasName = opt.Extract[opt.Name](opts); hasName {
		namePad = strings.Repeat(" ", len(name))
	}

	switch {
	case kind != reflect.Array && val.IsNil():
		return append(r, string(name)+" "+Nil)

	case val.Len() == 0:
		return append(r, string(name)+" <empty "+typeName+">")

	case len(name) > 0:
		name += " "
		namePad += " "
	}

	return appendItems(r, val, name, namePad, opts...)
}

func appendItems(r []string, v reflect.Value, name opt.Name, namePad string, opts ...any) []string {
	var (
		maxItems       = v.Len()
		remainingItems = 0
	)

	if limit, ok := opt.Get[opt.MaxItems](opts); ok && int(limit) < maxItems {
		remainingItems = maxItems - int(limit)
		maxItems = int(limit)
	}

	r = append(r, string(name)+"[ "+Value(v.Index(0).Interface(), opts...))

	for i := 1; i < maxItems; i++ {
		r = append(r, namePad+"[ "+Value(v.Index(i).Interface(), opts...))
	}

	if remainingItems > 0 {
		r = append(r, namePad+fmt.Sprintf("... and %d more", remainingItems))
	}

	return r
}
