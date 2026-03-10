package expect

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/blugnu/test"
	"github.com/blugnu/test/internal"
	"github.com/blugnu/test/opt"
	"github.com/blugnu/test/report"
)

// MapTest provides assertion methods for maps.
type MapTest[K comparable, V any] struct {
	got   map[K]V
	eopts []any
	mopts []any
}

// Map creates an expectation that may be used to perform assertions on a provided map.
// Any options provided are used for all assertions made using the returned expectation.
//
// # Supported Options
//
//	string                   // a name for the value, for use in any test failure report
//	opt.Name(string)         // }-- equivalents --
//	opt.Namef(s, args)       // }
//
//	opt.FailureReport(func)  // a function returning a custom failure report
//	                         // in the event that the test fails
//
//	opt.OnFailure(string)    // a simple string to output as the
//	                         // failure report if the test fails.
//
//	opt.IsRequired(bool)     // if true, the expectation is required to pass;
//	                         // no further expectations in the current test will be
//	                         // evaluated if the expectation fails.
//
//	opt.Required()           // equivalent to opt.IsRequired(true)
//
// # Failure Report Customization
//
// The built-in failure reports respond to several options to customize
// their output. If supplying a custom failure report function using
// [opt.FailureReport], consider supporting these options in your custom
// report function for consistency with built-in reports.
//
//	report.ValueAsString(func)  // a function that converts values to strings
//	                         // for inclusion in failure reports.
func Map[K comparable, V any](got map[K]V, opts ...any) MapTest[K, V] {
	mt := MapTest[K, V]{got: got}
	mt.eopts, mt.mopts = internal.SeparateOptions(opts)
	return mt
}

// DoesNotHaveKeys asserts that the map does not contain the specified keys.
func (mt MapTest[K, V]) DoesNotHaveKeys(keys ...K) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AnyItem)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		intruders := []K{}
		for _, k := range keys {
			if _, found := mt.got[k]; found {
				intruders = append(intruders, k)
			}
		}

		msg := []string{}
		msg = report.AppendSlice(msg, keys, opt.WithName(opts, "expected map to not have "+report.Pluralise("key", len(keys))+":")...)
		msg = report.AppendSlice(msg, intruders, opt.Force(opts, opt.Name("got:"))...)

		return msg
	}))

	test.Expect(test.KeysOfMap(mt.got), mt.eopts...).ToNot(test.ContainItems(keys), opts...)
}

// DoesNotHaveValues asserts that the map does not contain the specified values.
func (mt MapTest[K, V]) DoesNotHaveValues(values ...V) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AnyItem)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		mapValues := test.ValuesOfMap(mt.got)
		intruders := make([]V, 0, len(values))
		for _, v := range values {
			if found := slices.ContainsFunc(mapValues, func(mv V) bool {
				return reflect.DeepEqual(v, mv)
			}); found {
				intruders = append(intruders, v)
			}
		}
		msg := []string{}
		msg = report.AppendSlice(msg, values, opt.WithName(opts, "expected map to not have "+report.Pluralise("value", len(values))+":")...)
		msg = report.AppendSlice(msg, intruders, opt.Force(opts, opt.Name("got:"))...)
		return msg
	}))

	test.Expect(test.ValuesOfMap(mt.got), mt.eopts...).ToNot(test.ContainItems(values), opts...)
}

// Equals asserts that the map is equal to the expected map.
func (mt MapTest[K, V]) Equals(expected map[K]V) {
	test.T().Helper()
	test.Expect(mt.got, mt.eopts...).To(test.EqualMap(expected), mt.mopts...)
}

// HasAnyKey asserts that the map contains at least one of the specified keys.
func (mt MapTest[K, V]) HasAnyKey(keys ...K) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AnyItem)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = report.AppendSlice(result, keys, opt.WithName(opts, "expected map to contain at least one of these keys:")...)
		result = report.AppendSlice(result, test.KeysOfMap(mt.got), opt.WithName(opts, "got:")...)
		return result
	}))

	test.Expect(test.KeysOfMap(mt.got), mt.eopts...).To(test.ContainItems(keys), opts...)
}

