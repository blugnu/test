---
title: Installation
weight: 10
---

## Requirements

- Go 1.22 or later

## Add the module

```bash
go get github.com/blugnu/test
```

## Import style

`blugnu/test` is designed to be **dot-imported**. This means all exported
identifiers — `With`, `Expect`, `Equal`, `Run`, and so on — are available
unqualified in your test files, giving assertions a natural, prose-like
quality.

```go
import . "github.com/blugnu/test"
```

{{< alert title="Note" color="info" >}}
Dot-imports are only recommended in `_test.go` files, where they cannot affect
production code. Using a dot-import in a regular source file is generally
discouraged by the Go community.
{{< /alert >}}

If you prefer, or if a symbol name conflicts with something in your test package, you
can use a named import instead:

```go
import "github.com/blugnu/test"
```

and prefix every call: `test.Expect(...)`, `test.With(t)`, etc.

## Companion packages

The module ships several sub-packages you can import alongside the root package as needed:

| Package | Purpose |
| ------- | ------- |
| `github.com/blugnu/test/expect` | Convenience top-level assertion functions (`expect.Error`, `expect.True`, …) |
| `github.com/blugnu/test/require` | Same as `expect`, but each assertion halts the test on failure |
| `github.com/blugnu/test/opt` | Option constructors (`opt.Name`, `opt.IsRequired`, `opt.OnFailure`, …) |

These packages are not intended to be dot-imported.

## Verify the installation

Create a quick smoke test:

```go
package mypackage_test

import (
    "testing"
    . "github.com/blugnu/test"
)

func TestInstallation(t *testing.T) {
    With(t)
    Expect(1 + 1).To(Equal(2))
}
```

Run it:

```bash
go test ./...
```

You should see `ok` — you are ready to go.
