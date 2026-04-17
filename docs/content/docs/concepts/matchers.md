---
title: Matchers
weight: 30
---

A **matcher** evaluates whether a value satisfies some condition. Matchers are returned by
factory functions and passed to `To`, `ToNot`, `Should`, or `ShouldNot` on an expectation.

## Type-safe vs any-type matchers

There are two matcher interfaces:

```go
// matcher.ForType[T] — type-safe; used with To / ToNot
type ForType[T any] interface {
    Match(got T, opts ...any) bool
}

// matcher.ForAny — works on any value; used with Should / ShouldNot
type ForAny interface {
    Match(got any, opts ...any) bool
}
```

The generic `ForType[T]` interface catches type mismatches at compile time. `ForAny` is used
for matchers that are inherently type-agnostic (such as nil checks).

---

## Choosing `To` vs `Should`

`To` and `ToNot` are for strongly typed matchers (`matcher.ForType[T]`). Use these when the
matcher is constrained to a specific type at compile time.

```go
Expect(count).To(Equal(3))
Expect(name).ToNot(Equal(""))
Expect(score).To(BeGreaterThan(0))
```

`Should` and `ShouldNot` are for untyped matchers (`matcher.ForAny`). Use these when the matcher
is intentionally (or by necessity) type-agnostic.

```go
Expect(err).Should(BeNil())
Expect(value).ShouldNot(BeEmpty())
```

Rule of thumb:

- Prefer `To`/`ToNot` for compile-time type safety.
- Use `Should`/`ShouldNot` where the matcher is defined for `any`.
- If a matcher only implements `matcher.ForAny`, `To` will not compile.

---

## Built-in matchers reference

### Equality

| Factory | Description |
| ------- | ----------- |
| `Equal[T comparable](want T)` | equality via `==`, or `T.Equal(T)` if defined, or a custom comparison function |
| `DeepEqual[T any](want T)` | equality via `reflect.DeepEqual` — works on non-comparable types |

```go
Expect(result).To(Equal(42))
Expect(bigStruct).To(DeepEqual(expected))

// custom comparator supplied as a matcher option
Expect(result).To(Equal(expected), func(a, b MyType) bool {
    return a.ID == b.ID
})
```

### Nil

| Factory | Description |
| ------- | ----------- |
| `BeNil()` | value is nil (`matcher.ForAny`) |

```go
Expect(err).Should(BeNil())          // explicit
Expect(err).IsNil()                  // shorthand
Expect(ptr).IsNotNil()               // negated shorthand
```

### Boolean

| Factory | Description |
| ------- | ----------- |
| `BeTrue()` | value is true |
| `BeFalse()` | value is false |

```go
Expect(ok).To(BeTrue())
Expect(found).To(BeFalse())
```

### Emptiness

| Factory | Description |
| ------- | ----------- |
| `BeEmpty()` | collection/string has length zero |
| `BeEmptyOrNil()` | collection/string has length zero **or** is nil |

```go
Expect(items).Should(BeEmpty())
Expect(slice).Should(BeEmptyOrNil())
```

### Length

| Factory | Description |
| ------- | ----------- |
| `HaveLen(n int)` | collection/string has exactly n elements |

```go
Expect(items).To(HaveLen(3))
```

### Ordering (requires `cmp.Ordered` constraint)

| Factory | Description |
| ------- | ----------- |
| `BeGreaterThan[T](want T)` | value > want |
| `BeLessThan[T](want T)` | value < want |
| `BeBetween[T](lo T).And(hi T)` | lo ≤ value ≤ hi (inclusive) |

```go
Expect(n).To(BeGreaterThan(0))
Expect(score).To(BeBetween(1).And(100))
```

### Strings

| Factory | Description |
| ------- | ----------- |
| `ContainString(s string)` | string contains the substring s |
| `MatchRegEx(pattern string)` | string matches the regular expression |

```go
Expect(msg).To(ContainString("timeout"))
Expect(line).To(MatchRegEx(`^\d{4}-\d{2}-\d{2}`))
```

### Slices

| Factory | Description |
| ------- | ----------- |
| `ContainItem[T](item T)` | slice contains the item |
| `ContainItems[T](items []T)` | slice contains all of the given items |
| `ContainSlice[T](sub []T)` | slice contains the sub-slice as a contiguous run |
| `EqualSlice[T](want []T)` | slice is equal to want (element by element) |

```go
Expect(tags).To(ContainItem("go"))
Expect(ids).To(ContainItems([]int{1, 2, 3}))
Expect(data).To(EqualSlice(expected))
```

### Maps

| Factory | Description |
| ------- | ----------- |
| `ContainMap[K,V](want map[K]V)` | map contains all key-value pairs in want |
| `ContainMapEntry[K,V](key K, val V)` | map contains the specific key-value pair |
| `EqualMap[K,V](want map[K]V)` | maps are equal |

```go
Expect(m).To(ContainMapEntry("status", "ok"))
Expect(m).To(EqualMap(expected))
```

Helper functions for extracting keys/values before matching:

```go
Expect(KeysOfMap(m)).To(ContainItem("status"))
Expect(ValuesOfMap(m)).To(ContainItem("ok"))
```

### Errors

| Factory | Description |
| ------- | ----------- |
| `BeError[E error](targets ...E)` | error satisfies `errors.As[E]`; optional targets checked with `errors.Is` |

```go
Expect(err).Should(BeError[*os.PathError]())       // type check only
Expect(err).Should(BeError[*os.PathError](target)) // type + value check
```

See also `Expect(err).Is(target)` for `errors.Is` checks.

### Panics

Panic matchers are used slightly differently than other matchers since they *establish* an expectation
rather than verify that an expectation was met.

A `Panic` matcher is provided to the `Expect` function, and then verified with a deferred `DidOccur`
or `DidNotOccur` assertion:

```go
defer Expect(Panic("invalid input")).DidOccur()

MustParse("bad")
```

| Factory | Description |
| ------- | ----------- |
| `Panic()` | function panics, regardless of recovered value |
| `Panic(value)` | function panics with a specific recovered value |
| `Panic(nil)` | code under test did not panic; used in table-driven tests where a `nil` recovery value in a test case indicates no panic |
| `NilPanic()` | code under test called `panic(nil)` |

Panic testing is covered in more detail in [Panic Testing]({{< relref "/docs/concepts/panic-testing" >}}).

### Type assertions

| Factory | Description |
| ------- | ----------- |
| `BeOfType[T any]()` | value's dynamic type is T |

```go
Expect(val).Should(BeOfType[*MyStruct]())
```

For a type assertion that also returns the typed value for further testing, see the `expect.Type` and `require.Type` functions.

### Mock expectations

| Function | Description |
| -------- | ----------- |
| `ExpectationsWereMet(mock)` | all configured expectations on a `Mock` were satisfied |

```go
ExpectationsWereMet(myMock)
```

---

## Matcher options

Most matchers accept options forwarded from `To(matcher, opts...)` or `Should(matcher, opts...)`:

| Option | Effect |
| ------ | ------ |
| `opt.OnFailure(string)` or `opt.OnFailure([]string)` | fixed failure message |
| `opt.OnFailure([]string)` | multi-line, fixed failure message |
| `opt.OnFailure(func(...any) []string)` | dynamic failure message, accepting options |
| `opt.QuotedStrings` | quote string values (default) |
| `opt.UnquotedStrings` | do not quote string values in the failure report |

---

## Writing a custom matcher

See [Custom Matchers]({{< relref "/docs/advanced/custom-matchers" >}}) for the full guide.