// HasAnyValue asserts that the map contains at least one of the specified values.
func (mt MapTest[K, V]) HasAnyValue(values ...V) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AnyItem)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = report.AppendSlice(result, values, opt.WithName(opts, "expected map to contain at least one of these values:")...)
		result = report.AppendSlice(result, test.ValuesOfMap(mt.got), opt.WithName(opts, "got:")...)
		return result
	}))

	test.Expect(test.ValuesOfMap(mt.got), mt.eopts...).To(test.ContainItems(values), opts...)
}

// HasKey asserts that the map contains the specified key.
func (mt MapTest[K, V]) HasKey(key K) {
	test.T().Helper()

	opts := append(slices.Clone(mt.mopts), opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = report.AppendSlice(result, []K{key}, opt.WithName(opts, "expected map to contain key:")...)
		result = report.AppendSlice(result, test.KeysOfMap(mt.got), opt.Force(opts, opt.Name("got keys:"))...)
		return result
	}))

	test.Expect(test.KeysOfMap(mt.got), mt.eopts...).To(test.ContainItem(key), opts...)
}

// HasKeys asserts that the map contains all of the specified keys.
func (mt MapTest[K, V]) HasKeys(keys ...K) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AllItems)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = report.AppendSlice(result, keys, opt.WithName(opts, "expected map to contain keys:")...)
		result = report.AppendSlice(result, test.KeysOfMap(mt.got), opt.Force(opts, opt.Name("got keys:"))...)
		return result
	}))

	test.Expect(test.KeysOfMap(mt.got), mt.eopts...).To(test.ContainItems(keys), opts...)
}

// HasLength asserts that the map has the specified length (number of keys).
func (mt MapTest[K, V]) HasLength(n int) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{
			fmt.Sprintf("expected map with: %d %s", n, report.Pluralise("key", n)),
			fmt.Sprintf("got              : %d", len(mt.got)),
		}
		return result
	}))

	test.Expect(len(mt.got), mt.eopts...).To(test.Equal(n), opts...)
}

// HasValue asserts that the map contains the specified value. The keys of the
// map are not considered.
func (mt MapTest[K, V]) HasValue(value V) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = append(result, "expected map to contain value: "+report.Value(value, opts...))
		result = report.AppendMap(result, mt.got, opt.Force(opts, opt.Name("got:"))...)
		return result
	}))

	test.Expect(test.ValuesOfMap(mt.got), mt.eopts...).To(test.ContainItem(value), opts...)
}

// HasValues asserts that the map contains all of the specified values. The keys of the
// map are not considered.
func (mt MapTest[K, V]) HasValues(value ...V) {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.AllItems)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{}
		result = report.AppendSlice(result, value, opt.WithName(opts, "expected map to contain values:")...)
		result = report.AppendMap(result, mt.got, opt.Force(opts, opt.Name("got:"))...)
		return result
	}))

	test.Expect(test.ValuesOfMap(mt.got), mt.eopts...).To(test.ContainItems(value), opts...)
}

// IsEmpty asserts that the map is empty. A nil map will fail this assertion. To assert that a
// map is nil or empty, use HasLength(0) or IsEmptyOrNil().
func (mt MapTest[K, V]) IsEmpty() {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := report.AppendMap([]string{}, map[K]V{}, opt.WithName(opts, "expected:")...)
		result = report.AppendMap(result, mt.got, opt.Force(opts, opt.Name("got:"))...)
		return result
	}))

	test.Expect(mt.got, mt.eopts...).Should(test.BeEmpty(), opts...)
}

// IsEmptyOrNil asserts that the map is either nil or empty. To assert that a map is
// empty (but not nil), use IsEmpty().
func (mt MapTest[K, V]) IsEmptyOrNil() {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{"expected: <empty map> or <nil>"}
		result = report.AppendMap(result, mt.got, opt.Force(opts, opt.Name("got:"))...)
		return result
	}))

	test.Expect(mt.got, mt.eopts...).Should(test.BeEmptyOrNil(), opts...)
}

// IsNotEmpty asserts that the map is not empty.
func (mt MapTest[K, V]) IsNotEmpty() {
	test.T().Helper()

	opts := slices.Clone(mt.mopts)
	opts = append(opts, opt.FailureReport(func(opts ...any) []string {
		result := []string{"expected map to not be empty"}
		result = report.AppendMap(result, mt.got, opt.Force(opts, opt.Name("got:"))...)
		return result
	}))

	test.Expect(mt.got, mt.eopts...).ShouldNot(test.BeEmpty(), opts...)
}
