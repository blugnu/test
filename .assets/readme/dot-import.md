# dot-import

The `blugnu/test` is recommended to be dot-imported, enabling a more concise and fluent
syntax for writing tests without needing to prefix everything with `test.`.

For example:

```go
import (
  "testing"

  . "github.com/blugnu/test"
)

func TestDoSomething(t *testing.T) {
    // ARRANGE
    With(t) // establish the initial test frame

    // ACT
    got := DoSomething()

    // ASSERT
    Expect(got, "result").To(Equal("expected result")) // no need to prefix with test. when using an anonymous import
}
```

Without the dot-import, the same test would need to be written as:

```go
import (
  "testing"

  "github.com/blugnu/test"
)

func TestDoSomething(t *testing.T) {
    // ARRANGE
    test.With(t) // establish the initial test frame

    // ACT
    got := DoSomething()

    // ASSERT
    test.Expect(got, "result").To(test.Equal("expected result"))
}
```

This still reads _quite_ well, but is more verbose, less concise and is not as fluent.
