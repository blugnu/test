---
title: From Testify
weight: 20
---

This guide covers migrating from [`github.com/stretchr/testify`](https://github.com/stretchr/testify) —
the most widely used third-party assertion library for Go — to `blugnu/test`.

## Key differences

| | testify | blugnu/test |
| -- | ---- | ----------- |
| `t` threading | required by every assertion | managed by `With(t)`; not needed in helpers |
| Assertion style | function calls: `assert.Equal(t, ...)` | fluent: `Expect(...).To(Equal(...))` |
| Required assertions | `require.Equal(t, ...)` | `Require(...).To(Equal(...))` |
| Suites | `testify/suite` | `Run(Test(...))` or `Testcases(...)` |
| Mocking | `testify/mock` | `MockFn[A, R]` |

---

## About `expect` and `require` packages

If your codebase currently uses testify's package split:

- `assert` for non-fatal assertions (test continues after failure)
- `require` for fail-fast assertions (test stops immediately on failure)

`blugnu/test` supports the same mental model in two ways:

- Root package fluent API (primary): `Expect(...)` and `Require(...)`
- Companion packages (convenience): `github.com/blugnu/test/expect` and `github.com/blugnu/test/require`

The `expect` and `require` packages are convenience wrappers around especially
common fluent assertions.

This maps directly to testify usage:

| testify | blugnu/test equivalent |
| ------- | ---------------------- |
| `github.com/stretchr/testify/assert` | `github.com/blugnu/test/expect` or root `Expect(...)` |
| `github.com/stretchr/testify/require` | `github.com/blugnu/test/require` or root `Require(...)` |

Tests can use either the fluent API or the companion packages, or a mix of both. Generally, when mixing
both, the companion packages are used for simple control-flow assertions following execution of some
function under test, with the fluent API then used for more complex assertions:

```go
user, err := svc.CreateUser(ctx, "alice")
require.Nil(err)
require.NotNil(user)

Expect(user.Name, "user name").To(Equal("alice"))
Expect(user.ID, "user ID").ShouldNot(BeEmpty())
```

## Example package-level migration:

```go
import (
    "testing"

    test "github.com/blugnu/test"
    "github.com/blugnu/test/expect"
    "github.com/blugnu/test/require"
)

func TestCreateUser(t *testing.T) {
    test.With(t)

    user, err := svc.CreateUser(ctx, "alice")
    require.Nil(err)
    require.NotNil(user)
    expect.Equal(user.Name, "alice")
}
```

If you prefer, you can skip the companion packages entirely and use only the
root fluent API (`Expect`/`Require`).

This framework is designed around the fluent API as the primary model. In some
cases, fluent expectations can feel more verbose than package-level helper
functions, but they provide an important architectural benefit: matcher-based
extensibility.

You can add new matcher behavior for fluent assertions without modifying
`github.com/blugnu/test` itself. By contrast, adding new top-level
`expect.*`/`require.*` helper functions requires direct changes to the module.

---

## Equivalency table

### assert / require

| testify | blugnu/test fluent API | blugnu/test `expect`/`require` |
| ------- | ---------------------- | ------------------------------ |
| `assert.Equal(t, expected, got)` | `Expect(got).To(Equal(expected))` | `expect.Equal(got, expected)` |
| `assert.NotEqual(t, unexpected, got)` | `Expect(got).ToNot(Equal(unexpected))` | - |
| `assert.Nil(t, val)` | `Expect(val).IsNil()` | `expect.Nil(val)` |
| `assert.NotNil(t, val)` | `Expect(val).IsNotNil()` | `expect.NotNil(val)` |
| `assert.True(t, cond)` | `Expect(cond).To(BeTrue())` | `expect.True(cond)` |
| `assert.False(t, cond)` | `Expect(cond).To(BeFalse())` | `expect.False(cond)` |
| `assert.Error(t, err)` | `Expect(err).IsNotNil()` | `expect.Error(err)` |
| `assert.NoError(t, err)` | `Expect(err).IsNil()` | `expect.NoError(err)` |
| `assert.ErrorIs(t, err, target)` | `Expect(err).Is(target)` | `expect.Error(err, target)` |
| `assert.ErrorAs(t, err, &target)` | `Expect(err).Should(BeError[T]())` | `expect.ErrorAs[T](err)` |
| `assert.Contains(t, s, item)` | `Expect(s).To(ContainItem(item))` | `expect.Slice(got).Contains(item)` |
| `assert.Contains(t, s, slice)` | `Expect(s).To(ContainSlice(slice))` | `expect.Slice(got).Contains(slice...)` |
| `assert.Len(t, coll, n)` | `Expect(coll).To(HaveLen(n))` | - |
| `assert.Empty(t, coll)` | `Expect(coll).Should(BeEmpty())` | - |
| `assert.NotEmpty(t, coll)` | `Expect(coll).ShouldNot(BeEmpty())` | - |
| `assert.EqualValues(t, exp, got)` | `Expect(got).To(DeepEqual(exp))` | - |
| `assert.Panics(t, fn)` | `defer Expect(Panic()).DidOccur()` | - |
| `assert.PanicsWithValue(t, val, fn)` | `defer Expect(Panic(val)).DidOccur()` | - |

### Named assertions

testify uses a `msgAndArgs` trailing variadic parameter for context. `blugnu/test`
uses a name option supplied when creating the expectation. `opt.Namef` is used to
format a name with arguments:

```go
// testify
assert.Equal(t, expected, got, "response")
assert.Equal(t, expected, got, "response from '%s'", endpoint)

// blugnu/test - fluent
Expect(got, "response").To(Equal(expected))
Expect(got, opt.Namef("response from '%s'", endpoint)).To(Equal(expected))

// blugnu/test - expect package
expect.Equal(got, expected, "response")
expect.Equal(got, expected, opt.Namef("response from '%s'", endpoint))
```

---

## Worked example

**Before (testify):**

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
    user, err := svc.CreateUser(ctx, "alice")
    require.NoError(t, err)
    require.NotNil(t, user)
    assert.Equal(t, "alice", user.Name)
    assert.NotEmpty(t, user.ID)
}
```

**After (blugnu/test):**

```go
import (
    "testing"
    . "github.com/blugnu/test"
)

func TestCreateUser(t *testing.T) {
    With(t)

    user, err := svc.CreateUser(ctx, "alice")
    Require(err).IsNil()
    Require(user).IsNotNil()
    Expect(user.Name).To(Equal("alice"))
    Expect(user.ID).ShouldNot(BeEmpty())
}
```

---

## Suites

testify suites use `SetupTest` / `TearDownTest` and embed `suite.Suite`.
`blugnu/test` achieves the same with idiomatic Go setup functions and subtests:

**Before:**

```go
type UserServiceSuite struct {
    suite.Suite
    svc *UserService
}

func (s *UserServiceSuite) SetupTest() {
    s.svc = NewUserService(testDB)
}

func (s *UserServiceSuite) TestCreateUser() {
    user, err := s.svc.CreateUser(...)
    s.Require().NoError(err)
    // ...
}
```

**After:**

```go
func TestUserService(t *testing.T) {
    With(t)

    setup := func() *UserService {
        return NewUserService(testDB)
    }

    Run(Test("creates a user", func() {
        svc := setup()
        user, err := svc.CreateUser(...)
        Require(err).IsNil()
        // ...
    }))

    Run(Test("returns error for duplicate", func() {
        svc := setup()
        // ...
    }))
}
```

---

## Removing testify

Once all assertions have been migrated, remove the dependency:

```bash
go mod tidy
```
