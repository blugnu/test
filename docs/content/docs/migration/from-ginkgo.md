---
title: From Ginkgo / Gomega
weight: 30
---

This guide covers migrating from [`github.com/onsi/ginkgo`](https://onsi.github.io/ginkgo/)
and [`github.com/onsi/gomega`](https://onsi.github.io/gomega/) to `blugnu/test`.

## Philosophical differences

Ginkgo is a BDD framework with its own test runner (`ginkgo` CLI), `Describe`/`Context`/`It` hierarchies,
and a separate `BeforeEach`/`AfterEach` lifecycle. `blugnu/test` is intentionally simpler: it layers a
fluent assertion API on top of the standard `testing` package rather than replacing it.

If your team relies heavily on the BDD vocabulary (Describe, Context, When, It),
`blugnu/test` is a different philosophy rather than a one-to-one replacement.
That said, the migration is straightforward — see
[Mapping the vocabulary](#mapping-the-vocabulary) below.

---

## Gomega matcher equivalents

| Gomega | blugnu/test |
| ------ | ----------- |
| `Expect(x).To(Equal(y))` | `Expect(x).To(Equal(y))` |
| `Expect(x).NotTo(Equal(y))` | `Expect(x).ToNot(Equal(y))` |
| `Expect(x).To(BeNil())` | `Expect(x).IsNil()` |
| `Expect(x).NotTo(BeNil())` | `Expect(x).IsNotNil()` |
| `Expect(x).To(BeTrue())` | `Expect(x).To(BeTrue())` |
| `Expect(x).To(BeFalse())` | `Expect(x).To(BeFalse())` |
| `Expect(err).To(HaveOccurred())` | `Expect(err).IsNotNil()` |
| `Expect(err).NotTo(HaveOccurred())` | `Expect(err).IsNil()` |
| `Expect(err).To(MatchError(target))` | `Expect(err).Is(target)` |
| `Expect(s).To(ContainSubstring(sub))` | `Expect(s).To(ContainString(sub))` |
| `Expect(s).To(MatchRegexp(r))` | `Expect(s).To(MatchRegEx(r))` |
| `Expect(coll).To(HaveLen(n))` | `Expect(coll).To(HaveLen(n))` |
| `Expect(coll).To(BeEmpty())` | `Expect(coll).Should(BeEmpty())` |
| `Expect(coll).To(ContainElement(x))` | `Expect(coll).To(ContainItem(x))` |
| `Expect(fn).To(Panic())` | `defer Expect(Panic()).DidOccur()` |

{{< alert title="Note" color="info" >}}
Gomega's `Expect` and `blugnu/test`'s `Expect` have the same name but
different signatures. During migration, you need to update the import and
adjust the chaining syntax where they differ.
{{< /alert >}}

## Strongly typed matchers

Unlike Gomega, most `blugnu/test` matchers are strongly typed.  This makes them more
robust, less reliant on reflection, and easier to extend but comes at the cost of some
matcher flexibility.

See also:

- [Type-safe vs any-type matchers]({{< relref "/docs/concepts/matchers#type-safe-vs-any-type-matchers" >}}) for
  the rationale behind this design choice.
- [Extending matchers]({{< relref "/docs/advanced/custom-matchers" >}}) for how to add new matcher behavior.

---

## Mapping the vocabulary

### Describe / Context / When → Test

Ginkgo's `Describe` / `Context` hierarchy maps to nested `Run(Test(...))` calls:

**Before:**

```go
var _ = Describe("UserService", func() {
    Context("when creating a user", func() {
        It("returns the new user", func() {
            user, err := svc.CreateUser(ctx, "alice")
            Expect(err).NotTo(HaveOccurred())
            Expect(user.Name).To(Equal("alice"))
        })
        It("returns an error for a duplicate", func() {
            _, err := svc.CreateUser(ctx, "alice")
            Expect(err).To(MatchError(ErrDuplicate))
        })
    })
})
```

**After:**

```go
func TestUserService(t *testing.T) {
    With(t)

    Run(Test("creating a user", func() {
        Run(Test("returns the new user", func() {
            user, err := svc.CreateUser(ctx, "alice")
            Require(err).IsNil()
            Expect(user.Name).To(Equal("alice"))
        }))
        Run(Test("returns an error for a duplicate", func() {
            _, err := svc.CreateUser(ctx, "alice")
            Expect(err).Is(ErrDuplicate)
        }))
    }))
}
```

### BeforeEach → closure / setup function

**Before:**

```go
var svc *UserService

BeforeEach(func() {
    svc = NewUserService(testDB)
})
```

**After:**

```go
setup := func() *UserService {
    return NewUserService(testDB)
}

Run(Test("...", func() {
    svc := setup()
    // ...
}))
```

This is idiomatic Go: each subtest gets its own instance through a factory
function, avoiding shared mutable state between tests.

### Table-driven specs → Testcases

Ginkgo's `DescribeTable` / `Entry`:

**Before:**

```go
DescribeTable("parsing integers",
    func(input string, expected int) {
        got, err := Parse(input)
        Expect(err).NotTo(HaveOccurred())
        Expect(got).To(Equal(expected))
    },
    Entry("positive", "42", 42),
    Entry("zero",     "0",   0),
    Entry("negative", "-1", -1),
)
```

**After:**

```go
type parseCase struct {
    Scenario string
    input    string
    want     int
}

Run(Testcases(
    ForEach(func(tc parseCase) {
        got, err := Parse(tc.input)
        Require(err).IsNil()
        Expect(got).To(Equal(tc.want))
    }),
    Cases([]parseCase{
        {Scenario: "positive", input: "42", want: 42},
        {Scenario: "zero",     input: "0",  want: 0},
        {Scenario: "negative", input: "-1", want: -1},
    }),
))
```

---

## Removing Ginkgo

Once migrated, remove both dependencies:

```bash
go mod tidy
```

Also remove the Ginkgo test package bootstrap:

```go
func TestSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "My Suite")
}
```

Standard `go test ./...` is all you need.
