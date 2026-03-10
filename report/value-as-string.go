package report

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/blugnu/test/opt"
)

const Nil = "<nil>"

// Type returns the type of a value as a string, with special handling for some
// types to improve readability.
//
//   - nil values are represented as `<nil>`
//   - `interface {}` is represented as `any`
//   - spaces are removed from around { } in struct type names
func Type(v any) string {
	if v == nil {
		return Nil
	}

	// get the type name as usual
	typeName := fmt.Sprintf("%T", v)

	// replace "interface {}" with "any" for readability
	typeName = strings.ReplaceAll(typeName, "interface {}", "any")

	// remove spaces from around { } in struct type names for readability
	//
	// e.g. struct { a int } becomes struct{a int}"
	typeName = strings.ReplaceAll(typeName, " {}", "{}")
	typeName = strings.ReplaceAll(typeName, " { ", "{")
	typeName = strings.ReplaceAll(typeName, " }", "}")
	return typeName
}

// Value converts a value to string representation, according to the
// provided options. Special handling is provided for some types:
//
//   - nil values are represented as `<nil>`
//   - `string` values are quoted by default (%q), unless [opt.UnquotedStrings] is set (%v)
//   - slices are formatted as type (%T) with length, capacity and items (limited to opt.MaxItems if specified)
//   - arrays and maps are formatted as type (%T) with length and items (limited to opt.MaxItems if specified)
//   - channels are formatted as type (%T) with length
//   - functions and interfaces are formatted as type (%T)
//   - structs are formatted as type (%T) with field names and values
//   - pointers are formatted as dereferenced value (or `<nil>`) and pointer type
//   - all other types are formatted as unquoted strings (%v)
//
// This function should be used in test failure methods to format expected and
// actual values consistently with other matcher failure reports.
//
// Supported options:
//
//	opt.MaxItems(n)         // maximum number of items to display in a collection (default is no limit)
//	opt.QuotedStrings       // strings are formatted using %q (default)
//	opt.PlainStrings        // strings are formatted using %v (default is %q)
//	opt.ValueAsDeclaration  // non-string values are formatted using %#v (default is %v)
//	opt.NoTypeNames         // disables type names in output (default is to include type names)
func Value(v any, opts ...any) string {
	if v == nil {
		return Nil
	}

	// errors have special handling to avoid the verbose struct representation
	// when using fmt.Sprintf("%#v", err)
	//
	// e.g. &errors.errorString{s: "this is an error"}
	//
	// instead, we want to render them as
	//
	//   *errors.errorString "this is an error"
	//
	// the full declaration can be forced using opt.ValueAsDeclaration
	if err, ok := v.(error); ok {
		return errorAsString(err, opts...)
	}

	// if ValueAsDeclaration is set, this overrides all other formatting behaviors
	// and options
	if opt.IsSet(opts, opt.ValueAsDeclaration) {
		return fmt.Sprintf("%#v", v)
	}

	return valueAsString(v, opts...)
}

func arrayAsString(v any, typeName string, opts ...any) string {
	val := reflect.ValueOf(v)

	if val.Len() == 0 {
		return typeName + " (len=0)"
	}

	return fmt.Sprintf("%s (len=%d) [%s]", typeName, val.Len(), valueAsCSV(v, opts...))
}

func defaultValueAndType(v any, typeName string) string {
	return fmt.Sprintf("%v [%s]", v, typeName)
}

func defaultValueAsString(v any, typeName string, opts ...any) string {
	if opt.IsSet(opts, opt.NoTypeNames) ||
		typeName == "bool" ||
		typeName == "float64" ||
		typeName == "int" {
		return fmt.Sprintf("%v", v)
	}

	return defaultValueAndType(v, typeName)
}

