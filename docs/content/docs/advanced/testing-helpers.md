---
title: Testing Test Helpers
weight: 10
---

A *test helper* is a function called from tests that internally uses `blugnu/test` assertions.
To verify that a helper behaves correctly — that it passes when it should and fails when it
should — `blugnu/test` provides the `TestHelper` function and the `HelperTests` runner.

## The concept

Test helpers accept a function that exercises the code under test. `TestHelper` captures the
outcome of running that function — whether it passed, failed, or panicked — along with any
logged output, and returns a result value `R` that you can assert against.

---

## Running a single helper test

```go
func TestMyHelper(t *testing.T) {
    With(t)

    result := TestHelper(func() {
        // exercise the helper under test
        CheckPositive(-1)
    })

    result.Expect(test.Failed)
}
```

### The `R` type

`TestHelper` returns an `R` value with the following fields:

| Field | Type | Description |
| ----- | ---- | ----------- |
| `Outcome` | `test.Outcome` | `test.Passed`, `test.Failed`, or `test.Panicked` |
| `Report` | `[]string` | failure report lines emitted by the helper |
| `Log` | `[]string` | log output captured from the helper |
| `Recovered` | `any` | recover value (if the helper panicked) |
| `Stack` | `[]byte` | goroutine stack (if the helper panicked) |
| `FailedTests` | `[]string` | names of any subtests that failed |

### `R.Expect`

Assert the outcome and optionally match against expected failure report lines:

```go
// Assert the test passed
result.Expect(test.Passed)

// Assert the test failed
result.Expect(test.Failed)

// Assert it failed and the failure message contains a specific line
result.Expect(test.Failed, "expected: 1")

// Assert it failed and contains multiple specific lines
result.Expect(test.Failed, "expected: 1", "got:      -1")

// Assert the test panicked
result.Expect(test.Panicked)
result.Expect(test.Panicked, "index out of range")
```

{{< alert title="Note" color="info" >}}
`result.Expect` uses partial string matching — a line in `Report` only has to *contain* the
expected string, not equal it.
{{< /alert >}}

---

## The `HelperTests` runner

For helpers with several pass/fail scenarios, use `HelperTests` to drive declarative test tables:

```go
func TestCheckPositive(t *testing.T) {
    With(t)

    Run(HelperTests([]HelperScenario{
        {Scenario: "positive value passes",
            Act: func() { CheckPositive(1) },
            Assert: func(r R) { r.Expect(test.Passed) },
        },
        {Scenario: "zero fails",
            Act: func() { CheckPositive(0) },
            Assert: func(r R) { r.Expect(test.Failed, "expected: positive") },
        },
        {Scenario: "negative value fails",
            Act: func() { CheckPositive(-5) },
            Assert: func(r R) { r.Expect(test.Failed, "got: -5") },
        },
    }...))
}
```

Each scenario becomes a subtest named after the `Scenario` field. The `Act` function
runs the helper, and `Assert` receives the captured `R` value to assert against.

---

## Complete example

Suppose you have a helper that checks a string is a valid email address:

```go
// checkEmail.go
func CheckEmail(email string) {
    if !strings.Contains(email, "@") {
        Expect(email).Should(MatchRegEx(`^[^@]+@[^@]+\.[^@]+$`))
    }
}
```

The tests for this helper:

```go
// checkEmail_test.go
func TestCheckEmail(t *testing.T) {
    With(t)

    Run(HelperTests(
        HelperScenario{
            Scenario: "valid email passes",
            Act:    func() { CheckEmail("alice@example.com") },
            Assert: func(r R) { r.Expect(test.Passed) },
        },
        HelperScenario{
            Scenario: "missing @ fails",
            Act:    func() { CheckEmail("notanemail") },
            Assert: func(r R) { r.Expect(test.Failed) },
        },
        HelperScenario{
            Scenario: "empty string fails",
            Act:    func() { CheckEmail("") },
            Assert: func(r R) { r.Expect(test.Failed) },
        },
    ))
}
```

---

## When to use TestHelper vs HelperTests

| Scenario | Recommendation |
| -------- | -------------- |
| Testing one specific edge case in isolation | `TestHelper` |
| Testing many pass/fail scenarios for one helper | `Run(HelperTests(...))` |
| Testing output content (log lines, failure messages) | Either — `R.Expect` accepts expected lines |
