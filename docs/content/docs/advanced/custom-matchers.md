---
title: Custom Matchers
weight: 20
---

When the built-in matchers do not cover a specific assertion you need to write repeatedly,
you can implement your own matcher. Custom matchers plug directly into `Expect(...).To(matcher)`
and benefit from the same option handling, failure reporting, and composability as the built-in ones.

---

## The matcher interfaces

Matchers implement one of two interfaces from `github.com/blugnu/test/matchers/matcher`:

```go
// ForType[T] — type-safe matcher for values of type T
type ForType[T any] interface {
    Match(got T, opts ...any) bool
}

// ForAny — matcher for any value (use when the type is not known at compile time)
type ForAny interface {
    Match(got any, opts ...any) bool
}
```

Implement `ForType[T]` whenever possible — it prevents the matcher being used with an incompatible type at compile time.

---

## Minimal implementation

```go
type bePrimeMatcher struct{}

func BePrime() matcher.ForType[int] {
    return &bePrimeMatcher{}
}

func (m *bePrimeMatcher) Match(got int, opts ...any) bool {
    if got < 2 {
        return false
    }
    for i := 2; i*i <= got; i++ {
        if got%i == 0 {
            return false
        }
    }
    return true
}
```

Usage:

```go
Expect(7).To(BePrime())     // passes
Expect(8).To(BePrime())     // fails: no custom message, default is used
```

---

## Adding failure messages

Implement `OnTestFailure` to control what appears in the failure report when `Match` returns `false`:

```go
func (m *bePrimeMatcher) OnTestFailure(got int, opts ...any) []string {
    return []string{
        fmt.Sprintf("expected: a prime number"),
        fmt.Sprintf("got:      %d (not prime)", got),
    }
}
```

The `opts` slice contains the options passed with the assertion (e.g. `opt.Name`, `opt.OnFailure`),
making it possible to integrate option-based customisation.

---

## Using the `Expected` field convention

Many built-in matchers carry an `Expected` field that provides a default failure message
when `OnTestFailure` is not implemented. The framework prints this field if it is present
via reflection.

This is a convention, not an interface:

```go
type rangeContainsMatcher[T constraints.Ordered] struct {
    Expected string   // used by default failure report
    low, high T
}

func BeBetween[T constraints.Ordered](low T) *rangeBuilder[T] {
    return &rangeBuilder[T]{low: low}
}

type rangeBuilder[T constraints.Ordered] struct{ low T }

func (r *rangeBuilder[T]) And(high T) matcher.ForType[T] {
    return &rangeContainsMatcher[T]{
        Expected: fmt.Sprintf("a value between %v and %v", r.low, high),
        low: r.low,
        high: high,
    }
}

func (m *rangeContainsMatcher[T]) Match(got T, opts ...any) bool {
    return got >= m.low && got <= m.high
}
```

---

## Complete example

A custom matcher that checks whether a `*http.Response` has a given status code:

```go
package matchers

import (
    "fmt"
    "net/http"
    "strconv"
    "github.com/blugnu/test/matchers/matcher"
)

type httpStatusMatcher struct {
    Expected int
}

// HaveStatus returns a matcher that checks an *http.Response has the given status code.
func HaveStatus(code int) matcher.ForType[*http.Response] {
    return &httpStatusMatcher{
        Expected: code,
    }
}

func (m *httpStatusMatcher) Match(got *http.Response, opts ...any) bool {
    return got != nil && got.StatusCode == m.Expected
}

func (m *httpStatusMatcher) OnTestFailure(got *http.Response, opts ...any) []string {
    if got == nil {
        return []string{
            "expected status: " + strconv.Itoa(m.Expected),
            "got            : <nil response>",
        }
    }

    return []string{
        fmt.Sprintf("expected status: %d", m.Expected),
        fmt.Sprintf("got            : %d (%s)", got.StatusCode, http.StatusText(got.StatusCode)),
    }
}
```

Using it:

```go
resp, err := http.Get(server.URL + "/ping")
Require(err).IsNil()
Expect(resp).To(HaveStatus(http.StatusOK))
```

Failure output:

```text
expected status: 200
got            :       404 (Not Found)
```

---

## Negation

All matchers automatically support negation — `ToNot(HaveStatus(200))` passes when `Match`
returns `false`. No extra code is needed.

When a matcher is used in a negated assertion, the `opt.ToNotMatch` option is injected into
the options passed to `OnTestFailure`. This allows you to return a different failure message
for the negated case by testing for this option:

```go
func (m *httpStatusMatcher) OnTestFailure(got *http.Response, opts ...any) []string {
    if opt.IsSet(opts, opt.ToNotMatch) {
        return []string{
            fmt.Sprintf("did not expect status: %d", got.StatusCode),
        }
    }

    // message for the non-negated case when got is nil
    if got == nil {
        return []string{
            "expected status: " + strconv.Itoa(m.Expected),
            "got            : <nil response>",
        }
    }

    return []string{
        fmt.Sprintf("expected status: %d", m.Expected),
        fmt.Sprintf("got            : %d (%s)", got.StatusCode, http.StatusText(got.StatusCode)),
    }
}

```

---

## Packaging custom matchers

For reuse across your own projects, package matchers as a separate Go module or internal
package. The `blugnu/test` module itself follows this pattern: each matcher family lives
in a sub-package under `matchers/` and is re-exported from the root package with
constructors like `Equal`, `BeNil`, etc.
