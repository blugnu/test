---
title: Flaky Tests
weight: 30
---

A **flaky test** is one that passes and fails non-deterministically, usually due to timing,
external service availability, or environmental conditions. `blugnu/test` provides `FlakyTest`
to contain, manage, and report on such tests without either hiding them or causing persistent
CI noise.

A `FlakyTest` will attempt a specified test a maximum number of times in a specified period,
with a delay between attempts.

The test is retried until either an attempt passes, the maximum number of attempts is reached,
or the maximum duration is exceeded.

If an attempt passes then the test is considered to have passed; no further attempts are made
and no failures are reported even for the failed attempts.

If all attempts fail, the report includes the outcome of each attempt.

## Basic usage

```go
func TestEventuallyConsistent(t *testing.T) {
    With(t)

    Run(FlakyTest("cache is eventually populated",
      func() {
        entry, ok := cache.Get("mykey")
        Expect(ok).To(BeTrue())
        Expect(entry.Value).To(Equal("expected"))
      }, 
      MaxAttempts(5),
      MaxDuration(2*time.Second)),
      WaitBetweenAttempts(100*time.Millisecond)
    )
}
```

The options `MaxAttempts`, `MaxDuration`, and `WaitBetweenAttempts` are all optional
and have sensible defaults if not provided (see below).

Any combination of these options may be specified, or none at all.

---

## Default behaviour

| Setting | Default |
| ------- | ------- |
| Maximum attempts | 3 |
| Maximum duration | 1 second |
| Wait between attempts | 10 ms |

Limits are checked before each attempt:

- stops when `MaxAttempts` is reached;  
- stops when `MaxDuration` is reached;
- if both limits are zero, attempts will continue until either the test
  passes or the `go test` timeout is reached.

---

## Configuration options

Use functional options as the third argument to `FlakyTest`:

```go
Run(FlakyTest("eventually consistent", fn,
    MaxAttempts(5),
    MaxDuration(2*time.Second),
    WaitBetweenAttempts(100*time.Millisecond),
))
```

### MaxAttempts

```go
MaxAttempts(5)    // stop after 5 attempts
MaxAttempts(0)    // no attempt limit
```

### MaxDuration

```go
MaxDuration(500*time.Millisecond)  // stop after 500ms total
MaxDuration(0)                      // no duration limit
```

### WaitBetweenAttempts

```go
WaitBetweenAttempts(50*time.Millisecond)  // wait 50ms between tries
WaitBetweenAttempts(0)                     // retry immediately
```

---

## Failure output

When all attempts fail, the test produces a report that describes every individual failure:

```text
Flaky test failed after 3 attempts in 1.023s

attempt 1:
  expected: true
  got:      false

attempt 2:
  expected: true
  got:      false

attempt 3:
  expected: true
  got:      false
```

This makes it easy to tell whether the test is consistently failing (always the same error) or
genuinely flaky (different errors across attempts).

---

## When to use FlakyTest

`FlakyTest` is appropriate when:

- The system under test involves **timing** (e.g. async processing, caching, event propagation)
- An **external dependency** may be slow to respond under test conditions
- The test is **known to be non-deterministic** and you cannot make it deterministic without
  major refactoring

It is **not** a substitute for fixing a genuinely broken test. If a test consistently fails, fix
the code — do not wrap it in `FlakyTest`.

{{< alert title="Tip" color="success" >}}
If you find yourself setting `MaxAttempts` above 5 or `MaxDuration` above a few seconds, consider
whether there is a more deterministic approach: polling with a channel, using a synchronisation
primitive, or injecting a test clock.
{{< /alert >}}

---

## Complete example

```go
func TestQueueProcessor(t *testing.T) {
    With(t)

    queue := NewQueue()
    processor := NewProcessor(queue)
    go processor.Start()

    // Enqueue an item and wait for it to be processed.
    queue.Push("item-1")

    Run(FlakyTest("item is processed", func() {
        result, ok := processor.GetResult("item-1")
        Require(ok, "result found").To(BeTrue())
        Expect(result.Status).To(Equal("done"))
    },
        MaxAttempts(10),
        MaxDuration(3*time.Second),
        WaitBetweenAttempts(200*time.Millisecond),
    ))
}
```
