package test

import (
	"fmt"
	"reflect"
)

// Expecting[T] is a struct that holds a value of type T which may be
// accessed via an Expected() method. It may be used by Matcher
// implementations to embed an Expected interface implementation:
//
//	type MyMatcher[T any] struct {
//		Expecting[T]
//	}
type Expecting[T any] struct {
	value T
}

// Expected returns the value of the Expecting struct.
func (e Expecting[T]) Expected() any {
	return e.value
}

// expectation[T] is a type that represents an expectation in a test. It
// holds the test runner in scope at the time the expectation was
// expressed, the name of the expectation (optional) and the value to
// which the expectation applies.
type expectation[T any] struct {
	t    TestRunner
	got  T
	name string
}

func (e expectation[T]) fail(matcher Matcher[T], formatter any) {
	GetT().Helper()

	if formatter == nil {
		formatter = matcher
	}

	switch f := formatter.(type) {
	case CustomOneLineReport:
		e.Error([]string{f.Format()})
	case CustomReport:
		e.Error(f.Format())
	case CustomFormatExpectedAndGot:
		e.Error(f.Format(f.Expected(), e.got))
	case FormatExpectedAndGot:
		ex := f.Format(f.Expected())
		gt := f.Format(e.got)

		// if the expected and got values are small, use a one line report
		if len(ex) < 10 && len(gt) < 10 {
			e.Errorf("expected %s, got %s", ex, gt)
			return
		}
		// otherwise, use a multi-line report
		e.Errorf("expected: %s\ngot     : %s", ex, gt)
	case OneLineExpected:
		e.Errorf("expected %v, got %v", f.Expected(), e.got)
	case Expected:
		e.Errorf("expected: %v\ngot     : %v", f.Expected(), e.got)
	default:
		e.Errorf("got %v", e.got)
	}
}

// Error fails a test with an optional message.  If specified, the
// message may be given as a string or []string.  Only the first msg
// is used; any additional msg args are ignored.
//
// If no msg is supplied, the test fails with no message.  If the
// expectation has a name, the test fails with the message "<name> failed".
//
// If the first msg is a string it is used as the message. If the
// expectation has a name, it is prepended to the string.
//
// If the first msg is a []string, it is used as the message with each
// string in the slice on a new line. If the expectation has a name,
// it is prepended to the first string in the slice.
func (e expectation[T]) Error(msg ...any) {
	e.t.Helper()

	if len(msg) == 0 {
		if e.name != "" {
			e.t.Errorf("\n" + e.name + " failed")
			return
		}
		e.t.Fail()
		return
	}

	switch msg := msg[0].(type) {
	case string:
		if e.name != "" {
			msg = e.name + ": " + msg
		}
		e.t.Error(msg)

	case []string:
		if len(msg) == 0 {
			e.Error()
			return
		}

		var rpt string
		if e.name != "" {
			msg[0] = e.name + ": " + msg[0]
		}
		for _, s := range msg {
			rpt += "\n" + s
		}
		e.t.Error(rpt)
	}
}

// Errorf fails a test with an optional message.  If specified, the
// message may be given as a string containing formatting verbs and
// a variadic number of arguments.
//
// If a message string is given without formatting verbs, it is used
// as the message and any additional args are ignored.
//
// If the expectation is named, the name is prepended to the message.
func (e expectation[T]) Errorf(s string, args ...any) {
	e.t.Helper()
	if e.name != "" {
		e.t.Errorf("\n"+e.name+": "+s, args...)
	} else {
		e.t.Errorf("\n"+s, args...)
	}
}

// Expect creates an expectation for the given value.  The value
// may be of any type.  Additional options may be applied to the
// expectation as follows:
//
//	string    // a name for the expectation; the name is used in
//		      // the failure message if the expectation fails.
func Expect[T any](value T, opts ...any) expectation[T] {
	t, _ := TestFrame()
	if t == nil {
		panic(fmt.Errorf("Expect: %w", ErrNoTestFrame))
	}

	return expectation[T]{t, value, optName(opts)}
}

// func (e expectation[T]) ToNot(matcher Matcher[T], opts ...any) {
// 	GetT().Helper()

// 	if matcher.Match(e.got) {
// 		formatter := getFormatter(opts...)

// 		switch special := matcher.(type) {
// 		case interface{ MatchAny(any) bool }:
// 			if !special.MatchAny(e.got) {
// 				e.fail(matcher, formatter)
// 			}
// 		default:
// 			e.fail(matcher, formatter)
// 		}
// 	}
// }

// To applies a matcher to the expectation.  If the matcher
// does not match the value, the test fails.  The matcher may
// be any type that implements the Matcher interface.
//
// A matcher may also implement the MatchAny interface, to
// match using a value that is different to the type of the
// value in the expectation itself.
//
// Options may be specified as follows:
//
//	TBC
//
// Any options additional to those above are passed to the
// matcher to use as it sees fit. Refer to the documentation
// for the matcher for details.
func (e expectation[T]) To(matcher Matcher[T], opts ...any) {
	GetT().Helper()

	if matcher.Match(e.got, opts...) {
		switch special := matcher.(type) {
		case interface{ MatchAny(any, ...any) bool }:
			if special.MatchAny(e.got, opts...) {
				return
			}
		default:
			return
		}
	}
	e.fail(matcher, getFormatter(opts...))
}