// errorAsString formats an error value as a string, with special handling for some cases:
//
// # Standard Formatting (opt.ValueAsDeclaration NOT set):
//
// Errors are rendered as "error message [type]", where "error message" is the result of
// `err.Error()` and "type" is the type of the error.  The error string is always quoted,
// i.e. [opt.UnquotedStrings] is ignored for error values.
//
// # Declaration Formatting (opt.ValueAsDeclaration set):
//
// Errors are rendered as "error message [declaration]", where "error message" is the result of
// `err.Error()` and "declaration" is the full declaration of the error value as returned by
// `fmt.Sprintf("%#v", err)`, with the following special handling:
//
//   - if the declaration is the same as the error message (e.g. for error types that are based
//     on string), it is rendered as "error message [type]" to avoid redundant information
//
// This function is used internally by ValueAsString to format error values.
func errorAsString(err error, opts ...any) string {
	var (
		typeName = Type(err)
		decl     = fmt.Sprintf("%#v", err)
		errq     = fmt.Sprintf("%q", err.Error())
	)

	if decl != errq && opt.IsSet(opts, opt.ValueAsDeclaration) {
		typeName = strings.ReplaceAll(decl, "&", "*")
	}

	return defaultValueAndType(errq, typeName)
}

// pointerAsString formats a pointer value as a string, with special handling for
// nil pointers and type names.
func pointerAsString(v any, typeName string, opts ...any) string {
	val := reflect.ValueOf(v)

	if val.IsNil() {
		if !opt.IsSet(opts, opt.NoTypeNames) {
			return fmt.Sprintf("%s [%s]", Nil, typeName)
		}

		return Nil
	}

	// render the dereferenced value as a string
	value := Value(val.Elem().Interface(), opt.Force(opts, opt.NoTypeNames)...)

	// the pointer-ness of a value is always significant
	//
	// i.e. if the value is a string, we want to render it as "*string "value""
	// rather than just "value" so that it is clear that the value is a pointer
	// to a string, not just a string value.
	if !opt.IsSet(opts, opt.NoTypeNames) && !strings.HasPrefix(value, typeName[1:]) {
		return defaultValueAndType(value, typeName)
	}

	// the rendered value already has the type name (or type names are disabled),
	// so simply prepend a "*" to indicate that it is a pointer
	return "*" + value
}

// sliceAsString formats a slice value as a string, with special handling for
// length, capacity and items.
func sliceAsString(v any, typeName string, opts ...any) string {
	val := reflect.ValueOf(v)

	if val.Len() == 0 {
		switch val.Cap() {
		case 0:
			return typeName + " (len,cap=0)"

		default:
			return fmt.Sprintf("%s (len=%d, cap=%d)", typeName, val.Len(), val.Cap())
		}
	}

	// since the slice already indicates the type of the items in the slice,
	// the items themselves are rendered as values (not declarations)
	//
	// if the slice has items of `any` type, they are rendered with type names,
	// otherwise type names are omitted for readability since the slice type
	// already indicates the item type
	itemsOpts := opt.Unset(opts, opt.ValueAsDeclaration)
	switch {
	case strings.HasSuffix(typeName, "]any"):
		itemsOpts = opt.Unset(itemsOpts, opt.NoTypeNames)
	default:
		itemsOpts = opt.Force(opts, opt.NoTypeNames)
	}

	switch {
	case val.Len() == val.Cap():
		return fmt.Sprintf("%s (len,cap=%d) [%s]",
			typeName,
			val.Len(),
			valueAsCSV(v, itemsOpts...),
		)

	default:
		return fmt.Sprintf("%s (len=%d, cap=%d) [%s]",
			typeName,
			val.Len(),
			val.Cap(),
			valueAsCSV(v, itemsOpts...),
		)
	}
}

// stringAsString formats a string value as a string, with special handling for
// quoting and type names.
func stringAsString(v any, typeName string, opts ...any) string {
	var s string

	switch {
	case opt.IsSet(opts, opt.UnquotedStrings):
		s = fmt.Sprintf("%v", v)
	default:
		s = fmt.Sprintf("%q", v)
	}

	if opt.IsSet(opts, opt.NoTypeNames) ||
		typeName == "string" {
		return s
	}

	return defaultValueAndType(s, typeName)
}

