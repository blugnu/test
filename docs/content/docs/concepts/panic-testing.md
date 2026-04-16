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

> `Panic(nil)` and `NilPanic()` are not interchangeable.  `Panic(nil)` will match when no panic
> occurs, while `NilPanic()` will match when a panic occurs with a `nil` value.

Before Go 1.21, `panic(nil)` was indistinguishable from no panic at all, since the recovered value
would be `nil` in both cases. Starting with Go 1.21, an explicit call to `panic(nil)` is intercepted
by the runtime which replaces the `nil` recovery value with a sentinel value (`NilPanic`).

In `blugnu/test`, `Panic(nil)` is treated as an expectation that no panic occurs; in the unlikely
that a test wishes to match an explicit `panic(nil)`, the `NilPanic()` matcher must be used.

This allows table-driven tests to specify an expected panic recovery value of `nil` to indicate
that no panic should occur, without needing to use an indicator field or other mechanism:

```go
defer Expect(Panic(tc.ExpectedPanic)).DidOccur()
```

where `tc.ExpectedPanic` is `nil` for test cases where no panic is expected, and non-`nil` for
test cases where a panic is expected.

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
