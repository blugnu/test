---
title: Writing Tests
weight: 10
---

This page covers everyday patterns for writing tests with `blugnu/test`.

## Test structure

Every test function starts with a single call to `With(t)` to initialise the test frame for the current
goroutine. All `blugnu/test` functions depend on this being set.

```go
func TestMyFeature(t *testing.T) {
    With(t)

    // your test body
}
```

If you forget to call `With(t)`, the test will panic with a clear error message:

```text
panic: no test frame; did you forget to call With(t)?
```

See: [Test Frames]({{< relref "/docs/concepts/testframe" >}}) for more information on test frames.

---

## Matcher assertions

`Expect(value)` creates an expectation for a value which may then be tested. `Require(value)` also
creates an expectation with the same capabilities, but which will halt the test immediately if the
expectation is not met.

### Is, To or Should

Expectations support `Is` and `IsNot` assertions for common conditions such as nil checks, error
identity, and equality.  For more complex conditions, expectations also provide `To` and `Should`
assertions which apply a specified matcher to the expectation.

```go
result, err := Compute(42)

Require(err).IsNil()
Expect(result).To(BeGreaterThan(84))
```

### Type-Safe matchers: To / ToNot

Where possible, matchers are strongly typed to the value they operate on, allowing for type-safe
assertions. In the framework these are called *typed-matchers* and are used with the `To` and `ToNot`
assertion methods.

Typed-matchers use generic constraints to specify the types they operate on.

```go
var count int = GetCount()

Expect(count).To(Equal(3))
Expect(count).To(Equal("3")) // compile-time error
```

### Untyped matchers: Should / ShouldNot

If necessary, a matcher can be implemented to operate on `any` values, referred to as *untyped-matchers*.
These are used with the `Should` and `ShouldNot` assertion methods.

Untyped-matchers are used when the types supported by a matcher cannot be expressed as generic constraints.

For example, the `BeNil` matcher is untyped, since the variety of types that could be `nil` cannot be
expressed as a generic constraint.  The following will not compile (unless `value` is explicitly declared
as type `any`):

```go
Expect(value).To(BeNil())
```

Whereas this will compile regardless of the declared type of `value`:

```go
Expect(value).Should(BeNil())
```

In either case, if the underlying concrete type of `value` is not nilable, the test will fail at
runtime as an invalid test, rather than a failed assertion.

> Note: `To` and `Should` are not interchangeable. A typed-matcher that implements `matcher.ForType[T]`
> will not compile when used with `Should`, and an untyped-matcher that implements `matcher.ForAny` will
> not compile when used with `To` (unless the subject of the expectation is explicitly of type `any`).

### Require vs Expect

If an `Expect` assertion fails, the test will continue to evaluate further expectations.  If the
test should fail immediately, you can:

- use `Require` instead of `Expect`;
- pass `opt.Required()` to the assertion method;
- pass `opt.IsRequired(true)` to the assertion method.

```go
Require(result).To(Equal(84))                      // fail immediately on mismatch
Expect(result).To(Equal(84), opt.Required())       // fail immediately on mismatch
Expect(result).To(Equal(84), opt.IsRequired(true)) // fail immediately on mismatch
```

The `opt.IsRequired` option may be useful in table-driven tests where the required status of an
assertion may vary on each test-case.

{{< alert title="Tip" color="success" >}}
Prefer `Require` over `Expect` for precondition checks — values that, if wrong, would make the rest of the test meaningless
or vulnerable to runtime errors. For example, always `Require(result).IsNotNil()` before any tests that dereference `result`.
{{< /alert >}}

---

## Naming expectations

Add a name as the second argument to `Expect` or `Require`. The name appears in the failure report, making it clear which
assertion failed:

```go
Expect(user.Name, "user.Name").To(Equal("Alice"))
Expect(user.Email, "user.Email").To(Equal("alice@example.com"))
```

Failure output:

```text
user.Email:
  expected: "alice@example.com"
  got:      "alice@example.COM"
```

For dynamically named expectations, use `opt.Namef`:

```go
for i, item := range items {
    Expect(item.Valid, opt.Namef("items[%d].Valid", i)).To(BeTrue())
}
```

---

## Structural assertions

### Nil checks

```go
Expect(err).IsNil()           // passes if err == nil
Expect(result).IsNotNil()     // passes if result != nil
```

### Error identity

`Is` uses `errors.Is` under the hood:

```go
Expect(err).Is(ErrNotFound)     // passes if errors.Is(err, ErrNotFound)
Expect(err).IsNot(ErrTimeout)   // passes if !errors.Is(err, ErrTimeout)
```

### Error type checking

```go
Expect(err).Should(BeError[*ValidationError]())  // errors.As check
```

---

## Subtests

Use `Run` with a `Test` runner to create a subtest:

```go
func TestUser(t *testing.T) {
    With(t)

    Run(Test("creation succeeds", func() {
        user, err := NewUser("alice", "alice@example.com")
        Require(err).IsNil()
        Expect(user.Name).To(Equal("alice"))
    }))

    Run(Test("validation rejects empty name", func() {
        _, err := NewUser("", "alice@example.com")
        Expect(err).IsNotNil()
    }))
}
```

Each `Run(Test(...))` call creates a Go subtest. The failure output uses the subtest name,
and `go test -run` can target individual subtests.

---

## Working with context

Matchers are provided for testing context keys and values.

```go
Expect(ctx).To(HaveContextKey("userID"))
Expect(ctx).To(HaveContextValue("userID", "alice"))
```

---

## Recording stdout/stderr

The `Record` function can be used to capture `stdout` and `stderr` output during a test.
The `Record` function does not perform any assertions itself, returning the captured output
as a pair of `[]string` which can then be tested with matchers as usual:

```go
stdout, stderr := Record(func() {
    fmt.Println("Hello, world!")
})

Expect(stdout).To(Equal([]string{"Hello, world!"}))
Expect(stderr).To(BeEmpty())
```

---

## Common patterns

### Test + setup/teardown

```go
func TestDatabase(t *testing.T) {
    With(t)

    db, err := openTestDB()
    Require(err, "open test database").IsNil()
    defer db.Close()

    Run(Test("inserts a record", func() {
        err := db.Insert(record)
        Expect(err).IsNil()
    }))

    Run(Test("returns not found for missing ID", func() {
        _, err := db.Find("nonexistent")
        Expect(err).Is(ErrNotFound)
    }))
}
```

### Nested subtests

```go
func TestHTTPHandler(t *testing.T) {
    With(t)

    Run(Test("POST /items", func() {
        Run(Test("returns 201 on success", func() {
            // ...
        }))
        Run(Test("returns 400 on invalid payload", func() {
            // ...
        }))
    }))
}
```
