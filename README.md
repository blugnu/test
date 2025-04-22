<div align="center" style="margin-bottom:20px">
  <!-- <img src=".assets/banner.png" alt="logger" /> -->
  <div align="center">
    <a href="https://github.com/blugnu/test/actions/workflows/release.yml"><img alt="build-status" src="https://github.com/blugnu/test/actions/workflows/release.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/blugnu/test" ><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/test"/></a>
    <a><img alt="go version >= 1.20" src="https://img.shields.io/github/go-mod/go-version/blugnu/test?style=flat-square"/></a>
    <a href="https://github.com/blugnu/test/blob/master/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/blugnu/test?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/blugnu/test?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/test?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/blugnu/test"><img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/test"/></a>
    <hr/>
  </div>
</div>

<br>

# blugnu/test

Provides test helpers for use with the standard library `testing` package; it is not intended
to be a replacement for the `testing` package or a complete testing framework in its own right.

# CAUTION

Feel free to use this package if you find it useful, but be aware that it is still in development
and the API may change without notice.  The package is in active use and is constantly being revised
and refined as problems and annoyances are identified and resolved in the API.

The API will remain as stable as possible, but until the package hits v1.0 this is only an aspiration,
not a commitment.

## Installation

```bash
go get github.com/blugnu/test
```

## Quick Start

To perform tests use either a testable factory or a test function:

