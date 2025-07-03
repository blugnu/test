<div align="center" style="margin-bottom:20px">
  <img src=".assets/banner.png" alt="logger" />
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

Provides a concise, fluent, type-safe API over the standard testing framework, simplifying common tests
in an extensible fashion whilst maintaining compatibility with the standard testing package.

_Friends don't let friends write tests that are hard to read, hard to maintain, or that
don't fail when they should_.

Don't do this:

```go
func TestDoSomething(t *testing.T) {
  // act
  err := DoSomething()

  // assert
  result, err := DoSomething()
  if err != nil {
    t.Errorf("unexpected error: %v", err)
  }

  expected := 42
  if result != expected {
    t.Errorf("expected result %v, got %v", expected, result)
  }
}
```

Do this instead:

```go
func TestDoSomething(t *testing.T) {
  With(t)

  // act
  result, err := DoSomething()

  // assert
  Expect(err).IsNil()
  Expect(result, "result").To(Equal(42))
}
```

## Features

- **Clean**: No more constantly referencing a `*testing.T`;
- **Concise**: Provides a fluent API that is concise and easy to read, reducing boilerplate code;
- **Type-Safe**: Uses Go's type system to ensure that valid tests are performed, reducing runtime
errors and false positives (or negatives);
- **Compatible**: Compatible with the standard library `testing` package;
- **Matchers**: Provides a rich set of matchers for common assertions, making it easy to express
expectations;
- **Extensible**: Supports custom matchers, enabling functionality to be extended as needed;
- **Panic Testing**: Provides a way to test for expected panics, ensuring that code behaves
correctly under error conditions;
- **Mocking Utilities**: Provides types to assist with implementing mock functions and replacing
  dependencies in tests, allowing for isolated testing of components;
- **Console Recording**: Supports recording console output (`stdout` and `stderr`), to facilitate
testing of log messages and other output;
- **Meta-Testing**: Provides methods for testing a test (used by the package to test itself).

## :construction_worker: &nbsp; Under Construction

_This package is not yet considered stable._

Feel free to use this package if you find it useful, but be aware that the API may change
without notice.

Having said that, the API has just been through a major overhaul to address all of the shortcomings
and annoyances that existed in the previous version, providing a stronger foundation for future
development.

The API will remain as stable as possible, but until the package hits v1.0 this should still be
considered an aspiration, not a commitment.

## :warning: &nbsp; Goroutine IDs

This module uses `runtime.Stack()` to determine the ID of the current goroutine.  This is required
to maintain a reliable per-goroutine stack of `*testing.T` values (a.k.a 'test frames').

_**This mechanism is not guaranteed to be stable and may change in future versions of Go**_.

If you are using this module in a production environment, be aware that changes in future
versions of Go may require break this mechanism for determining a goroutine id, requiring changes. This
may hamper the ability of dependent code to update to a later go version until those changes have been
made.

# :runner: &nbsp; Quick Start

## Installation

```bash
go get github.com/blugnu/test
```

## Read This First

1. [dot-importing the `blugnu/test` package](.assets/readme/dot-import.md)
2. [Test Frames](.assets/readme/test-frames.md)

### Writing a Test: With(t)

The `With()` function is used to set up a test frame for the current test.  This
function is typically called at the start of a test function, passing the `*testing.T` value
from the test function as an argument:

```go
  func TestDoSomething(t *testing.T) {
     With(t) // establishes the initial test frame
     // ...
  }
```

> :bulb: Calling `Parallel(t)` is equivalent to calling `With(t)` followed by `Parallel()` or `t.Parallel()`.

There is no cleanup required after calling `With(t)`; the test frame is automatically cleaned up
when the test completes.

If you use the `blugnu/test` package functions for running table-driven tests or explicit subtests
the test frame stack is managed for you:

```go
  func TestDoSomething(t *testing.T) {
     With(t) // establishes the initial test frame

     Run("subtest", func() {
        // no need to call With(t) here; it is managed automatically
        // ...
     })
  }
```

If a new test frame is created outside of the `test` package, then the `With(t)` function
must be called again to push that test frame onto the stack.  For example, if you choose to create
a subtest using `testing.T.Run()` and want to use the `blugnu/test` functions in that subtest:

```go
  func TestDoSomething(t *testing.T) {
     With(t) // establishes the initial test frame

     // using the testing.T.Run method...
     t.Run("subtest", func(t *testing.T) {
        With(t) // so a new test frame must be established
        // ...
     })
  }
```

Generally speaking it is much easier to use the `blugnu/test` package functions to avoid having
to use `With(t)` (or `Parallel(t)`) for anything other than establishing the initial test frame
for each test function.

