package report

import (
	"fmt"
	"reflect"
	"strings"

	"slices"

	"github.com/blugnu/test/opt"
)

// AppendMap appends a formatted representation of the provided map `m` to the
// `result` slice of strings. It handles nil and empty maps, and formats the map
// entries in a stable order (sorted keys).
func AppendMap[K comparable, V any](result []string, m map[K]V, opts ...any) []string {
	var (
		hasName bool
		name    opt.Name
		namePad string
	)
	if name, opts, hasName = opt.Extract[opt.Name](opts); hasName {
		namePad = strings.Repeat(" ", len(name))
	}

	switch {
	case m == nil:
		return append(result, string(name)+" <nil>")
	case len(m) == 0:
		return append(result, string(name)+" <empty map>")
	case len(name) > 0:
		name += " "
		namePad += " "
	}

	// for stable ordering of the map, we first transform the map keys to strings,
	// sort the keys and append the rendered map in key order
	sortedMap := make(map[string]any, len(m))
	sortedKeys := make([]string, 0, len(m))
	for k, v := range m {
		key := Value(k, opts...)
		sortedKeys = append(sortedKeys, key)
		sortedMap[key] = v
	}
	slices.Sort(sortedKeys)

	var renderValue func(k string, v any) []string

	render := func(v any, name opt.Name) []string {
		valueKind := reflect.ValueOf(v).Kind()
		if valueKind == reflect.Slice || valueKind == reflect.Array {
			return AppendSliceOrArray(nil, v, opt.Force(opts, name)...)
		}
		return []string{fmt.Sprintf("%s %s", name, Value(v, opts...))}
	}

	renderValue = func(k string, v any) []string {
		renderValue = func(k string, v any) []string {
			return render(v, opt.Namef(namePad+"( %s =>", k))
		}

		return render(v, opt.Namef(string(name)+"( %s =>", k))
	}

	for _, k := range sortedKeys {
		result = append(result, renderValue(k, sortedMap[k])...)
	}

	return result
}