| Category | Description |
| --- | --- |
| [Testable Factories](#testable-factories) | functions that return a '_testable_' value providing functions to perform tests on that value |
| [Test Helpers](#test-helpers) | functions that directly perform a test and report the outcome |

To quickly understand the difference you might find it helpful to read: [Testables vs Test Helpers](#testables-vs-test-helpers)

### Basic Usage

- [Tests for Errors and Comparable Values](#tests-for-errors-and-comparable-values)
- [Testing Maps and Slices](#testing-maps-and-slices)
- [Mocking Functions](#mocking-functions)

### Advanced Usage

In addition to performing common, basic tests, the `test` package also provides support for more advanced testing scenarios:

| Category | Description |
| --- | --- |
| [Capture and Test Console Output](#capture-and-test-console-output) | capture output of a function that writes to `stdout` and/or `stderr` |
| [Mocking Functions](#mocking-functions) | mock functions for testing |
| [Test for Expected Panics](#test-for-expected-panics) | test that a function panics as expected |
| [Test for an Expected Type](#test-for-an-expected-type) | test that a value is of an expected type |
| [Testing a Test Helper](#testing-a-test-helper) | test your own test helper functions |
| [Testing Context Values](#testing-context-values) | test values stored in a `context.Context` |

<br/>
<hr/>

# Testable Factories

Testable factories are used to create testable values ('testables') that provide test functions
appropriate to, or specialised for, the type of value being tested.

<!-- markdownlint-disable MD013 // line length -->
| Factory Function | Description |
| --- | --- |
| `test.Bytes(t *testing.T, got []byte, opts ...any) *Bytes` | returns a testable `[]byte` |
| `test.Error(t *testing.T, got error, opts ...any) *Error` | returns a testable `error` |
| `test.Map[K comparable, V any](t *testing.T, got map[K]V, opts ...any) *Map[K, V]` | returns a testable `map[K,V]` |
| `test.Slice[T comparable](t *testing.T, got []T, opts ...any) *Slice[T]` | returns a testable slice of values satisfying the `comparable` constraint |
| `test.Strings(t *testing.T, got []string, opts ...any) *Strings` | returns a testable `[]string` |
| `test.Value[T comparable](t *testing.T, got T, opts ...any) *Value[T]` | returns a testable value of a type satisfying the `comparable` constraint |
<!-- markdownlint-enable MD013 -->

Note that some testable factories are generic functions with a constrained type parameter
which may make them [unsuitable for use with certain values](#working-with-or-around-constraints).

Testable factories accept a minimum of two arguments:

- `t *testing.T` - the `*testing.T` to be used by any test functions provided by the testable
- `got` - the value to be tested

Note that the type parameter for generic testable factories is able to be inferred
from the `got` parameter; there is no need to specify the type.  For example (assuming
that `DoSomething` returns a value of a type that satisfies the `comparable` constraint):

```golang
func TestDoSomething(t *testing.T) {
  // ACT
  got := DoSomething()

  // ASSERT
  test.Value(t, got).Equals("foo")
}
```

Testable factories also support additional options that may be provided as additional parameters.
For details of the options supported by each testable, see [Testable Factory Options](#testable-factory-options).

## Working With (or Around) Constraints

If you need to test a value which does not satisfy the type constraint of a testable factory, it
should be possible to implement an equivalent test using a factory that is specialised for the
value involved, one that is unconstrained, or by using a test helper.  For example, `[]byte` does
not satisfy the `comparable` constraint and so cannot be tested using a `test.Value()` testable.
Alternatives in this case are:

- use the `test.Bytes()` testable factory (_testable factory specialised for `[]byte`_)
- use the `test.Slice[byte]()` testable factory (_`[]byte` does not satisfy `comparable`, but `byte` does_)
- use the `test.DeepEqual()` test helper (_unconstrained test helper_)

## Testable Factory Options

Testable factory options are always passed after the mandatory parameters.  Optional parameters
are discriminated by _type_. The following types are supported:

<!-- markdownlint-disable MD013 // line length -->
| Parameter Type | Name | Description | Notes |
| --- | --- | --- | --- |
| `string` | _value name_ | a name for the value being tested | |
| `test.Format` | _format verb_ | the format verb to be used when manifesting values in a test failure report  | _ignored if a _format function_ is specified_ |
| `test.BytesFormat` | _format verb_ | the format verb for formatting `[]byte` values in a test failure report | only supported by `test.Bytes()`<br/><br/>_ignored if a _format function_ is specified_ |
| `func(*testing.T) string` | _format function_ | a function that returns a string representation of the value being tested.  _The type T varies according to the type of the testable value_. | not supported by `test.Bool()` |
<!-- markdownlint-enable MD013 -->

If multiple values of any of these types are supplied in a given call to a factory only the first
is significant.  For example, in the following call the additional `"some other value"` name
parameter (`string`) will be ignored:

```go
test.Value(t, got, "some value", "some other value").Equals(expected)
```

<br>
<hr>

# Test Helpers

Test helpers are functions that directly perform a test and report the outcome.

<!-- markdownlint-disable MD013 // line length -->
| Test Helper | Description |
| --- | --- |
| `test.DeepEqual[T any](t *testing.T, got, wanted *testing.T, opts ...any)` | fails if `got` is not equal to `wanted`, based on `reflect.DeepEqual()` comparison |
| `test.Equal[T comparable](t *testing.T, got, wanted *testing.T, opts ...any)` | fails if `got` is not equal to `wanted` |
| `test.IsNil(t *testing.T, got any, name ...string)` | fails if `got` is not `nil` |
| `test.IsNotNil(t *testing.T, got any, name ...string)` | fails if `got` is `nil` |
| `test.NotDeepEqual[T any](t *testing.T, got, wanted *testing.T, opts ...any)` | fails if `got` is equal to `wanted`, based on `reflect.DeepEqual()` comparison |
| `test.NotEqual[T comparable](t *testing.T, got, wanted *testing.T, opts ...any)` | fails if `got` is equal to `wanted` |
<!-- markdownlint-enable MD013 -->

With the exception of `IsNil` and `IsNotNil`, a test helper accepts a minimum of _three_ arguments:

- `t *testing.T` - the `*testing.T` value passed to the test function
- `got` - the value to be tested
- `wanted` - the value to be compared with `got`

Additional parameters are optional and are always passed after the mandatory parameters.

If multiple optional parameters are supported, they are discriminated by _type_.  The following
optional parameters are supported:

<!-- markdownlint-disable MD013 // line length -->
| Parameter Type | Name | Description | Notes |
| --- | --- | --- | --- |
| `string` | _name_ | a name for the test | this is the _only_ parameter supported by `IsNil` and `IsNotNil` |
| `test.Equality` | _comparison method_ | the method used to compare `got` and `wanted` | only supported by `Equal` and `NotEqual`<br/><br/>_ignored if a _comparison function_ is specified_ |
| `test.Format` | _format verb_ | the format verb to be used when manifesting values in a test failure report  | _ignored if a _format function_ is specified_ |
| `func(got, wanted T) bool` | _comparison function_ | a function that compares two values of type `T` for equality, returning `true` if considered equal otherwise `false`.  _The type T varies according to the type of the testable value_. | only supported by `test.IsEqual` and `test.NotEqual` |
| `func(*testing.T) string` | _format function_ | a function that returns a string representation of the value being tested.  _The type T varies according to the type of the testable value_. | not supported by `test.IsNil` or `test.IsNotNil` |
<!-- markdownlint-enable MD013 -->

If multiple values of a given type are supplied, only the first is significant.  For example,
in the following call the `"some other value"` name parameter will be ignored:

```go
  test.Equal(t, got, wanted, "some value", "some other value")
```

<br>
<hr>

# Additional Information

## Testables vs Test Helpers

There are often multiple ways of performing a given test using either a testable or a helper function.
There is no "right" way; use whichever is most appropriate or intuitive in a specific case.

For example, to test that an `error` returned by a function is `nil` you could use either of the tests
illustrated here:

```go
func TestDoSomething(t *testing.T) {
  // ACT
  err := DoSomething()

  // ASSERT
  test.IsNil(t, err)
  test.Equal(t, err, nil)
  test.Error(t, err).IsNil()
}
```

All three tests in this example are testing the same thing, though the first two use test helpers
while the third uses a `test.Error` testable.

The third test is strongly typed and is arguably more readable and intuitive than the first two,
but the first two are more concise.

### Naming a SUT (Subject Under Test)

All factory functions support an optional `string` parameter to provide a name for the value being
tested.  If not specified, each factory function will assume a default name.

Example:
  
  ```go
  func TestDoSomething(t *testing.T) {
    // ACT
    got := DoSomething()

    // ASSERT
    test.Value(t, got).Equals(expected)           // will produce a test named: TestDoSomething/value/equals
    test.Value(t, got, "result").Equals(expected) // will produce a test named: TestDoSomething/result/equals
  }
  ```

## Features and Examples

- [blugnu/test](#blugnutest)
- [CAUTION](#caution)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
    - [Basic Usage](#basic-usage)
    - [Advanced Usage](#advanced-usage)
- [Testable Factories](#testable-factories)
  - [Working With (or Around) Constraints](#working-with-or-around-constraints)
  - [Testable Factory Options](#testable-factory-options)
- [Test Helpers](#test-helpers)
- [Additional Information](#additional-information)
  - [Testables vs Test Helpers](#testables-vs-test-helpers)
    - [Naming a SUT (Subject Under Test)](#naming-a-sut-subject-under-test)
  - [Features and Examples](#features-and-examples)
    - [Tests for Errors and Comparable Values](#tests-for-errors-and-comparable-values)
    - [Testing Maps and Slices](#testing-maps-and-slices)
      - [test.Map](#testmap)
      - [test.Slice](#testslice)
      - [test.Bytes](#testbytes)
    - [Mocking Functions](#mocking-functions)
      - [Fake Function Results](#fake-function-results)
    - [Capture and Test Console Output](#capture-and-test-console-output)
    - [Testing a Test Helper](#testing-a-test-helper)
    - [Testing Context Values](#testing-context-values)
    - [Test for Expected Panics](#test-for-expected-panics)
    - [Test for an Expected Type](#test-for-an-expected-type)

### Tests for Errors and Comparable Values

A `test.Error()` factory is provided that returns a testable `error` supporting the following tests:

- `Is(wanted)` - fails if the error is not `wanted` (using `errors.Is()`)
- `IsNil()` - fails if the error is not `nil`

In addition, the `test.IsNil()` function provides specific support for testing for `nil` errors and
so may be more convenient to use when performing a simple test for an unexpected error.

The following snippets demonstrate these tests:

```go
func TestDoSomething(t *testing.T) {
  // ACT
  err := DoSomething()
  test.IsNil(t, err)  // will fail with "unexpected error: <type>: <error string>" if err is not nil

  // ASSERT
  test.Error(t, err, "returned error").IsNil()    // equivalent to the above but with an explicit name
  test.Error(t, err).Is(io.EOF)                   // equivalent to errors.Is(err, io.EOF)
}
```

> _NOTE: If the value supplied to the `test.IsNil()` function is of a type that does not have a
> meaningful `nil` value, the test will fail as an invalid test. Types that may be tested for
> `nil` are: `chan`, `func`, `interface`, slices, maps, and pointers._

For values of a comparable type, the `test.Value[T comparable]()` factory returns a testable
value of a type satisfying the `comparable` constraint.  The returned value provides the
following tests:

- `Equals(wanted)` - fails if the value is not equal to `wanted`
- `IsNil()` - fails if the value is not `nil`
- `IsNotNil()` - fails if the value is `nil`

```go
func TestDoSomething(t *testing.T) {
  // ACT
  got, err := DoSomething()
  test.IsNil(err)

  // ASSERT
  test.Value(t, "returned value", got).Equals("foo")  // fails if got is not "foo"
}
```

> _NOTE: If the value supplied to the `test.Value()` factory function is of a type that does not
> have a meaningful `nil` value, the `IsNil()` and `IsNotNil()` tests will fail as invalid.
> Types that may be tested for `nil` are: `chan`, `func`, `interface`, slices, maps, and pointers._

### Testing Maps and Slices

Three factory functions are provided for creating values for testing maps and slices:

- `test.Map[K comparable, V any]()` for testing a map
- `test.Slice[T comparable]()` for testing a slice of values satisfying the `comparable` constraint
- `test.Bytes()` for testing a `[]byte` specifically

> _Tests that are available using a `test.Bytes()` test could also be performed using a `test.Slice[byte]()`.
> However, a `test.Bytes()` test provides options that are more useful when working with `[]byte` values,
> together with more helpful formatting of test failure reports._

#### test.Map

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := map[string]string{
    "foo": "bar",
  }

  // ACT
  got := DoSomething()

  // ASSERT
  test.Map(t, got).Equals(expected)
}
```

#### test.Slice

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := []string{"foo", "bar"}

  // ACT
  got := DoSomething()

  // ASSERT
  test.Slice(t, got).Equals(expected)
}
```

#### test.Bytes

When testing `[]byte`, an optional format argument may be used to specify the format of the
expected and actual values in any test failure report.

The default format is `BytesHex` (hexadecimal):

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := []byte("foo")
  got := []byte("bar")

  // ACT & ASSERT
  test.Bytes(t, got).Equals(expected)               // displays values in a failure report as hexadecimal
  test.Bytes(t, got, test.BytesBinary).Equals(expected)  // displays values in a failure report as binary
}
```

Any format may be specified by casting a string as a `BytesFormat` if needed; sensible values
are provided as constants.

### Mocking Functions

When testing functions that call other functions, it is often necessary to mock the functions
being called to ensure that the tests are isolated and that the functions being tested are
not dependent on the behaviour of the functions being called.

The `test.MockFn[A, R]` type provides a way to mock a function accepting arguments of type A and
returning a result of type R.

All `test.MockFn` values support an optional `error` value which may simply be ignored/not used
if the mocked function does not return an error.

If the function being mocked does not return any value other than an `error`, the result type `R`
should be `any` and ignored.  Similarly if the function being mocked does not require any
arguments, the argument type `A` should be `any` and ignored.

#### Fake Function Results

The `test.MockFn` type can provide fake results for a mocked function.  Fake results may be setup
in two different ways:

- expected calls mode.  `ExpectCall()` is used configure an expected call; this returns a value
  with a `WillReturn` method to setup a result to be returned for that call.  In this mode, calls
  to the mocked function that do not match the expected calls will cause the test to fail.

- mapped result more.  `WhenCalledWith(args A)` is used to setup a result to be returned when the
  mocked function is called with the specified arguments.  In this mode, calls to the mocked
  function that do not match any of the mapped results will cause the test to fail.

```go

#### Multiple Arguments/Result Values

If a function being mocked accepts multiple arguments and/or returns multiple result values (in
addition to an error), the types A and/or R should be a `struct` type with fields for the arguments
and result values required:

```go
type fooArgs struct {
  A int
  B string
}

type fooResult struct {
  X int
  Y string
}

type mockFoo struct {
  foo test.MockFn[fooArgs, fooResult]
}

func (mock *mockFoo) Foo(A int, B string) (int, string, error) {
  result, err := mock.foo.RecordCall(fooArgs{A, B})
  return result.X, result.Y, err
}
```




### Capture and Test Console Output

The `test.CaptureOutput` function captures the output of a function that writes to `stdout` and/or `stderr`
and returns the captured output as a `CapturedOutput` value.

The `CapturedOutput` value provides the following tests:

- `Contains(wanted)` - fails if the captured output does not contain `wanted`
- `DoesNotContain(wanted)` - fails if the captured output contains `wanted`
- `Equals(wanted)` - fails if the captured output is not equal to `wanted`
- `IsNil()` - fails if the captured output is not `nil`
- `IsNotNil()` - fails if the captured output is `nil`

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  var err error

  // ACT
  stdout, stderr :=test.CaptureOutput(t, func (*testing.T) {
    _, err := DoSomething()
    return err
  })

  // ASSERT
  test.UnexpectedError(t, err)
  test.Equals(t, "foo", got)
  stdout.Contains("some expected log message")
}
```

### Testing a Test Helper

The `test.Helper` function combines the execution of a test helper with the testing of the
outcome of the helper.  The outcome of the helper is specific using `test.ShouldPass` or
`test.ShouldFail` or providing a `test.*Panic` if the helper is expected to panic.

The output of the test helper is returned as `CapturedOutput` (both `stdout` and `stderr`)
so that the presentation of test failure messages in the log can also be tested and verified.

```go
func TestUnexpectedError(t *testing.T) {
  // ARRANGE
  err := errors.New("some error")

  // ACT & ASSERT
  stdout, stderr := test.Helper(t, func(st *testing.T) {
    test.UnexpectedError(st, err)
  }, test.ShouldPass)

  stdout.Contains(nil)  // no output expected for a PASS
}
```

> **NOTE:** _It is important that the helper function being tested is called with
> the `*testing.T` passed to the function that runs it (`st` in the example above)
> and not the `T` of the test (`t` in the example)._

### Testing Context Values

When a module under test uses a `context.Context` to store values, functions are usually
provided to set and retrieve those values.  In addition to returning any value from a context,
the retrieval functions often also return an indicator value which can be used to identify
whether the value was found in the context or not (to differentiate between a non-existent
value and a value that is present with a zero value).

This makes testing the retrieval functions more cumbersome than it might otherwise be:

```go
func TestGetValue(t *testing.T) {
  // ARRANGE
  ctx := context.Background()

  // ACT
  ctx := SomeFuncModifyingContext(ctx, args)

  // ASSERT
  got, ok := GetValue(ctx)
  test.Bool(t, ok).IsTrue()
  test.That(t, got).Equals(value)
}
```

To simplify such tests, two functions are provided:

- `test.ContextIndicator`
- `test.ContextValue`

Both of these function are generic, accpting type parameters `T` and `I` for the value and
indicator types respectively.

In addition to the usual `*testing.T`, these function accept a context to be tested and the
retrieval function; each function returns a testable for the indicator or value returned by
the retrieval function:

```go
func TestGetValue(t *testing.T) {
  // ARRANGE
  ctx := context.WithValue(context.Background(), key, value)

  // ACT & ASSERT
  test.ContextIndicator(t, ctx, GetValue).IsTrue()
  test.ContextValue(t, ctx, GetValue).Equals(value)
}
```

### Test for Expected Panics

Panic tests must be deferred to ensure that the panic is captured and tested.
The `test.ExpectPanic` function returns a `*Panic` value with an `Assert` function
that can be deferred to test for an expected panic.

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  err := errors.New("some error")
  defer test.ExpectPanic(err).Assert(t)

  // ACT
  panic(err)
}
```

The `Assert` function may be called on a `nil` receiver to test that no panic was
recovered, which is useful in table-driven tests:

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  err := errors.New("panicked")

  testcases := []struct {
    name string
    error
    panic  *test.Panic
  }{
    {name: "panic expected", error: err, panic: test.ExpectPanic(err)},
    {name: "no panic expected", err: nil},
  }
  for _, tc := range testcases {
    t.Run(tc.name, func(t *testing.T) {
      // ARRANGE
      defer tc.panic.Assert(t)

      // ACT
      panic(tc.error)
    })
  }
}
```

### Test for an Expected Type

You can test that some value is of an expected type using the `test.Type` function.
This function returns the value as the expected type if the test passes, otherwise it
returns `nil` and the test fails.

If the value is of the expected type, further tests may then be performed on the
returned value as that type:

```go
func TestDoSomething(t *testing.T) {
  // ACT
  result := DoSomething()

  // ASSERT
  if got, ok := test.Type[Customer](t, result); ok {
    // further assertions on got (of type Customer)
  }
}
```
