---
title: Parallel Tests
weight: 40
---

`blugnu/test` fully supports Go's parallel testing model. This page covers the three ways
to run tests in parallel.

## Marking a whole test as parallel

Use `Parallel(t)` instead of `With(t)` when an entire top-level test function should run
in parallel with other top-level tests:

```go
func TestFeatureA(t *testing.T) {
    Parallel(t)   // equivalent to With(t) + t.Parallel()

    result, err := FeatureA()
    Require(err).IsNil()
    Expect(result).To(Equal("A"))
}

func TestFeatureB(t *testing.T) {
    Parallel(t)

    result, err := FeatureB()
    Require(err).IsNil()
    Expect(result).To(Equal("B"))
}
```

Both tests run concurrently. `IsParallel()` returns `true` inside any test marked this way.

{{< alert title="Warning" color="warning" >}}
`Parallel(t)` must not be called from a test that is already parallel (i.e., from inside
a `ParallelTest` or `ParallelCases` subtest).
{{< /alert >}}

---

## Running parallel subtests

Use `ParallelTest` to run individual subtests in parallel within a non-parallel parent:

```go
func TestConcurrentOperations(t *testing.T) {
    With(t)

    // All three subtests run concurrently
    Run(ParallelTest("operation A", func() {
        result := OperationA()
        Expect(result).To(Equal("A"))
    }))

    Run(ParallelTest("operation B", func() {
        result := OperationB()
        Expect(result).To(Equal("B"))
    }))

    Run(ParallelTest("operation C", func() {
        result := OperationC()
        Expect(result).To(Equal("C"))
    }))
}
```

---

## Running parallel table-driven tests

Use `ParallelCases` instead of `Testcases` to run all cases concurrently:

```go
func TestParallelCases(t *testing.T) {
    With(t)

    Run(ParallelCases(
        ForEach(func(tc myCase) {
            result := Process(tc.input)
            Expect(result).To(Equal(tc.want))
        }),
        Cases([]myCase{
            {Scenario: "A", input: "a", want: "A"},
            {Scenario: "B", input: "b", want: "B"},
            {Scenario: "C", input: "c", want: "C"},
        }),
    ))
}
```

Or mix sequential and parallel cases within a single `Testcases` call using `ParallelCase`:

```go
Run(Testcases(
    ForEach(func(tc myCase) { /* ... */ }),
    Case("sequential-only", myCase{...}),
    ParallelCase("can run concurrently", myCase{...}),
    ParallelCase("also concurrent",      myCase{...}),
))
```

---

## Constraints

`Parallel(t)`, `ParallelTest` and `ParallelCases` are all not allowed in the context of a test
that is already parallel. Attempting to use them in such a context will fail the test as invalid.

Use `IsParallel()` to conditionally adjust behaviour, if required:

```go
if !IsParallel() {
    Run(ParallelTest("concurrent setup", func() { /* ... */ }))
}
```

---

## Data safety in parallel tests

Each parallel subtest should work on its own independent data. Avoid sharing mutable variables across subtests:

```go
// BAD — shared variable mutated by parallel tests
var result string
Run(ParallelTest("A", func() { result = "A" }))
Run(ParallelTest("B", func() { result = "B" }))
Expect(result).To(Equal("A"))  // race condition

// GOOD — each subtest uses its own scope
Run(ParallelTest("A", func() {
    result := ComputeA()
    Expect(result).To(Equal("A"))
}))
Run(ParallelTest("B", func() {
    result := ComputeB()
    Expect(result).To(Equal("B"))
}))
```
