---
title: Table-Driven Tests
weight: 20
---

`blugnu/test` provides first-class support for table-driven testing through the `Testcases` runner family.
This page covers all the ways to define and run test cases.

## Basic structure

Every `Testcases` invocation requires two things:

1. An **executor** — a function that runs a single test case
2. One or more **case registrations** — values or structs describing individual test cases

```go
Run(Testcases(
    executor,
    registrations...,
))
```

---

## Defining test cases

### Using a struct with Scenario / Name field

The recommended approach: define a struct with a `Scenario` or `Name` string field. This field is
automatically used as the subtest name.

```go
type divideCase struct {
    Scenario string
    dividend int
    divisor  int
    want     int
    wantErr  bool
}

cases := []divideCase{
    {Scenario: "positive result",  dividend: 10, divisor: 2,  want: 5},
    {Scenario: "integer rounding", dividend: 7,  divisor: 2,  want: 3},
    {Scenario: "divide by zero",   dividend: 1,  divisor: 0,  wantErr: true},
}
```

### Using `Cases` (slice)

Register all cases from a slice in one call:

```go
Run(Testcases(
    ForEach(func(tc divideCase) {
        /* ... */
    }),
    Cases([]divideCase{
        {Scenario: "positive result",  dividend: 10, divisor: 2,  want: 5},
        {Scenario: "integer rounding", dividend: 7,  divisor: 2,  want: 3},
        {Scenario: "divide by zero",   dividend: 1,  divisor: 0,  wantErr: true},
    }),
))
```

### Using `Case` (individual, named)

Register one case at a time with an explicit name:

```go
Run(Testcases(
    ForEach(func(tc divideCase) {
        /* ... */
    }),
    Case("even division", divideCase{dividend: 10, divisor: 2, want: 5}),
    Case("rounding towards zero", divideCase{dividend: 7, divisor: 2, want: 3}),
    Case("divide by zero", divideCase{dividend: 1, divisor: 0, wantErr: true}),
))
```

> When registering cases individually, the name is provided with the case registration;
> any `Scenario` or `Name` field in each case is ignored, unless the name provided in the
> registration is empty.

---

## Executor functions

### `ForEach` — name not needed

The simplest executor. The executor function is called once per case, with the case value as its argument.
Any variations in the test must be determined by inspecting the case value:

```go
ForEach(func(tc divideCase) {
    got, err := Divide(tc.dividend, tc.divisor)
    if tc.wantErr {
        Expect(err).IsNotNil()
        return
    }
    Require(err).IsNil()
    Expect(got).To(Equal(tc.want))
})
```

### `For` — name supplied

The `For` executor is called once per case with case name and case value. This enables tests to be
varied based on the case name without needing a `Scenario` field and without having to provide
additional test settings in the case struct itself:

```go
For(func(name string, tc divideCase) {
    switch name {
    case "divide by zero":
        _, err := Divide(tc.dividend, tc.divisor)
        Expect(err).IsNotNil()

    default:
        got, err := Divide(tc.dividend, tc.divisor)
        Require(err).IsNil()
        Expect(got).To(Equal(tc.want))
    }
    // ...
})
```

---

## Controlling individual cases

### Debug

`Debug` marks a case as a debug-only run. When one or more debug cases are present, **only those cases execute** and
the overall test is failed with a warning. Use this to focus on a single failing case without commenting out all the others.

```go
Run(Testcases(
    ForEach(func(tc divideCase) {
        /* ... */
    }),
    Case("positive result",  divideCase{...}),
    Debug("failing case", divideCase{dividend: 7, divisor: 0, wantErr: true}),  // only this runs
    Case("integer rounding", divideCase{...}),
))
```

You can also set a `Debug bool` or `debug bool` field on the test case struct to the same effect — `Cases` will
detect it automatically.

{{< alert title="Warning" color="warning" >}}
`Debug` cases always cause the test suite to fail, even when the case itself passes. This ensures that the
presence of debug cases is clearly visible in test results, preventing accidental commits of tests with cases
marked as debug.  The failure will be clearly identifiable as being due to the presence of debug cases.
{{< /alert >}}

### Skip

```go
Skip("edge case", divideCase{...})
```

Or set a `Skip bool` / `skip bool` field on the struct.

{{< alert title="Warning" color="warning" >}}
`Skip` cases always cause the test suite to fail; this is necessary in order to produce a warning identifying,
the presence of skipped cases.  The test failure will be clearly marked as being due to the case being skipped.
{{< /alert >}}

### ParallelCase

```go
ParallelCase("independent case", divideCase{...})
```

{{< alert title="Warning" color="warning" >}}
`ParallelCase` cannot be used from a test that is already parallel.
{{< /alert >}}

---

## Running cases in parallel

Use `ParallelCases` instead of `Testcases` to run every case in parallel:

```go
Run(ParallelCases(
    ForEach(func(tc scenario) { /* ... */ }),
    Cases(scenarios),
))
```

{{< alert title="Warning" color="warning" >}}
`ParallelCases` cannot be used from a test that is already parallel.
{{< /alert >}}

To run only specific cases in parallel while others run sequentially, use `ParallelCase` within
a regular `Testcases` call.

---

## Complete example

```go
package calculator_test

import (
    "testing"
    . "github.com/blugnu/test"
)

type divideCase struct {
    Scenario string
    a, b     int
    want     int
    wantErr  bool
}

func TestDivide(t *testing.T) {
    With(t)

    Run(Testcases(
        ForEach(func(tc divideCase) {
            got, err := Divide(tc.a, tc.b)
            if tc.wantErr {
                Expect(err).IsNotNil()
                return
            }
            Require(err).IsNil()
            Expect(got, "result").To(Equal(tc.want))
        }),
        Cases([]divideCase{
            {Scenario: "positive values",  a: 10, b: 2,  want: 5},
            {Scenario: "negative dividend", a: -6, b: 3, want: -2},
            {Scenario: "integer rounding",  a: 7,  b: 2, want: 3},
            {Scenario: "divide by zero",    a: 1,  b: 0, wantErr: true},
        }),
    ))
}
```

Output with `go test -v`:

```text
--- PASS: TestDivide (0.00s)
    --- PASS: TestDivide/positive_values (0.00s)
    --- PASS: TestDivide/negative_dividend (0.00s)
    --- PASS: TestDivide/integer_rounding (0.00s)
    --- PASS: TestDivide/divide_by_zero (0.00s)
PASS
```
