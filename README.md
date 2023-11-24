<div align="center" style="margin-bottom:20px">
  <!-- <img src=".assets/banner.png" alt="logger" /> -->
  <div align="center">
    <a href="https://github.com/blugnu/test/actions/workflows/pipeline.yml"><img alt="build-status" src="https://github.com/blugnu/test/actions/workflows/pipeline.yml/badge.svg?branch=master&style=flat-square"/></a>
    <a href="https://goreportcard.com/report/github.com/blugnu/test" ><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/test"/></a>
    <a><img alt="go version >= 1.20" src="https://img.shields.io/github/go-mod/go-version/blugnu/test?style=flat-square"/></a>
    <a href="https://github.com/blugnu/test/blob/master/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/blugnu/test?color=%234275f5&style=flat-square"/></a>
    <a href="https://coveralls.io/github/blugnu/magpack?branch=master"><img alt="coverage" src="https://img.shields.io/coveralls/github/blugnu/test?style=flat-square"/></a>
    <a href="https://pkg.go.dev/github.com/blugnu/test"><img alt="docs" src="https://pkg.go.dev/badge/github.com/blugnu/test"/></a>
    <hr/>
  </div>
</div>

<br>

# blugnu/test

Provides some simple test helpers for use with the standard library testing package.  It is not a replacement for the testing package or complete testing framework.

## Features

- [x] Test comparable values (e.g. `[]byte`, `map[]`)
- [x] Capture console output
- [x] Test console output
- [x] Test test helpers
- [x] Test for expected panics

## Installation

```bash
go get github.com/blugnu/test
```

## Examples

- [Test for Unexpected Errors and Equality of Comparable Values](#test-for-unexpected-errors-and-equality-of-comparable-values)
- [Testing `map` and `[]byte`](#testing-maps-and-byte-slices)
- [Capture and Test Console Output](#capture-and-test-console-output)
- [Testing a Test Helper](#testing-a-test-helper)
- [Test for Expected Panics](#test-for-expected-panics)
- [Test for an Expected Type](#test-for-an-expected-type)


### Test for Unexpected Errors and Equality of Comparable Values

If a test is not expected to return an error you can use the `test.UnexpectedError` function to test for this.  Similarly, if a test is expected to return a specific value you can use the `test.Equals` function to test for this.  The `test.Equals` function can be used to test for equality of any comparable value.

```go
func TestDoSomething(t *testing.T) {
  // ACT
  got, err := DoSomething()

  // ASSERT
  test.UnexpectedError(t, err)
  test.Equal(t, "foo", got)
}
```

To test for a specific error you can use the `test.ErrorIs` function:

```go
func TestDoSomething(t *testing.T) {
  // ACT
  err := DoSomething()

  // ASSERT
  test.ErrorIs(t, ErrSomething, err)
}
```

An optional argument may be used to specify the format of the expected and actual values in any test failure report produced by `test.Equal()`.  The default format is `FormatDefault`:

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := "foo"
  got := "bar"

  // ACT & ASSERT
  test.Equal(t, expected, got)               // displays values in a failure as default (%v)
  test.Equal(t, expected, got, FormatHex)    // displays values in a failure as hexadecimal
}
```

Any format may be specified by casting a string as a `Format` if needed; sensible values are provided as constants.

### Testing Maps and Byte Slices

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := map[string]string{
    "foo": "bar",
  }

  // ACT
  got := DoSomething()

  // ASSERT
  test.Maps(t, expected, got)
}
```

When testing `[]byte`, an optional format argument may be used to specify the format of the expected and actual values in any test failure report.  The default format is `BytesHex` (hexadecimal):

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  expected := []byte("foo")
  got := []byte("bar")

  // ACT & ASSERT
  test.Bytes(t, expected, got)               // displays values in a failure as hexadecimal
  test.Bytes(t, expected, got, BytesBinary)  // displays values in a failure as binary
}
```

Any format may be specified by casting a string as a `BytesFormat` if needed; sensible values are provided as constants.


### Capture and Test Console Output

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  var err error

  // ACT - discards stderr and result of DoSomething
  stdout, _ := test.CaptureOutput(t, func (*testing.T) {
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

The `test.Helper` function combines the execution of a test helper with the testing of the outcome of the helper.  The outcome of the helper is specific using `test.ShouldPass` or `test.ShouldFail` or providing a `test.*Panic` if the helper is expected to panic.

The output of the test helper is returned as `CapturedOutput` (both `stdout` and `stderr`) so that the presentation of test failure messages in the log can also be tested and verified. 

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

> **NOTE:** _It is important that the helper function being tested is called with the `*testing.T` passed to the function that runs it (`st` in the example above) and not the `T` of the test (`t` in the example)._

### Test for Expected Panics

Panic tests must be deferred to ensure that the panic is captured and tested.  The `test.ExpectPanic` function returns a `*Panic` value with an `IsRecovered` function that can be deferred to test for an expected panic.

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  err := errors.New("some error")
  defer test.ExpectPanic(err).IsRecovered(t)

  // ACT
  panic(err)
}
```

The `IsRecovered` function may be called on a `nil` receiver to test that no panic was recovered, which is useful in table-driven tests:

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
      defer tc.panic.IsRecovered(t)

      // ACT
      panic(tc.error)
    })
  }
}
```

## Test for an Expected Type

You can test that some value is of an expected type using the `test.Type` function.  This function returns the value as the expected type if the test passes, otherwise it returns `nil` and the test fails. If the value is of the expected type, further tests may then be performed on the returned value as that type:

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
