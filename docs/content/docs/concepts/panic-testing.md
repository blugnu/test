---
title: Panic Testing
weight: 35
---

Panic assertions are different from ordinary value assertions because the panic value is only
available during stack unwinding. In practice this means panic assertions are written as
deferred expectations.

Another way to think about this is that panic assertions are more like "established expectations
that are verified later (deferred)".

## Basic patterns

These patterns are the idiomatic way to test panics in `blugnu/test`.

### Expect any panic

This test will fail if no panic occurs.  It will pass if any panic occurs, regardless of
the recovered value.

```go
func TestPanics(t *testing.T) {
    With(t)

    defer Expect(Panic()).DidOccur()
    panic("boom")
}
```

### Expect a specific panic value

This test will fail if no panic occurs or if the panic value does not match the expected value.
It will pass only if a panic occurs and the expected value is recovered.

```go
func TestPanicsWithValue(t *testing.T) {
    With(t)

    defer Expect(Panic("boom")).DidOccur()
    panic("boom")
}
```

### Expect a specific panic error

This test will fail if no panic occurs or if the recovered panic value is not an `error` that
satisfies `errors.Is` for the expected `error`. It will pass only if a panic occurs and
an `error` is recovered that satisfies `errors.Is` with respect to the expected `error`.

```go
func TestPanicsWithError(t *testing.T) {
    With(t)

    defer Expect(Panic(ErrInvalidArgument)).DidOccur()
    panic(ErrInvalidArgument)
}
```

### Assert that no panic occurs

```go
func TestDoesNotPanic(t *testing.T) {
    With(t)

    defer Expect(Panic()).DidNotOccur()
    // code under test
}
```

## Panic(nil) and NilPanic

> `Panic(nil)` and `NilPanic()` are not interchangeable.  `Panic(nil)` expects no panic,
> while `NilPanic()` expects a panic with a nil value.
>
> Production code should not call `panic(nil)`, but this distinction is important for testing
> legacy code that may do so.

In earlier versions of Go (prior to Go 1.21), `panic(nil)` was indistinguishable from no panic at all,
since the recovered value would be `nil` in both cases. Starting with Go 1.21, an explicit call to
`panic(nil)` is treated as a distinct case and the recovered value is replaced by a sentinel value
that indicates a nil panic.

In `blugnu/test`, `Panic(nil)` is treated as an expectation that no panic occurs; if an explicit
`panic(nil)` is expected, the `NilPanic()` matcher establishes this expectation and matches a panic
with the nil panic sentinel value.

```go
func TestNilPanicValue(t *testing.T) {
    With(t)

    defer Expect(NilPanic()).DidOccur()
    panic(nil)
}
```

## Important constraints

- Use at most one panic expectation per function scope.
- Avoid combining panic assertions with other deferred `recover()` logic in the same scope.
- Keep the deferred expectation close to the statement that may panic.

## See also

- [Matchers]({{< relref "/docs/concepts/matchers" >}})
- [Expectations]({{< relref "/docs/concepts/expectations" >}})
