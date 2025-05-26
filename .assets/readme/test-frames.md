# Test Frames

The `test` package maintains a stack of `*testing.T` values for each goroutine, updating the stack
when a new `*testing.T` is introduced in the stack, such as when running a subtest.

The `*testing.T` value at the top of the stack at any given time is called the _current **test frame**_.

## Establishing the Initial Test Frame

Each test function must capture the initial `*testing.T` value to establish the initial test frame.
This is done by calling either `test.With(t *testing.T)` or `test.Parallel(t *testing.T)` at the
beginning of the test function:


<!-- markdownlint-disable MD013 // line length -->
| Function | Description |
| --- | --- |
| `test.With(t *testing.T)` | establishes a `*testing.T` value as the current test frame for the goroutine |
| `test.Parallel(t *testing.T)` | establishes a `*testing.T` value as the current test frame for the goroutine and marks the test for parallel execution |
<!-- markdownlint-enable MD013 -->

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  test.With(t) // capture the *testing.T value for this test

  // ... do something with the test ...
}
```

## Subtests and New Test Frames

If a new `*testing.T` is introduced, the `test.With()` or `test.Parallel()` function
must be called to establish it as the current test frame:

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  test.With(t) // capture the *testing.T value for this test

  // ACT
  t.Run("some subtest", func(t *testing.T) {
    test.With(t) // capture the *testing.T value for this subtest
    // ... do something with the subtest ...
  })

  // ... do something with the test ...
}
```

Alternatively, `#blugnu/test` functions for running subtests can be used to automatically
establish the new test frame:

```go
func TestDoSomething(t *testing.T) {
  // ARRANGE
  test.With(t) // capture the *testing.T value for this test

  // ACT
  test.Run("some subtest", func() {
    // ... do something with the subtest ...
  })

  // ... do something with the test ...
}
```

When using these methods, there is no need to call `test.With()` or `test.Parallel()` and there
is no need to pass the `*testing.T` anywhere.
