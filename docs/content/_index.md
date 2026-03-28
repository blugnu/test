---
title: blugnu/test
---

{{< blocks/cover title="blugnu/test" height="min" color="primary" >}}
<p class="lead mt-3">A concise, fluent, type-safe testing framework for Go.</p>
<p class="mt-3">
  <a href="https://github.com/blugnu/test/actions/workflows/release.yml"><img alt="build" src="https://github.com/blugnu/test/actions/workflows/release.yml/badge.svg?branch=master"/></a>
  &nbsp;
  <a href="https://goreportcard.com/report/github.com/blugnu/test"><img alt="go report" src="https://goreportcard.com/badge/github.com/blugnu/test"/></a>
  &nbsp;
  <a href="https://pkg.go.dev/github.com/blugnu/test"><img alt="pkg.go.dev" src="https://pkg.go.dev/badge/github.com/blugnu/test"/></a>
</p>
<div class="mt-4">
  <a class="btn btn-lg btn-light me-3 mb-4" href="{{< relref "/docs/getting-started/installation" >}}">
    Get started <i class="fas fa-arrow-alt-circle-right ms-2"></i>
  </a>
  <a class="btn btn-lg btn-secondary me-3 mb-4" href="https://pkg.go.dev/github.com/blugnu/test">
    API Reference <i class="fab fa-golang ms-2"></i>
  </a>
</div>
{{< /blocks/cover >}}

{{% blocks/section color="white" %}}

<div class="row g-4 align-items-start">
  <div class="col-12 col-lg-6">
  <h2>Instead of this</h2>
{{< highlight go >}}
  result, err := DoSomething()

  if err != nil {
    t.Errorf("unexpected error: %v", err)
  }
  if result != 42 {
    t.Errorf("expected 42, got %v", result)
  }

  t.Run("should panic", func(t *testing.T) {
    defer func() {
      if r := recover(); r == nil {
        t.Errorf("did not panic")
      }
    }()

    DoSomethingThatPanics()
  })
{{< /highlight >}}
  </div>
  <div class="col-12 col-lg-6">
  <h2>Do this</h2>
{{< highlight go >}}
  result, err := DoSomething()

  Expect(err).IsNil()
  Expect(result).To(Equal(42))

  Run(Test("should panic", func() {
    defer Expect(Panic()).DidOccur()

    DoSomethingThatPanics()
  }))
{{< /highlight >}}
  </div>
</div>

{{% /blocks/section %}}

{{% blocks/section color="light" class="home-ladder-section" %}}

<div class="container home-ladder py-4">
  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6">
      <h2>Clean test scaffolding</h2>
      <p class="lead">
        No constant <code>t *testing.T</code> threading through every function call.
        Set up a test frame once with <code>With(t)</code> and keep test code focused
        on behavior.
      </p>
      <p>
        <a class="btn btn-outline-primary" href="{{< relref "/docs/concepts/testframe" >}}">Read testframe concepts</a>
      </p>
    </div>
    <div class="col-12 col-lg-6 text-center home-ladder-image">
      <img
        src="images/home/clean.svg"
        alt="Illustration representing clean test setup"
        class="img-fluid"
      />
    </div>
  </div>

  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6 text-center order-lg-1 home-ladder-image">
      <img
        src="images/home/fluent.svg"
        alt="Illustration representing fluent, readable assertions"
        class="img-fluid"
      />
    </div>
    <div class="col-12 col-lg-6 order-lg-0">
      <h2>Fluent expectations</h2>
      <p class="lead">
        Assertions read like English:
        <code>Expect(result).To(Equal(42))</code>,
        <code>Expect(err).IsNil()</code>,
        <code>Require(id).IsNotNil()</code>.
      </p>
      <p>
        <a class="btn btn-outline-primary" href="{{< relref "/docs/concepts/expectations" >}}">See expectation model</a>
      </p>
    </div>
  </div>

  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6">
      <h2>Type-safe by default</h2>
      <p class="lead">
        Matchers are generic so type mismatches are caught at compile-time,
        not during a flaky CI run.
      </p>
      <p>
        <a class="btn btn-outline-primary" href="{{< relref "/docs/concepts/matchers" >}}">Browse built-in matchers</a>
      </p>
    </div>
    <div class="col-12 col-lg-6 text-center home-ladder-image">
      <img
        src="images/home/typesafe.svg"
        alt="Illustration representing compile-time safety"
        class="img-fluid"
      />
    </div>
  </div>

  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6 text-center order-lg-1 home-ladder-image">
      <img
        src="images/home/table.svg"
        alt="Illustration representing table-driven testcases"
        class="img-fluid"
      />
    </div>
    <div class="col-12 col-lg-6 order-lg-0">
      <h2>Table-driven without boilerplate</h2>
      <p class="lead">
        Use <code>Testcases</code>, <code>Case</code>, <code>Cases</code>,
        <code>Debug</code>, and <code>Skip</code> to keep the test intent clear
        and the loop mechanics out of your way.
      </p>
      <p>
        <a class="btn btn-outline-primary" href="{{< relref "/docs/usage/table-driven-tests" >}}">Learn Testcases</a>
      </p>
    </div>
  </div>

  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6">
      <h2>Composable options and extensibility</h2>
      <p class="lead">
        Fine-tune failure reports, required assertions, and output formatting with
        options. Add custom matchers by implementing a tiny interface.
      </p>
      <p>
        <a class="btn btn-outline-primary me-2" href="{{< relref "/docs/usage/options" >}}">Options guide</a>
        <a class="btn btn-outline-secondary" href="{{< relref "/docs/advanced/custom-matchers" >}}">Custom matchers</a>
      </p>
    </div>
    <div class="col-12 col-lg-6 text-center home-ladder-image">
      <img
        src="images/home/options.svg"
        alt="Illustration representing configurable options and extension points"
        class="img-fluid"
      />
    </div>
  </div>

  <div class="row align-items-center g-5 py-4 home-ladder-item">
    <div class="col-12 col-lg-6 text-center order-lg-1 home-ladder-image">
      <img
        src="images/home/batteries.svg"
        alt="Illustration representing batteries-included testing features"
        class="img-fluid"
      />
    </div>
    <div class="col-12 col-lg-6 order-lg-0">
      <h2>Batteries included</h2>
      <p class="lead">
        Flaky retries, parallel runners, mock helpers, console recording,
        and helper-testing utilities are built in so you can standardize on one toolkit.
      </p>
      <p>
        <a class="btn btn-outline-primary me-2" href="{{< relref "/docs/usage/flaky-tests" >}}">Flaky tests</a>
        <a class="btn btn-outline-secondary" href="{{< relref "/docs/usage/mocking" >}}">Mocking</a>
      </p>
    </div>
  </div>
</div>

{{% /blocks/section %}}