> :warning: Neither `With(t)` nor `Parallel(t)` should be called unless required
> to establish a test frame for a new `*testing.T` value.

### Writing a Test: Expect

Almost all tests are written using `Expect` to create an expectation
over some value (the _subject_).  `Expect` returns an _expectation_ with
methods for testing the subject.

Some expectation methods test the value directly, such as `IsEmpty()`, `IsNil()` and `IsNotNil()`:

```go
  err := DoSomething()
  Expect(err).IsNil()
```

### Using Matchers

The `To` method of an expectation delegates the evaluation of a test
to a type-safe matcher, usually provided by a factory function where the
factory function itself is named to fluently express the expected outcome,
e.g.:

```go
  Expect(got).To(Equal(expected))
```

In this example, the `Equal()` function is a factory function that
returns an `equal.Matcher`, used to test that the subject is equal to
some expected value.

The `Should` method provides the same functionality as `To()`, but
accepts matchers that are not type-safe, referred to as _any-matchers_ as
they accept `any` as the subject type.  This is necessary for matchers which
do not accept any arguments or where the compatible arguments cannot be
expressed as a generic type constraint.

An example of an any-matcher is `BeEmpty()`, which can be used to
test whether a slice, map, channel or string is empty:

```go
  Expect(got).Should(BeEmpty())
```

