---
title: Test Runners
weight: 40
---

A **runner** is any value that implements the `Runnable` interface:

```go
type Runnable interface {
    Run()
}
```

The `Run(r Runnable)` function executes a runner inside the current test frame. The
framework provides three built-in runner types: `Test`, `FlakyTest`, and `Testcases`.

## Test

`Test` creates a named subtest. It is the primary building block for structuring tests:

```go
func TestOrderService(t *testing.T) {
    With(t)

    Run(Test("places an order", func() {
        order, err := svc.PlaceOrder(ctx, item)
        Require(err).IsNil()
        Expect(order.Status).To(Equal("pending"))
    }))

    Run(Test("rejects duplicate orders", func() {
        _, err := svc.PlaceOrder(ctx, item)
        Expect(err).Is(ErrDuplicate)
    }))
}
```

The inner function receives no arguments. The test frame is pushed automatically before the
function runs and popped when it returns.

### ParallelTest

`ParallelTest` creates a subtest that runs in parallel with other parallel subtests of the
same parent. It cannot be used from a test that is already parallel.

```go
func TestConcurrentAccess(t *testing.T) {
    With(t)

    Run(ParallelTest("goroutine A", func() { /* ... */ }))
    Run(ParallelTest("goroutine B", func() { /* ... */ }))
}
```

---

## FlakyTest

`FlakyTest` wraps a test function that may be non-deterministic or timing-dependent.
It retries the function until it passes, up to a configured limit.

```go
Run(FlakyTest("cache is eventually consistent", func() {
    entry, ok := cache.Get("key")
    Expect(ok).To(BeTrue())
    Expect(entry.Value).To(Equal("expected"))
}))
```

### Default behaviour

| Parameter | Default |
| --------- | ------- |
| Maximum attempts | 3 |
| Maximum duration | 1 second |
| Wait between attempts | 10 ms |

The test passes as soon as any attempt succeeds; failure reports from preceding attempts
are discarded. If every attempt fails, all failure reports are included in the output.

### Options

```go
Run(FlakyTest("eventually consistent", fn,
    MaxAttempts(5),                      // allow up to 5 attempts
    MaxDuration(500*time.Millisecond),   // … but stop after 500ms
    WaitBetweenAttempts(50*time.Millisecond), // wait 50ms between attempts
))
```

Setting either limit to `0` disables it. Setting both to `0` allows unlimited retries until the `go test` timeout.

---

## Testcases

`Testcases` is the table-driven test runner. It pairs a set of **test cases** with
a **test executor** function and runs each case as its own subtest.

```go
type scenario struct {
    Scenario string
    input    string
    want     int
    wantErr  bool
}

func TestParse(t *testing.T) {
    With(t)

    Run(Testcases(
        ForEach(func(tc scenario) {
            got, err := Parse(tc.input)
            if tc.wantErr {
                Expect(err).IsNotNil()
                return
            }
            Require(err).IsNil()
            Expect(got).To(Equal(tc.want))
        }),
        Cases([]scenario{
            {Scenario: "valid integer",  input: "42",  want: 42},
            {Scenario: "negative value", input: "-1",  want: -1},
            {Scenario: "empty string",   input: "",    wantErr: true},
        }),
    ))
}
```

### Providing a test executor

| Function | Description |
| -------- | ----------- |
| `ForEach(func(T))` | executor receives just the test case |
| `For(func(string, T))` | executor receives the test case name and the test case — useful when execution varies by name |

### Registering test cases

| Function | Description |
| -------- | ----------- |
| `Case[T]("name", tc)` | single named case |
| `Cases[T]([]T{…})` | slice of cases; names are inferred from `Scenario`/`Name` field if present |
| `Debug[T]("name", tc)` | mark this case as a debug case; only debug cases will run (test fails with a warning) |
| `Skip[T]("name", tc)` | skip this case |
| `ParallelCase[T]("name", tc)` | run this case in parallel |

### Automatic name inference

When using `Cases(slice)`, if the element type is a struct (or pointer to struct) with
a `string` field named `Scenario`, `scenario`, `Name`, or `name`, that field's value is used
as the subtest name. Otherwise names are generated as `testcase-001`, `testcase-002`, etc.

### ParallelCases

`ParallelCases` is identical to `Testcases` except that each case runs in parallel. It cannot
be used from a test that is already parallel.

```go
Run(ParallelCases(
    ForEach(func(tc scenario) { /* … */ }),
    Cases(scenarios),
))
```

---

## The Runnable interface and custom runners

Because `Run` accepts the `Runnable` interface, you can implement your own runner types:

```go
type retryRunner struct {
    attempts int
    fn       func()
}

func (r retryRunner) Run() {
    for i := 0; i < r.attempts; i++ {
        // ...
    }
}

Run(retryRunner{attempts: 3, fn: func() { /* ... */ }})
```

This is effectively how `FlakyTest` is implemented internally.
