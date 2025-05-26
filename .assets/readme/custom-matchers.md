# Custom Matcher Implementation

Documentation for implementing custom matchers in the testing framework will be provided here.

Almost all of the matchers provided by the `test` package are implemented as custom matchers
with factory functions in the `test` package itself.  These matchers and factories may be used
as reference implementations.

These may be useful as an introduction:

<!-- markdownlint-disable MD013 -->
| Matcher (factory) | Implemented by |
|---|---|
| `Equal(value)`                  | [`matchers/equal.Matcher`](../../matchers/equal/equal.go) |
| `DeepEqual(value)`              | [`matchers/equal.DeepMatcher`](../../matchers/equal/deepEqual.go) |
| `HaveContextKey(key)`           | [`matchers/contexts.KeyMatcher`](../../matchers/contexts/keyMatcher.go) |
| `HaveContextValue(key, value)`  | [`matchers/contexts.ValueMatcher`](../../matchers/contexts/valueMatcher.go) |
<!-- markdownlint-enable MD013 -->

## Reporting Test Failures

When implementing a custom matcher, the `Match` method returns false when the test has failed.

To report details about the failure, the matcher may implement one of a number of supported
failure reporting methods:

<!-- markdownlint-disable MD013 -->
| Method Signature | Description |
|---|---|
| OnTestFailure(opts ...any) string           | return a formatted string describing the failure |
| OnTestFailure(opts ...any) []string         | return a slice of formatted strings describing the failure |
| OnTestFailure(got T, opts ...any) string    | return a formatted string describing the failure, including the value of `got` |
| OnTestFailure(got T, opts ...any) []string  | return a slice of formatted strings describing the failure, including the value of `got` |
<!-- markdownlint-enable MD013 -->

Where `T` is the type of the value being matched against.

If a matcher does not implement any of these methods, the framework will use a default
failure message that indicates that the matcher failed with the subject value and
expected value, if available.

## Matcher Expected Values

If a matcher is implemented as a struct with an `Expected` field, the framework will
automatically use that in the failure message if no failure report is implemented by
the matcher.

If a matcher only needs to report expected vs actual values, it can avoid the need to
implement a failure report by simply implementing the `Match` method and providing the
`Expected` field.

## Formatting Values for Failure Reports

When implementing failure reports, the `opt` package provides functions to help
with formatting failure messages and applying relevant options:

```go
opt.ValueAsString(value, opts...)    // formats a value as a string for reporting, respecting any supported options
```

Matcher packages may also provide their own options for formatting failure messages,
for example the `matchers/slices` package provides the `slices.AppendToReport` function.
