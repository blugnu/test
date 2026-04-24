---
title: Your First Test
weight: 20
---

This page walks through writing a complete test from scratch using `blugnu/test`.

## The function under test

Suppose you have a small package that parses an integer from a string:

```go
// parse.go
package mypackage

import (
    "fmt"
    "strconv"
)

var ErrNegative = fmt.Errorf("value must be non-negative")

func ParsePositive(s string) (int, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("ParsePositive: %w", err)
    }
    if n < 0 {
        return 0, ErrNegative
    }
    return n, nil
}
```

## Step 1 — Establish a test frame

Every test file that uses `blugnu/test` starts with a single call to `With(t)`:

```go
// parse_test.go
package mypackage_test

import (
    "testing"
    . "github.com/blugnu/test"
)

func TestParsePositive(t *testing.T) {
    With(t)   // establishes the test frame; no cleanup required
}
```

`With(t)` registers the `*testing.T` with a goroutine-local stack managed by the package.
Every subsequent call to `Expect`, `Run`, and other helpers uses this registered value
automatically — you never need to pass `t` around again.

## Step 2 — Assert a result

Use `Expect(value)` to create an expectation. Chain an assertion method on it:

```go
func TestParsePositive(t *testing.T) {
    With(t)

    result, err := ParsePositive("42")

    Expect(err).IsNil()
    Expect(result).To(Equal(42))
}
```

`IsNil()` is a direct assertion method. `To(Equal(42))` delegates to a type-safe
matcher — the compiler will reject `Equal("42")` here because `result` is an `int`,
not a `string`.

## Step 3 — Test the error path

```go
func TestParsePositive_negative(t *testing.T) {
    With(t)

    _, err := ParsePositive("-1")

    Expect(err).IsNotNil()
    Expect(err).Is(ErrNegative)
}
```

`Is(target)` wraps `errors.Is`, so it correctly traverses wrapped error chains.

## Step 4 — Add named subtests

Use `Run(Test("name", fn))` to create readable subtests without passing `t` back in:

```go
func TestParsePositive(t *testing.T) {
    With(t)

    Run(Test("valid input", func() {
        result, err := ParsePositive("7")
        Expect(err).IsNil()
        Expect(result).To(Equal(7))
    }))

    Run(Test("negative value", func() {
        _, err := ParsePositive("-1")
        Expect(err).Is(ErrNegative)
    }))

    Run(Test("non-numeric input", func() {
        _, err := ParsePositive("abc")
        Expect(err).IsNotNil()
    }))
}
```

Each `Test(...)` becomes a proper Go subtest visible in `go test -v` output.

## Step 5 — Halt on a required assertion

When a test only makes sense to continue if an earlier assertion passes, use
`Require` instead of `Expect`:

```go
Run(Test("valid input", func() {
    result, err := ParsePositive("7")
    Require(err).IsNil()          // stops the subtest here if err != nil
    Expect(result).To(Equal(7))  // only runs when err is nil
}))
```

## Running the tests

```bash
go test -v ./...
```

```text
--- PASS: TestParsePositive (0.00s)
    --- PASS: TestParsePositive/valid_input (0.00s)
    --- PASS: TestParsePositive/negative_value (0.00s)
    --- PASS: TestParsePositive/non-numeric_input (0.00s)
PASS
```

## What's next

- Read the [Concepts]({{< relref "/docs/concepts" >}}) section to understand
    the underlying model.
- Jump to [Table-Driven Tests]({{< relref "/docs/usage/table-driven-tests" >}})
    to eliminate the repetitive subtest pattern above.
- See the full [Matcher Reference]({{< relref "/docs/concepts/matchers" >}})
    for all built-in matchers.
- [Migrating From Testify]({{< relref "/docs/migration/from-testify" >}}).
- [Migrating From Ginkgo/Gomega]({{< relref "/docs/migration/from-ginkgo" >}}).
