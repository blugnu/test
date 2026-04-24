---
title: Expectations
weight: 20
---

An **expectation** is a value together with an assertion about it. The `Expect` and
`Require` functions create expectations; methods on the returned value perform the assertion.

## Creating an expectation

```go
Expect(someValue)          // unnamed expectation
Expect(someValue, "label") // named expectation — label appears in failure output
Expect(someValue, opt.Name("label"))   // equivalent
Expect(someValue, opt.Namef("item %d", i))  // formatted name
```

The value can be of any type — the returned `*expectation[T]` is generic and carries the type information forward.

## Expect vs Require

| | `Expect` | `Require` |
| - | ------ | --------- |
| On failure | records the failure and **continues** the test | records the failure and **halts** the current test |
| Use when | you want to collect all failures in one run | subsequent assertions only make sense if this one passes |

```go
// collect all failures
Expect(err).IsNil()
Expect(result.ID).To(Equal(expectedID))
Expect(result.Name).To(Equal(expectedName))

// stop early when it doesn't make sense to continue
Require(err).IsNil()           // test stops here if err != nil
Expect(result.ID).To(Equal(expectedID))   // only reached when err is nil
```

`Require(value)` is equivalent to `Expect(value, opt.Required())`.

You can also mark an individual assertion as required using an option:

```go
Expect(value).To(Equal(expected), opt.Required())
Expect(value).To(Equal(expected), opt.IsRequired(true))
```

## Direct assertion methods

These methods are defined directly on the expectation and do not involve a separate matcher:

| Method | Passes when |
| ------ | ----------- |
| `IsNil()` | value is nil |
| `IsNotNil()` | value is not nil |
| `Is(target)` | `errors.Is(value, target)` returns true |
| `IsNot(target)` | `errors.Is(value, target)` returns false |

`Is` and `IsNot` are primarily for error assertions with wrapped errors.

## Matcher-based assertions

Most assertions are performed via a matcher:

| Method | Matcher type | Use when |
| ------ | ------------ | -------- |
| `To(m)` | `matcher.ForType[T]` | matcher and expectation share the same type |
| `ToNot(m)` | `matcher.ForType[T]` | negated type-safe assertion |
| `Should(m)` | `matcher.ForAny` | matcher operates on `any` (e.g. `BeNil`) |
| `ShouldNot(m)` | `matcher.ForAny` | negated any-type assertion |

```go
Expect(n).To(Equal(42))           // type-safe: Equal[int]
Expect(n).ToNot(Equal(0))         // negated
Expect(n).Should(BeNil())         // any-type matcher
Expect(n).ShouldNot(BeNil())      // negated
```

Most built-in factory functions return `matcher.ForType[T]` matchers and therefore
require the `To` / `ToNot` methods. The `BeNil` and `BeEmpty` family return
`matcher.ForAny` matchers and require `Should` / `ShouldNot`.

{{< alert title="Tip" color="success" >}}
`IsNil()` and `IsNotNil()` are convenience wrappers around `Should(BeNil())` and
`ShouldNot(BeNil())` respectively — they exist to read more naturally in code.
{{< /alert >}}

## Failure messages and options

All assertion methods accept variadic options:

```go
Expect(result).To(Equal(expected), opt.OnFailure("result was wrong"))
Expect(result).To(Equal(expected), opt.OnFailure(func(opts ...any) []string {
    return []string{
        fmt.Sprintf("wanted: %v", expected),
        fmt.Sprintf("got:    %v", result),
    }
}))
```

Commonly used options:

| Option | Effect |
| ------ | ------ |
| `opt.OnFailure(...)` | replaces the default failure message with a `string`, `[]string` or `opt.FailReporter` function |
| `opt.IsRequired(bool)` / `opt.Required()` | marks an assertion as required |
| `opt.QuotedStrings` | quotes string values in failure reports (default) |
| `opt.UnquotedStrings` | removes quotes from string values in failure reports |
