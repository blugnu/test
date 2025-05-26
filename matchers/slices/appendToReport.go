package slices

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/blugnu/test/opt"
)

type JustifyPrefix bool

func AppendToReport(r []string, s any, pfx string, opts ...any) []string {
	v := reflect.ValueOf(s)
	switch {
	case v.Kind() != reflect.Slice:
		return append(r, pfx+" <not a slice>")
	case v.IsNil():
		return append(r, pfx+" nil")
	case v.Len() == 0:
		return append(r, pfx+" <empty slice>")
	}

	spec := "%v"
	if reflect.TypeOf(v.Index(0).Interface()).Kind() == reflect.String && !opt.IsSet(opts, opt.QuotedStrings(false)) {
		spec = "%q"
	}

	if opt.IsSet(opts, opt.PrefixInlineWithFirstItem(true)) {
		ifx := strings.Repeat(" ", len(pfx)+1)
		r = append(r, pfx+fmt.Sprintf(" | "+spec, v.Index(0).Interface()))
		for i := 1; i < v.Len(); i++ {
			r = append(r, fmt.Sprintf(ifx+"| "+spec, v.Index(i).Interface()))
		}
		return r
	}

	r = append(r, pfx)
	for i := 0; i < v.Len(); i++ {
		r = append(r, fmt.Sprintf("| "+spec, v.Index(i).Interface()))
	}
	return r
}