> Further information on any-matchers is provided in the section on [Type-Safety: Any-Matchers](#type-safety-any-matchers)

### Type-Safety: Matcher Compatibility

`Expect()` is a generic function, where the type `T` is inferred from
the subject value; the `To()` function will only accept matchers that
are compatible with the type of the subject value.

For example, in the previous example, the `equal.Matcher` uses the `==` operator
to determine equality, so is constrained to types that satisfy `comparable`. As
a result, values of non-comparable type cannot be tested using this matcher:

```go
  Expect(got).To(Equal([]byte("expected result"))) // ERROR: cannot use `Equal` with `[]byte`
```

In this case, two alternatives exist:

1. `DeepEqual()` returns a `equal.DeepMatcher` that may be used with _any_ type, using `reflect.DeepEqual` for equality;

2. `EqualBytes` returns a `bytes.EqualMatcher` which provides test failure reports that are specific to `[]byte` values.

```go
  // DeepEqual can be used with any type; uses reflect.DeepEqual for equality
  // but can result in verbose failure reports since both expected and got
  // values are printed in full
  Expect(got).To(DeepEqual([]byte("expected result")))

  // EqualBytes is a type-safe matcher specifically for []byte values
  // providing failure reports that report and highlight differences between
  // byte slices accurately and concisely
  Expect(got).To(EqualBytes([]byte("expected result")))
```

### Type-Safety: Any-Matchers

Not all matchers are constrained by types; some matchers accept `any` as the
subject type, allowing them to be used with any value, referred to as
"_any-matchers_".

There are two main use cases for any-matchers:

- the matcher supports testing values of a variety of types that cannot be
described in a generic type constraint

> **Why?**: if a matcher is designed to work with a set of types that
> cannot be described in a generic type constraint, it must accept `any` as the
> subject type.
>
> For example, the `BeEmpty()` matcher can be used to test whether a slice, map,
> channel or string is empty (and more), but there is no way to express that
> set of types in a type constraint.

- the type of the expected value is _implicit_ in the test, rather than
explicit

> **Why?**: without an explicit expected value, it is not possible to infer the
> corresponding subject type; the matcher must accept `any` as the subject type

Whilst any-matchers _can_ be used with `To()`, by casting the subject to `any`,
this can be cumbersome and disrupts the fluency of the test:

```go
  Expect(any(got)).To(BeEmpty())
```

As an alternative, the `Should()` and `ShouldNot()` methods provide the same
functionality as `To()`/`ToNot()`, but accepting any-matchers rather than
type-safe matchers:

```go
  Expect(got).Should(BeEmpty())
  Expect(got).ShouldNot(BeEmpty())
```

### Type-Safety: Invalid Tests fail as Invalid

If an expectation method is called inappropriately on a subject, the
test will often fail as an invalid test.  For example, if the `IsNil()`
method is called on a value of a type that does not support a meaningful
`nil` value the test will fail, not because the value is not `nil` but
because the test itself is invalid:

```go
  Expect(42).IsNil() // <== INVALID TEST: `int` is not nilable
```

In general, expectation tests will attempt to provide a meaningful test
consistent with the intent, only failing as invalid if a meaningful test
is not possible.

# Guides

## Basic Usage

- [Setting Expectations](#setting-expectations)
- [Short-Circuit Evaluation](#short-circuit-evaluation)
- [Testing for Nil/Not Nil](#testing-nilnot-nil)
- [Testing Errors](#testing-errors)
- [Testing for Panics](#testing-for-panics)
  - [Panic(nil) vs NilPanic()](#panicnil-vs-nilpanic)
- [Testing for Emptiness](#testing-emptiness)
- [Testing With Matchers](#testing-with-matchers)
  - [Matcher Options](#matcher-options)
  - [Custom Matchers](#custom-matchers)
- [Testing Maps](#testing-maps)
- [Testing Slices](#testing-slices)
- [Testing Context](#testing-context)

## Advanced Usage

In addition to performing common, basic tests, the `test` package also provides support for more advanced testing scenarios:

| Category | Description |
| --- | --- |
| [Mocking Functions](#mocking-functions) | mock functions for testing |
| [Recording Console Output](#recording-console-output) | record output of a function that writes to `stdout` and/or `stderr` |
| [Test for an Expected Type](#test-for-an-expected-type) | test that a value is of an expected type |
| [Testing a Test](#testing-a-test) | test your own test helper functions |

------
</br>

# Setting Expectations

Almost all tests start with setting an expectation over some value (the _subject_).

The `Expect` function returns an _expectation_ with methods for testing the subject:

```go
  Expect(got)  // returns an expectation for the value `got`
```

In addition to a subject, the `Expect()` function accepts options to configure the
expectation, passed as variadic arguments.  Currently the only option is a name for the
subject, which is used in test failure reports to identify the subject being tested:

```go
  Expect(result, "result")  // returns an expectation named "result"
```

An expectation alone does not perform any tests; it simply provides a way to
express an expectation over a value.  The expectation is evaluated when a test method
is called on the expectation, such as `IsNil()`, `IsNotNil()`, `IsEmpty()`, `To()` or
`DidOccur()`.

# Short-Circuit Evaluation

If an expectation is critical to the test, it can be useful to short-circuit the test
execution if the expectation fails.  For example, if a value is expected to not be
`nil` and further tests on that value will panic or be guaranteed to fail:

```go
  Expect(value).ShouldNot(BeNil())
  Expect(value.Name).To(Equal("some name"))  // this test will panic if value is nil
```

To short-circuit the test execution if the expectation fails, the `opt.IsRequired(true)`
option can be passed to the expectation method:

```go
  Expect(value).ShouldNot(BeNil(), opt.IsRequired(true))
  Expect(value.Name).To(Equal("some name"))  // this test will not be executed if value is nil
```

Alternatively, the `Require()` function can be used to create the expectation:

```go
  Require(value).ShouldNot(BeNil())
  Expect(value.Name).To(Equal("some name"))  // this test will not be executed if value is nil
```

In both cases, if the expectation fails the current test exits without evaluating any further
expectations. Execution continues with the next test.

# Testing Nil/Not Nil

A nilness matcher is provided which may be used with the `Should()` or `ShouldNot()`
methods:

```go
  Expect(value).Should(BeNil())      // fails if value is not nil or of a non-nilable type
  Expect(value).ShouldNot(BeNil())   // fails if value is nil
```

Since these tests are common, `IsNil()` and `IsNotNil()` convenience methods are also
provided on expectations:

```go
    Expect(value).IsNil()      // fails if value is not nil or of a non-nilable type
    Expect(value).IsNotNil()   // fails if value is nil
```

> :bulb: _If `IsNil()`/`Should(BeNil())` is used on a subject of a type that does not have a
> meaningful `nil` value, the test will fail as invalid_.
>
> _Types that may be tested for `nil` are: `chan`, `func`, `interface`, `slice`, `map`, and
> `pointer`_.
>
> `IsNotNil()`/`ShouldNot(BeNil())` will **NOT** fail on a non-nilable subject.

```go
  var got 42
  Expect(got).IsNil()      // <== INVALID TEST: `int` is not nilable
  Expect(got).IsNotNil()   // <== VALID TEST: `int` is not nilable, so this test passes
```

# Testing Errors

## Testing that an Error did not occur

There are two ways to explicitly test that an error did not occur:

```go
  Expect(err).DidNotOccur()
  Expect(err).IsNil()
```

A third way to test that an error did not occur is to use the `Is()` method, passing `nil`
as the expected error:

```go
  Expect(err).Is(nil)
```

This is most useful when testing an error in a table driven test where each test case
may have an expected error or `nil`:

```go
  Expect(err).Is(tc.err)  // tc.err may be nil or an expected error
```

## Testing that an Error occurred (any error)

```go
  Expect(err).DidOccur()
  Expect(err).IsNotNil()
```

## Testing that a Specific Error Occurred

```go
  Expect(err).Is(expectedError) // passes if `errors.Is(err, expectedError)` is true
```

> _If `nil` is passed as the expected error, the test is equivalent to `IsNil()`_.

# Testing for Panics

Panics can be tested to ensure that an expected panic did (or did not) happen. Since
panic tests rely on the recovery mechanism in Go, they must be deferred to ensure
that the panic is captured and tested correctly.

> :warning: There must be at most **ONE** panic test per function; multiple
> panic tests (or other calls to `recover()` in general) in the same function
> will not work as expected.

When testing panics, recovered values may be significant or they may be ignored.

For example, if testing only that a panic occurred without caring about the
recovered value:

```go
  defer Expect(Panic()).DidOccur()
```

By contrast, if the recovered value is significant, it can be tested by specifying
the expected panic value.  The following tests will pass if a panic occurs and the
recovered value is equal to the expected string:

```go
  defer Expect(Panic("expected panic")).DidOccur()
```

> :warning: `Panic(nil)` is a special case.  see: [Panic(nil) vs NilPanic()](#panicnil-vs-nilpanic)

## Testing a Panic with a Recovered Error

When testing for a panic that recovers an error and the expected recovered value
is specified, the test will pass if the recovered value is an error and
`errors.Is(recovered, expectedErr)` is true:

```go
  defer Expect(Panic(expectedErr)).DidOccur()
```

## Testing that a Panic did NOT occur

It is also possible to explicitly test that a panic did not occur:

```go
  defer Expect(Panic()).DidNotOccur()
```

> :bulb: `Expect(Panic(nil)).DidOccur()` is a special case that is exactly equivalent to
> the above test for no panic. See: [Panic(nil) vs NilPanic()](#panicnil-vs-nilpanic)

Again, if the recovered value is significant, it can be tested by specifying
the expected panic value.

> :warning: If a value is recovered from a panic that is different to that
> expected, the test will fail as an `unexpected panic`.

## Panic(nil) vs NilPanic()

Prior to go 1.21, `recover()` could return `nil` if a `panic(nil)` had been called,
making it impossible to distinguish from no panic having occurred.

From go 1.21 onwards, a `panic(nil)` call is now transformed such that a specific
runtime error will be recovered.

`Panic(nil)` is treated as a special case that is used to test that a panic did NOT occur.
i.e. the following are exactly equivalent:

```go
  // test that a panic did NOT occur
  defer Expect(Panic(nil)).DidNotOccur()
```

and

```go
  // also test that a panic did NOT occur
  defer Expect(Panic(nil)).DidOccur()
```

This may seem counter-intuitive, but there is a good reason for this.

The motivation is to simplify table-driven tests where each test case may
expect a panic or not.  Without this special case, the test would
require a conditional to determine whether to test for a panic or not:

```go
  if tc.expectPanic {
    defer Expect(Panic(tc.expectedPanic)).DidOccur()
  } else {
    defer Expect(Panic()).DidNotOccur()
  }
```

With the special case of `Panic(nil)`, the test can be simplified to:

```go
  defer Expect(Panic(tc.expectedPanic)).DidOccur()
```

Where `tc.expectedPanic` may be `nil` (panic not expected to occur) or an
expected value to be recovered from a panic.

In the unlikely event that you need to test specifically for a `panic(nil)`
having occured, the go 1.21+ runtime error can be tested for:

```go
  defer Expect(Panic(&runtime.PanicNilError{})).DidOccur()
```

To make even this unlikely case easier, the `NilPanic()` function
is provided, so the above can be simplified to:

```go
  defer Expect(NilPanic()).DidOccur()
```

### Testing Emptiness

The `BeEmpty()` and `BeEmptyOrNil()` matchers are provided to test whether a
value is considered empty, or not.  These are any-matchers for use with the
`Should()` or `ShouldNot()` methods.

```go
  Expect(value).Should(BeEmpty())         // fails if value is not empty or nil
  Expect(value).Should(BeEmptyOrNil())    // fails if value is not empty and not nil
  Expect(value).ShouldNot(BeEmpty())      // fails if value is empty
```

`BeEmpty()` and `BeEmptyOrNil()` are provided to differentiate between empty
and `nil` values where useful.

For example, if testing a slice, `IsEmpty()` will pass if the slice is empty
but will fail if the slice is `nil`, while `IsEmptyOrNil()` will pass in both cases.

Emptiness tests will fail as invalid if emptiness of the value cannot be determined.

Emptiness is defined as follows:

- for `string`, `slice`, `map`, `chan` and `array` types, emptiness is defined as
  `len(value) == 0`
- for all other types, emptiness is determine by the implementation of a `Count()`,
  `Len()` or `Length()` method returning 0 (zero) of type  `int`, `int64`, `uint`
  or `uint64`
- if testing a value for this emptiness cannot be determined, the test will fail
  as invalid.

# Testing With Matchers

The `To()`, `ToNot()`, `Should()` and `ShouldNot`() methods delegate the evaluation
of a test to a matcher, usually provided by a factory function.  Matcher factory
functions are typically named to describe the expected outcome in a fluent fashion
as part of the test expression, e.g.:

```go
  Expect(got).To(Equal(expected))  // uses an equal.Matcher{} from the 'blugnu/test/matchers/equal' package
```

> _**"Matching" Methods**_: For brevity, the `To()`, `ToNot()`, `Should()` and
> `ShouldNot()` methods are referred to generically as _Matching Methods_.

## Type-Safe Matchers

A type-safe matcher is a matcher that is compatible with a specific type of subject
value.  A type-safe matcher may be constrained to a single, explicit formal type,
or it may be a generic matcher where type compatability is expressed through the
constraints on the generic type parameter.

For example:

- `HasContextKey()`: is explicitly compatible only with `context.Context` values
- `Equal()`: uses a generic matcher that is compatible with any type `T` that
  satisfies the `comparable` constraint

By contrast:

- `BeEmpty()` is **NOT** type-safe: _it is compatible with (literally) any type_

## Any-Matchers

An any-matcher is a matcher that accepts `any` as the subject type, allowing it to be used
with literally any value.  Any-matchers are used with the `Should()` or `ShouldNot()` matching
methods.

> Any-matchers may also be used with the `To()` or `ToNot()` matching methods
> if the formal type of the subject is `any`, but this is not recommended.

# Built-In Matchers

A number of matchers are provided in the `test` package, including:

<!-- markdownlint-disable MD013 -->
| Factory Function | Subject Type | Description |
| --- | --- | --- |
| `BeEmpty()` | `any` | Tests that the subject is empty but not nil |
| `BeEmptyOrNil()` | `any` | Tests that the subject is empty or nil |
| `BeGreaterThan(T)` | `T cmp.Ordered` | Tests that the subject is greater than the expected value using the `>` operator |
| `BeLessThan(T)` | `T cmp.Ordered` | Tests that the subject is less than the expected value using the `<` operator |
| `BeNil()` | `any` | Tests that the subject is nil |
| `Equal(T)` | `T comparable` | Tests that the subject is equal to the expected value using the `==` operator |
| `DeepEqual(T)` | `T any` | Tests that the subject is deeply equal to the expected value using `reflect.DeepEqual` |
| `EqualBytes([]byte)` | `[]byte` | Tests that `[]byte` slices are equal, with detailed failure report highlighting different bytes |
| `EqualMap(map[K,V])` | `map[K,V]` | Tests that the subject is equal to the expected map |
| `ContainItem(T)` | `[]T` | Tests that the subject contains an expected item |
| `ContainItems([]T)` | `[]T` | Tests that the subject contain the expected items (in any order, not necessarily contiguously) |
| `ContainMap(map[K,V])` | `map[K,V]` | Tests that the subject contains the expected map (keys and values must match) |
| `ContainMapEntry(K,V)` | `map[K,V]` | Tests that the subject contains the expected map entry |
| `ContainSlice([]T)` | `[]T` | Tests that the subject contains the expected slice (items must be present contiguously and in order) |
| `ContainString(expected T)` | `T ~string` | Tests that the subject contains an expected substring |
| `HaveContextKey(K)` | `context.Context` | Tests that the context contains the expected key |
| `HaveContextValue(K,V)` | `context.Context` | Tests that the context contains the expected key and value |
<!-- markdownlint-enable -->

Matchers are used by passing the matcher to one of th expectation matching methods together
with options to control the behaviour of the expectation or the matcher itself.

A matcher is typically constructed by a factory function accepting any arguments required
by the matcher.  It is worth repeating that _options_ supported by the matcher are passed
as arguments to the matching method, _not_ the matcher factory:

```go
  // the opt.OnFailure option replaces the default error report
  // with the custom "failed!" message
  Expect(got).To(Equal(expected),
    opt.OnFailure("failed!"),
  )

  // override the use of the `==` operator with a custom comparison function
  // where the subject type is a hypothetical `MyStruct` type
  Expect(got).ToNot(Equal(expected),
    func(exp, got MyStruct) bool { return /* custom comparison logic */ },
  )
```

## Matcher Options

Matching methods accept options as variadic arguments following the matcher.

The matching methods themselves support options for customising the test error report
in the event of failure.

```go
  Expect(got).To(Equal(expected),
    opt.OnFailure(fmt.Sprintf("expected %v, got %v", expected, got)),
  )
```

Options supported by matching methods (and therefore _all_ matchers) include:

<!-- markdownlint-disable MD013 -->
| Option | Description |
| --- | --- |
| `opt.FailureReport(func)` | a function that returns a custom error report for the test failure; the function must be of type `func(...any) []string` |
| `opt.OnFailure(string)`   | a string to use as the error report for the test failure; this overrides the default error report for the matcher |
| `opt.AsDeclaration(bool)` | a boolean to indicate whether values (other than strings) in test failure reports should be formatted as declarations (`%#v` rather than `%v`) |
| `opt.QuotedStrings(bool)` | a boolean to indicate whether string values should be quoted in failure reports; defaults to `true` |
| `opt.IsRequired(bool)`    | a boolean to indicate whether the expectation is required; defaults to `false` |
<!-- markdownlint-enable -->

> `opt.OnFailure()` is a convenience function that returns an `opt.FailureReport` with a
> function that returns the specified string in the report.
>
> `opt.FailureReport` and `opt.OnFailure()` are mutually exclusive; if both are specified, only the
> first in the options list will be used.

The `...any` argument to an `opt.FailureReport` function is used to pass any options supplied to the
matcher, so that the error report can respect those options where appropriate.

> See the [Custom Failure Report Guide](.assets/readme/custom-failure-reports.md) for details.

Matchers may support options to modify their behaviour.  The specific options supported
by a matcher are documented on the relevant matcher factory function.

Examples of other options supported by some matchers include:

<!-- markdownlint-disable MD013 -->
| Option | Description | Default |
| --- | --- | --- |
| `opt.ExactOrder(bool)` | a boolean to indicate whether the order of items in a collection is significant | `false` |
| `opt.CaseSensitive(bool)` | a boolean to indicate whether string comparisons should be case-insensitive | `true` |
| `opt.AsDeclaration(bool)` | a boolean to indicate whether values other than strings should be formatted as declarations (`%#v` vs `%v`) | `false` |
| `opt.QuotedStrings(bool)` | a boolean to indicate whether string values should be quoted in failure reports | `true` |
| `func(T, T) bool` | a type-safe custom comparison function; the type `T` is the type of the subject value |  |
| `func(any, any) bool` | a custom comparison function accepting `expected` and `subject` values as `any` |  |
<!-- markdownlint-enable -->

> type-safe custom comparison functions are preferred over `any` comparisons.  Only one
> should be specified; if multiple comparison functions are specified, the first type-safe
> function will be used in preference over the first `any` function.

## Custom Matchers

Custom matchers may be implemented by defining a type that implements a `Match(T, ...any) bool` method.

`T` may be:

- an explicit, formal type
- a generic type parameter with constraints
- `any` if the matcher is not type-safe

Refer to the [Custom Matchers Implementation Guide](.assets/readme/custom-matchers.md) for details.

------
</br>

# Testing Maps

The `test` package provides matchers for testing maps, including the ability to test for
equality, containment of items or the existence of a specific key:value entry.

<!-- markdownlint-disable MD013 -->
| Matcher | Subject Type | Description |
| --- | --- | --- |
| `EqualMap(map[K,V])` | `map[K,V]` | tests that the subject is equal to the expected map |
| `ContainMap(map[K,V])` | `map[K,V]` | tests that the subject contains the expected map (keys and values must match, order is not significant) |
| `ContainMapEntry(K,V)` | `map[K,V]` | tests that the subject contains the expected map entry (key and value must match) |
<!-- markdownlint-enable -->

These matchers accept either a map or a pair of key and value parameters; type inference
ensures type-compatibility with expectation subjects that are maps.

## Testing Keys or Values in Isolation

When testing keys or values in isolation (a key without a value or vice versa), any matcher
would not have enough information to determine the type of both key and value to provide
compatibility with a map subject without explicitly instantiating with declared type
information.

To avoid this, functions are provided to extract keys or values as slices, enabling the use of
slice matchers in fluent fashion, e.g.:

```go
  Expect(KeysOfMap(m)).To(ContainItem("key"))     // tests that the map contains the key "key"
  Expect(ValuesOfMap(m)).To(ContainItem("value")) // tests that the map contains a key having the value "value"
```

------
</br>

# Testing Slices

The `test` package provides matchers for testing slices, including the ability to test for
equality, containment of items or the existence of a specific item in a slice.

<!-- markdownlint-disable MD013 -->
| Matcher | Subject Type | Description |
| --- | --- | --- |
| `EqualSlice([]T)` | `[]T` | tests that the subject is equal to the expected slice (items must be present contiguously and in order); options supported include `opt.ExactOrder(false)` or `opt.AnyOrder()`, to allow equality of slices where order is not significant |
| `ContainItem(T)` | `[]T` | tests that the subject contains an expected item |
| `ContainItems([]T)` | `[]T` | tests that the subject contains the expected items (in any order, not necessarily contiguously) |
| `ContainSlice([]T)` | `[]T` | tests that the subject contains the expected slice (items must be present contiguously and in order); options supported include `opt.ExactOrder(false)` or `opt.AnyOrder()`, to allow containment of slices where order is not significant |
<!-- markdownlint-enable -->

To test for an expected length of a slice, use `len(slice)` as the subject:

```go
  Expect(len(slice)).To(Equal(expectedLength)) // tests that the slice has the expected length
```

For emptiness tests where precise length is not significant (other than zero, for emptiness), the `IsEmpty()`
and `IsNotEmpty()` methods can be used on any slice subject:

```go
  Expect(slice).IsEmpty()      // tests that the slice is empty
  Expect(slice).IsNotEmpty()   // tests that the slice is not empty
```

------
</br>

# Testing Context

Passing values in a `context.Context` is a common pattern in Go, and the `test` package
provides a way to test that a context contains the expected values.

The `HaveContextKey(K)` and `HaveContextValue(K,V)` matchers can be used to test that a
context contains a specific key or key-value pair.

```go
  Expect(ctx).To(HaveContextKey("key"))            // tests that the context contains the key "key"
  Expect(ctx).To(HaveContextValue("key", "value")) // tests that the context contains the key "key" with value "value"
```

The type of the key is determined by the type parameter `K` of the matcher, which
must be the type of the key used in the context, not just compatible.

For example, if a custom `string` type has been used for the key, the value of any key
expected by the matcher must be cast to that type:

```go
  // e.g. where MyPackageContextKey is the type used for context keys:
  Expect(ctx).To(HaveContextKey(MyPackageContextKey("key")))   
```

------
</br>

# Mocking Functions

When testing functions that call other functions, it is often necessary to mock the functions
being called to ensure that the tests are isolated and that the functions being tested are
not dependent on the behaviour of the functions being called.

The `test.MockFn[A, R]` type provides a way to mock a function accepting arguments of type `A` and
returning a result of type `R`.

All `test.MockFn` values support an optional `error` value which may simply be ignored/not used
if the mocked function does not return an error.

If the function being mocked does not return any value other than an `error`, the result type `R`
should be `any` and ignored.  Similarly if the function being mocked does not require any
arguments, the argument type `A` should be `any` and ignored.

## Fake Function Results

The `test.MockFn` type can provide fake results for a mocked function.  Fake results may be setup
in two different ways:

- expected calls mode.  `ExpectCall()` is used configure an expected call; this returns a value
  with a `WillReturn` method to setup a result to be returned for that call.  In this mode, calls
  to the mocked function that do not match the expected calls will cause the test to fail.

- mapped result more.  `WhenCalledWith(args A)` is used to setup a result to be returned when the
  mocked function is called with the specified arguments.  In this mode, calls to the mocked
  function that do not match any of the mapped results will cause the test to fail.

## Multiple Arguments/Result Values

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

------

# Recording Console Output

The `test.Record` function records the output of a function that writes to `stdout`
and/or `stderr`, returning the output as a pair of `[]string` values for `stdout`
and `stderr` respectively.

> :bulb: The function does not return an error; it will panic if the redirection
> fails. This is an intentional design choice to ensure that a test fails if the
> mechanism is not working correctly, avoiding incorrect results without requiring
> a test to handle any error.

Since failed tests will write to `stdout`, the output will include any test
failures that occur during execution of the captured function. You may wish to
structure your code to perform tests outside of the recorded functions to avoid
this and simplify testing of the output:

```go
func TestDoSomething(t *testing.T) {
  With(t)

  // ACT
  var (
    err error
    result string
  )
  stdout, stderr := Record(func () {
    result, err := DoSomething()
    return err
  })

  // ASSERT
  Expect(err).IsNil()
  Expect(result).To(Equal("foo"))
  Expect(stdout).To(ContainItem("DoSomething wrote this to stdout"))
}
```

------

# Testing a Test Helper

If you write your own test helpers (or matchers), you should of course test them.
A `TestHelper()` function is provided to enable you to do just that.

`TestHelper()` accepts a function that executes your test helper; it is performed
using a separate test runner, independent of the current test.  This allows the
helper being tested to fail without affecting the outcome of the test that is
testing it.

The `TestHelper()` function returns a `test.R` value that contains information about
the outcome of the test.  You could test the information in this `R` value directly,
but the `R` type provides methods to make this easier.

## Testing the Outcome

To test the outcome of a test, without considering any output, you can pass the expected
outcome as an argument to the `R.Expect()` method:

```go
  result := TestHelper(func() {
    /* your test code here */
  })

  result.Expect(TestPassed)
```

## Testing Test Helper Output

It is recommended to test the output of your test helper or matcher when it fails. You
can do this by passing the expected lines of test output as strings to the `R.Expect()`
method:

```go
  result := Test(func() {
    /* your test code here */
  })

  result.Expect(
    "expected output line 1",
    "expected output line 2",
  )
```

> :bulb: By testing the output, the test is implicitly expected to fail, so the
> `R.Expect()` method in this case will also test that the outcome is `TestFailed`.

## Running Multiple Test Scenarios

A specialised version of `RunScenarios()` is provided to test a test helper or
custom matcher: `RunTestScenarios()`. This accepts a slice of `TestScenario` values,
where each scenario is a test case to be run against your test helper or matcher.

`RunTestScenarios()` implements a test runner function for you, so all you need
to do is provide a slice of scenarios, with each scenario consisting of:

- `Scenario`: a name for the scenario (scenario); each scenario is run in a subtest
   using this name;
- `Act`: a function that contains the test code for the scenario; this function has
   the signature `func()`;
- `Assert`: a function that tests the test outcome; this function has the signature
  `func(*R)` where `R` is the result of the test scenario.

The `Assert` function is optional; if not provided the scenario is one where the
test is expected to pass without any errors or failures.

### Debugging and Skipping Scenarios

When you have a large number of scenarios it is sometimes useful to focus on a
subset, or specific test, or to ignore scenarios that are not yet implemented or
known to be failing with problems which you wish to ignore while focussing on
other scenarios.

The `RunTestScenarios()` function and `TestScenario` type support this by
providing a `Skip` and `Debug` field on each `TestScenario`.

> :warning: &nbsp; When setting either `Debug` or `Skip` to `true`, it is important
to remember to remove those settings when you are ready to move on to the next
focus of your testing.

```go
  scenarios := []TestScenario{
    {
      Scenario: "test scenario 1",
      Act: func() { /* test code */ },
      Assert: func(r *R) { /* assertions */ },
    },
    {
      Scenario: "test scenario 2",
      Skip: true, //                              <== this scenario won't run
      Act: func() { /* test code */ },
      Assert: func(r *R) { /* assertions */ },
    },
  }
```

Setting `Skip` to `true` may be impractical if you have a large number of scenarios
and want to run only a few of them. In this case, you can use the `Debug` field to
focus on a single scenario or a subset of scenarios.

## Debugging Scenarios

> :bulb: &nbsp; Setting `Debug` does not invoke the debugger or subject a test scenario to
> any special treatment, beyond selectively running it.  The name merely reflects that it
> most likely to be of use when debugging.

When `Debug` is set to `true` on any one or more scenarios, the test runner will run
ONLY those scenarios, skipping all other scenarios:

```go
  scenarios := []TestScenario{
    {
      Scenario: "test scenario 1",
      Debug: true, //                              <== only this scenario will run
      Act: func() { /* test code */ },
      Assert: func(r *R) { /* assertions */ },
    },
    {
      Scenario: "test scenario 2",
      Act: func() { /* test code */ },
      Assert: func(r *R) { /* assertions */ },
    },
    {
      Scenario: "test scenario 3",
      Act: func() { /* test code */ },
      Assert: func(r *R) { /* assertions */ },
    },
  }
```

### :warning: &nbsp; If both `Debug` and `Skip` are set `true` the scenario is skipped

------

# Test for an Expected Type

You can test that some value is of an expected type using the `ExpectType` function.

This function returns the value as the expected type and `true` if the test passes;
otherwise the zero-value of the expected type is returned, with `false`.

A common pattern when this type of test is useful is to assert the type of some
value and then perform additional tests on that value appropriate to the type:

```go
func TestDoSomething(t *testing.T) {
  With(t)

  // ACT
  result := DoSomething()

  // ASSERT
  if cust, ok := ExpectType[Customer](result); ok {
    // further assertions on cust (type: Customer) ...
  }
}
```

This test can only be used to test that a value is of a type that can be expressed
through the type parameter on the generic function.

For example, the following test will fail as an invalid test:

```go
  type Counter interface {
    Count() int
  }
  ExpectType[Counter](result) // INVALID TEST: cannot be used to test for interfaces
```
