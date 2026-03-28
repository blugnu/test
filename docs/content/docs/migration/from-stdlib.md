---
title: From the Standard Library
weight: 10
---

This guide shows how to migrate tests written using only the standard
`testing` package to `blugnu/test`. The migration can be done incrementally —
existing `testing`-style tests continue to work alongside new-style tests.

## Before you start

Add the module and set up a dot-import in your test files:

```bash
go get github.com/blugnu/test
```

```go
import . "github.com/blugnu/test"
```

---

## Common patterns

### Error check

**Before:**

```go
result, err := DoSomething()
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

**After:**

```go
result, err := DoSomething()
Require(err).IsNil()   // halts test if err != nil, equivalent to t.Fatal
```

Or if you want to continue collecting failures:

```go
Expect(err).IsNil()    // equivalent to t.Error
```

---

### Equality check

**Before:**

```go
if got != expected {
    t.Errorf("expected %v, got %v", expected, got)
}
```

**After:**

```go
Expect(got).To(Equal(expected))
```

---

### errors.Is check

**Before:**

```go
if !errors.Is(err, ErrNotFound) {
    t.Errorf("expected ErrNotFound, got %v", err)
}
```

**After:**

```go
Expect(err).Is(ErrNotFound)
```

---

### Nil pointer check

**Before:**

```go
if result == nil {
    t.Fatal("expected non-nil result")
}
```

**After:**

```go
Require(result).IsNotNil()
```

---

### Subtests

**Before:**

```go
func TestThings(t *testing.T) {
    t.Run("case A", func(t *testing.T) {
        // ... assertions using t ...
    })
    t.Run("case B", func(t *testing.T) {
        // ... assertions using t ...
    })
}
```

**After:**

```go
func TestThings(t *testing.T) {
    With(t)

    Run(Test("case A", func() {
        // ... assertions — no t parameter needed ...
    }))
    Run(Test("case B", func() {
        // ...
    }))
}
```

---

### Table-driven tests

**Before (typical stdlib pattern):**

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  int
    }{
        {"valid", "42", 42},
        {"zero",  "0",   0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            if err != nil {
                t.Fatal(err)
            }
            if got != tt.want {
                t.Errorf("Parse(%q) = %d, want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

**After:**

```go
type parseCase struct {
    Scenario string
    input    string
    want     int
}

func TestParse(t *testing.T) {
    With(t)

    Run(Testcases(
        ForEach(func(tc parseCase) {
            got, err := Parse(tc.input)
            Require(err).IsNil()
            Expect(got).To(Equal(tc.want))
        }),
        Cases([]parseCase{
            {Scenario: "valid", input: "42", want: 42},
            {Scenario: "zero",  input: "0",  want: 0},
        }),
    ))
}
```

---

### Panic testing

**Before:**

```go
func assertPanics(t *testing.T, fn func()) {
    t.Helper()
    defer func() {
        if r := recover(); r == nil {
            t.Error("expected panic, but none occurred")
        }
    }()
    fn()
}

func TestPanic(t *testing.T) {
    assertPanics(t, func() {
        _ = DoSomethingThatPanics()
    })
}
```

**After:**

```go
func TestPanic(t *testing.T) {
    With(t)
    defer Expect(Panic()).DidOccur()
    _ = DoSomethingThatPanics()
}
```

---

## Incremental adoption

You do not have to convert every test at once. The two styles are fully
compatible — existing `if err != nil { t.Errorf(...) }` code continues to
work. Introduce `With(t)` at the top of any test function you want to convert
and replace assertions one file at a time.
