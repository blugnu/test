---
title: Mocking
weight: 60
---

`blugnu/test` provides two types for implementing test doubles: `FakeResult[R]` for
simple cases and `MockFn[A, R]` for spy and expectation-based mocking.  These are not
mocking frameworks in the traditional sense; they are simple, composable building blocks for
implementing the specific test doubles you need in your tests.  They are designed to be
easy to use and integrate with `blugnu/test` assertions, but do not generate mocks or
verify expectations.

The types are not concurrent safe; tests that run concurrently must use separate
instances.

---

## FakeResult[R] — simple fakes

`FakeResult[R]` is a generic struct for faking a function or method that returns a value
of type `R` and/or an error. It eliminates the boilerplate from common fakes that always
return the same result.

A `FakeResult[R]` has two fields: `Result` of type `R` and `Err` of type `error`. Either
or both can be used depending on the needs of the fake.

No tracking of call arguments or number of calls is provided by `FakeResult[R]` — it simply
returns the configured `Result` and `Err` whenever the fake function is called. For more
advanced mocking needs, see `MockFn[A, R]`, below.

```go
type FakeResult[R any] struct {
    Result R
    Err    error
}
```

> To fake a function that returns multiple values, define a struct for the return type
> with fields for each return value, and use that struct as the type parameter `R`.

### Returns

The `Returns` method accepts values in any order. It sets `Result` from the first `R` value
found and `Err` from the first `error` value found:

```go
fake.Returns(&User{ID: "123"}, nil)     // Result and nil error
fake.Returns(ErrNotFound)               // only an error
fake.Returns(&User{ID: "123"})          // only a result
```

### Reset

`FakeResult[R]` implements `Reset()` which zeroes all fields:

```go
fake.GetUser.Reset()
```

Embedding promotes `Reset` so that if your double only wraps a single `FakeResult[R]`, the
reset is available directly on the double struct.

### Basic usage

```go
type storeDouble struct {
    GetUserFn  FakeResult[*User]
    SaveUserFn FakeResult[*User]
}

func (f *storeDouble) GetUser(id string) (*User, error) {
    return f.GetUserFn.Result, f.GetUserFn.Err
}

// SaveUser returns only the faked error; Result is ignored
func (f *storeDouble) SaveUser(u *User) (error) {
    return f.SaveUserFn.Err
}
```

Set up the fake in your test:

```go
store := &storeDouble{}
store.GetUserFn.Returns(&User{ID: "123", Name: "Alice"})

// or set directly:
store.GetUserFn.Result = &User{ID: "123", Name: "Alice"}
store.GetUserFn.Err = nil
```

Use the double to inject the fake behaviour into the system under test:

```go
service := NewUserService(store)
user, err := service.GetUser("123")
```

---

## MockFn[A, R] — expectation-based mocks

`MockFn[A comparable, R any]` is a more powerful type for:

- capturing and asserting call arguments (spy)
- configuring what each successive call should return (expected-calls mode)
- mapping arguments to return values (mapped-results mode)

### Type parameters

- `A` — the argument type; use a struct for multiple arguments (must be `comparable`)
- `R` — the return type; use a struct for multiple return values

### Embedding in a fake

```go
type cacheDouble struct {
    Set MockFn[struct{ Key, Value string }, any]
    Get MockFn[string, string]
}

func (c *cacheDouble) Set(key, value string) error {
    _, err := c.Set.CalledWith(struct{ Key, Value string }{key, value})
    return err
}

func (c *cacheDouble) Get(key string) (string, bool) {
    result, err := c.Get.CalledWith(key)
    return result, err == nil
}
```

### Expected-calls mode

Configure the calls you expect in order, plus what each should return:

```go
mock.Set.
    ExpectCall().WillReturn(nil).         // first call succeeds
    ExpectCall().WillReturn(ErrFull)      // second call returns an error
```

On each call to `CalledWith`, `MockFn` returns the result configured for the next
expected call. After all expected calls have been consumed, verify them:

```go
ExpectationsWereMet(&mock)
```

### Mapped-results mode

Map specific arguments to specific return values:

```go
mock.Get.
    WhenCalledWith("key-a").Returns("value-a").
    WhenCalledWith("key-b").Returns("value-b")
```

All configured mappings must be used by the end of the test, or `ExpectationsWereMet` will fail.

### Verifying expectations

The `ExpectationsWereMet(mock)` assertion accepts any mock that implements the `Mock` interface:

```go
type Mock interface {
    ExpectationsWereMet() error
    Reset()
}
```

The `Mock` interface is implemented by `MockFn[A, R]` and can be implemented by any custom mock
type. The assertion reports any unmet expectations as a test failure.

```go
ExpectationsWereMet(&mock)
```

This checks that:

- every expected call was made (expected-calls mode)
- every mapped argument was actually called (mapped-results mode)

---

## Choosing between FakeResult and MockFn

| Scenario | Recommended type |
| -------- | ---------------- |
| Always returns the same value | `FakeResult[R]` |
| Number of calls not important | `FakeResult[R]` |
| Arguments not important | `FakeResult[R]` |
| Need to assert call arguments | `MockFn[A, R]` |
| Need to vary return values per call | `MockFn[A, R]` |
| Need to map arguments to return values | `MockFn[A, R]` |