// structAsString formats a struct value as a string, with special handling for field
// names and values.
func structAsString(v any, typeName string, opts ...any) string {
	if typeName == "struct{}" {
		return typeName
	}

	val := reflect.ValueOf(v)

	// for structs, we want to render the field names and values
	// using the type name and the field names and values
	// e.g. "report_test.testcase {value: <nil>, opts: []any (len,cap=0), result: ""}"
	addressable := reflect.New(val.Type()).Elem()
	addressable.Set(val)

	fieldOpts := opt.Force(opts, opt.QuotedStrings)
	fieldOpts = opt.Force(fieldOpts, opt.NoTypeNames)

	fields := make([]string, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name

		field := addressable.Field(i)
		field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem() //nolint: gosec // unsafe is required to read unexported fields

		value := Value(field.Interface(), fieldOpts...)
		fields[i] = fmt.Sprintf("%s:%s", fieldName, value)
	}

	if opt.IsSet(opts, opt.NoTypeNames) {
		return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
	}

	return fmt.Sprintf("{%s} [%s]", strings.Join(fields, ", "), typeName)
}

// valueAsCSV formats a value as a comma-separated list of items, with
// special handling for slices, arrays and maps. It respects the
// opt.MaxItems option to limit the number of items displayed.
//
// For slices and arrays, it formats the items as a comma-separated list,
// limited to any specified opt.MaxItems. If there are more items than the
// limit, it appends "... and N more" to indicate the number of additional
// items.
//
// For maps, it formats each key-value pair as "key => value", also limited
// by the maximum number of items.
//
// The function returns a string representation of the items, suitable for
// inclusion in test failure reports.
//
// Note: This function is used internally by ValueAsString to format slices,
// arrays and maps.
//
// Examples:
//
//	valueAsCSV([]int{1, 2, 3, 4, 5}, opt.MaxItems(3)) // returns "[1, 2, 3, ... and 2 more]"
//
//	valueAsCSV(map[string]int{"a": 1, "b": 2, "c": 3}, opt.MaxItems(2)) // returns "[a => 1, b => 2, ... and 1 more]"
func valueAsCSV(v any, opts ...any) string {
	val := reflect.ValueOf(v)
	kind := reflect.TypeOf(v).Kind()

	maxLen := val.Len()
	remaining := 0
	if limit, ok := opt.Get[opt.MaxItems](opts); ok && int(limit) < maxLen {
		remaining = maxLen - int(limit)
		maxLen = int(limit)
	}

	opts = opt.Unset(opts, opt.ValueAsDeclaration)

	items := make([]string, 0, maxLen)
	switch kind { //nolint:exhaustive // only to be used with arrays, slices and maps
	case reflect.Array, reflect.Slice:
		opts = opt.Force(opts, opt.QuotedStrings)
		for i := range maxLen {
			item := val.Index(i).Interface()
			items = append(items, Value(item, opts...))
		}

	case reflect.Map:
		keys := val.MapKeys()
		for i := range keys[:min(maxLen, len(keys))] {
			key := keys[i]
			item := val.MapIndex(key).Interface()
			items = append(items, fmt.Sprintf("%s => %s", Value(key.Interface(), opts...), Value(item, opts...)))
		}
	}

	if remaining > 0 {
		items = append(items, fmt.Sprintf("... and %d more", remaining))
	}

	return strings.Join(items, ", ")
}

func valueAsString(v any, opts ...any) string {
	var (
		typeName = Type(v)
		kind     = reflect.TypeOf(v).Kind()
	)

	switch kind { //nolint:exhaustive // default handling for non-special cases
	case reflect.Func, reflect.Interface:
		return typeName

	case reflect.Chan:
		return fmt.Sprintf("%s (len=%d)", typeName, reflect.ValueOf(v).Len())

	case reflect.Array, reflect.Map:
		return arrayAsString(v, typeName, opts...)

	case reflect.Pointer:
		return pointerAsString(v, typeName, opts...)

	case reflect.Slice:
		return sliceAsString(v, typeName, opts...)

	case reflect.String:
		return stringAsString(v, typeName, opts...)

	case reflect.Struct:
		return structAsString(v, typeName, opts...)

	default:
		return defaultValueAsString(v, typeName, opts...)
	}
}
