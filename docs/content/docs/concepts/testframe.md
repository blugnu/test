---
title: Test Frames
weight: 10
---

The test frame is the foundational concept of `blugnu/test`. Everything else — expectations,
matchers, runners — depends on it.

## The problem with `*testing.T`

Standard Go tests receive a `*testing.T` in every function:

```go
func TestThing(t *testing.T) {
    t.Run("subtest", func(t *testing.T) {
        result := doSomething()
        assertResult(t, result)     // every helper needs t
    })
}

func assertResult(t *testing.T, result int) {
    t.Helper()
    if result != 42 {
        t.Errorf("expected 42, got %d", result)
    }
}
```

Every helper function must accept `t *testing.T` as a parameter and call `t.Helper()` to
keep stack traces clean. This spreads boilerplate across the entire codebase and couples
every assertion utility to the standard `testing` package.

## The solution: a goroutine-local stack

`blugnu/test` maintains a goroutine-local stack of `*testing.T` values called the **test frame stack**.

When you call `With(t)` at the start of a test, the current `*testing.T` is pushed onto
the stack for that goroutine. Any call to the framework within that goroutine — regardless
of call depth — can retrieve it with `T()`.

```yaml
Goroutine stack:
  TestThing → [t₁]
    Run(Test("sub")) → [t₁, t₂]   (t₂ added automatically)
      inner function → [t₁, t₂]   (reads t₂)
    returns → [t₁]                (t₂ automatically popped)
```

## With(t): establishing the frame

`With(t)` pushes the supplied `*testing.T` onto the stack and registers a cleanup function
to pop it when the test completes:

```go
func TestMyThing(t *testing.T) {
    With(t)
    // From this point on, T() returns t
}
```

You should call `With(t)` exactly **once per `*testing.T` value** — at the start
of each top-level `Test...` function.

{{< alert title="Warning" color="warning" >}}
Do not call `With(t)` more than once for the same `t`, and do not call it when the
value has already been registered. If you use `Run(Test(...))` to create subtests,
`With` is called automatically — you do not need to call it again inside the
subtest function.
{{< /alert >}}

## T(): reading the current frame

`T()` returns the `TestingT` currently at the top of the stack:

```go
func assertPositive(n int) {
    T().Helper()
    if n <= 0 {
        T().Errorf("expected positive, got %d", n)
    }
}
```

Notice there is no `t *testing.T` parameter — the helper retrieves it from the stack.
You can write test utilities that are completely decoupled from the `testing` package,
making them trivially reusable across test packages.

`GetT()` is an alias for `T()`, useful in contexts where `T` clashes with a generic
type parameter name.

## How subtests work

When you write:

```go
Run(Test("my subtest", func() {
    Expect(someValue).To(Equal(expected))
}))
```

The `Test` function captures the inner function. When `Run` executes it, it
calls `t.Run(name, func(t *testing.T) { ... })` internally and pushes the new `*testing.T`
onto the stack before the inner function runs, then pops it afterward. The inner function
sees the correct `t` automatically.

If you instead create a subtest directly with `t.Run(...)` — bypassing the test
package — you must push the new frame yourself:

```go
t.Run("my subtest", func(t *testing.T) {
    With(t)   // required when using t.Run directly
    Expect(someValue).To(Equal(expected))
})
```

## Parallel(t): parallel test frames

`Parallel(t)` is a drop-in replacement for `With(t)` when an entire test function should
run in parallel:

```go
func TestMyThing(t *testing.T) {
    Parallel(t)   // equivalent to With(t) + t.Parallel()
    // ...
}
```

Only the top-level test frame should be marked parallel this way. For individual parallel
subtests, use `Run(ParallelTest(...))` instead.

## The TestingT interface

Rather than depending directly on `*testing.T`, the package uses a `TestingT` interface that
mirrors the methods used by the framework. This makes it possible to write test helpers that
work with any type satisfying that interface, which is useful in meta-testing scenarios
(see [Testing Test Helpers]({{< relref "/docs/advanced/testing-helpers" >}})).

## Goroutine safety

The test frame stack is keyed by goroutine ID, derived from `runtime.Stack()`. This makes the
stack safe for concurrent use across goroutines — each goroutine has its own independent stack.

{{< alert title="Warning" color="warning" >}}
The goroutine ID mechanism is based on an implementation detail of the Go runtime that is not
part of the public API. It has been stable across many Go releases, but is not guaranteed to
remain so. See the module README for the latest status.
{{< /alert >}}