// Is tests the value of the expectation against some expected
// value.  The function behaves differently depending on the values
// and types of the expected and actual values:
//
//   - If both expected and the value being tested are nil, the test
//     passes.
//
//   - If either the expected value or the value being tested is nil and
//     the other is not, the test fails.
//
//   - If both the expected value and the value being tested implement
//     the error interface, errors.Is() is used to compare the two
//     values.
//
// Otherwise, non-nil values are compared using DeepEqual.
func (e expectation[T]) Is(expected T) {
	GetT().Helper()

	switch {
	case any(expected) == nil && any(e.got) == nil:
		return
	case any(expected) == nil && any(e.got) != nil:
		e.Errorf("expected nil, got %v", e.got)
	case any(expected) != nil && any(e.got) == nil:
		e.Errorf("expected %v, got nil", expected)
	default:
		experr, _ := any(expected).(error)
		goterr, _ := any(e.got).(error)
		if experr != nil && goterr != nil {
			ExpectError(goterr, e.name).Is(experr)
			return
		}
		Expect(e.got, e.name).To(DeepEqual(expected))
	}
}

// IeEmpty checks that the value of the expectation is empty.
func (e expectation[T]) IsEmpty() {
	GetT().Helper()

	// regexMatch := func(regex, s string) map[string]string {
	// 	var comp = regexp.MustCompile(regex)
	// 	match := comp.FindStringSubmatch(s)

	// 	vars := make(map[string]string)
	// 	for i, name := range comp.SubexpNames() {
	// 		if i > 0 && i <= len(match) {
	// 			vars[name] = match[i]
	// 		}
	// 	}
	// 	return vars
	// }

	valsFn := func(l int) []string {
		// vars := regexMatch(`(?P<type>.*)\{(?P<values>.*)\}`, fmt.Sprintf("%#v", e.got))
		// vals := strings.Replace(vars["values"], `", "`, "\"\n\"", -1)

		report := []string{
			fmt.Sprintf("expected empty %T, got length %d", e.got, l),
		}
		// if len(vals) > 0 {
		// 	report = append(report, "------------")
		// 	report = append(report, strings.Split(vals, "\n")...)
		// 	report = append(report, "------------")
		// }
		return report
	}

	gt := reflect.TypeOf(e.got)
	gv := reflect.ValueOf(e.got)
	l := 0
	m := ""
	switch gt.Kind() {
	case reflect.Array:
		l = gt.Len()
	case reflect.Chan, reflect.Map, reflect.Slice:
		l = gv.Len()
	default:
		valsFn = func(l int) []string {
			return []string{
				fmt.Sprintf("expected empty %T.%s, got length %d", e.got, m, l),
			}
		}

		switch got := any(e.got).(type) {
		case string:
			l = len(got)
			valsFn = func(l int) []string {
				return []string{fmt.Sprintf("expected empty string, got length %d", l),
					fmt.Sprintf("  value: %q", got),
				}
			}

		case interface{ Count() int }:
			l = got.Count()
			m = "Count() int"
		case interface{ Len() int }:
			l = got.Len()
			m = "Len() int"
		case interface{ Length() int }:
			l = got.Length()
			m = "Length() int"
		case interface{ Count() int64 }:
			l = int(got.Count())
			m = "Count() int64"
		case interface{ Len() int64 }:
			l = int(got.Len())
			m = "Len() int64"
		case interface{ Length() int64 }:
			l = int(got.Length())
			m = "Length() int64"
		case interface{ Count() uint }:
			l = int(got.Count())
			m = "Count() uint"
		case interface{ Len() uint }:
			l = int(got.Len())
			m = "Len() uint"
		case interface{ Length() uint }:
			l = int(got.Length())
			m = "Length() uint"
		case interface{ Count() uint64 }:
			l = int(got.Count())
			m = "Count() uint64"
		case interface{ Len() uint64 }:
			l = int(got.Len())
			m = "Len() uint64"
		case interface{ Length() uint64 }:
			l = int(got.Length())
			m = "Length() uint64"
		default:
			Expect(got).Error([]string{
				fmt.Sprintf("IsEmpty: invalid test for type %T:", e.got),
				"  tested value must be string, array, slice, channel, map",
				"  or a type that implements Count(), Len() or Length()",
				"  returning int, int64, uint, or uint64",
			})
		}
	}

	if l > 0 {
		e.Error(valsFn(l))
	}
}

// IsNil checks that the value of the expectation is nil.  If the
// value is not nil, the test fails.  If the value is nil, the test
// passes.
//
// If the value being tested does not support a nil value the test
// will fail and produce a report similar to:
//
//	test.IsNil: invalid test:
//	   values of type '<type>' are not nilable
func (e expectation[T]) IsNil() {
	switch result := isNil(e.got).(type) {
	case bool:
		if !result {
			switch got := any(e.got).(type) {
			case error:
				e.Errorf("expected nil, got error: %v", got)
			default:
				if reflect.ValueOf(e.got).Kind() == reflect.Pointer {
					v := reflect.Indirect(reflect.ValueOf(e.got))
					if v.Kind() == reflect.String {
						e.Errorf("expected nil, got &(%q)", v)
						return
					}
					e.Errorf("expected nil, got &(%#v)", v)
					return
				}
				e.Errorf("expected nil, got %#v", got)
			}
		}
	case error:
		e.Error([]string{
			"test.IsNil: invalid test:",
			fmt.Sprintf("  values of type '%T' are not nilable", e.got),
		})
	}
}
