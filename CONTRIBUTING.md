# Contributing

Thanks for considering a contribution to `statstest`.

## Before you start

For anything beyond a small fix (a new test, a change in behavior, a new
dependency), please open an issue first to discuss the approach. It's a
lot easier to agree on the shape of a change before code is written than
after.

## Development

```bash
go build ./...
go vet ./...
go test ./...
gofmt -l .   # should print nothing
```

All of the above must pass before a PR will be merged; CI runs them
against the minimum supported Go version (see the `go` directive in
[go.mod](go.mod)) as well as newer releases.

## Adding or changing a statistical method

This package's core value is that its numbers are verifiably correct, not
just "it runs." If you add or modify a test, effect size, or distribution
computation, please include a **reference-value test** in
[reference_test.go](reference_test.go) that checks the implementation's
output against an independent, external source — R, SciPy, or a published
worked example from a textbook or paper. A test that only checks the code
doesn't panic and returns a p-value in `[0, 1]` is not sufficient on its
own.

If a method has more than one valid convention in the statistical
literature (e.g. which tie-correction formula, which variance estimator),
say so in a doc comment and note the convention this package follows, the
way existing functions like `ProportionOneSample` and `LevenesTest` do.

Also add:
- Validation-path unit tests (invalid input, edge cases) in a
  `<file>_test.go` alongside the implementation.
- A runnable `Example` function in a root-level `example_*_test.go` file
  (`package statstest_test`), so it's attached to the function's
  documentation on pkg.go.dev.

## Reporting a bug

Please include the input that triggers it and, if possible, the value you
expected versus what the package returned — ideally cross-checked against
R or SciPy.
