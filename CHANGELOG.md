# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project follows Semantic Versioning while the API is still evolving.

## [Unreleased]

### Added
- Committed, permanent benchmarks (`benchmark_test.go`) for every
  exported function, runnable with `go test -bench=. -benchmem ./...`,
  so performance is measured going forward instead of asserted once.
- CI now also runs golangci-lint on every push/PR (separate `lint` job),
  and `go test` runs with `-race`.
- `TukeyHSDContext(ctx, conf, groups...)`: a context-aware variant of
  `TukeyHSD` for callers on a deadline (e.g. an HTTP handler). It checks
  `ctx` before starting and stops dispatching further pairwise
  comparisons as soon as `ctx` is done, returning a zero `TukeyHSDResult`
  and `ctx.Err()` rather than a partial result that could be mistaken for
  a complete one. `TukeyHSD` is now a thin wrapper around it using
  `context.Background()`, so its behavior is unchanged.
- Benchmarked `ShapiroWilk` at its largest supported size (n=5000, ~0.3ms)
  and `DunnTest` at k=20 groups (~0.08ms) before deciding cancellation
  support wasn't worth adding there — both finish far too fast for it to
  matter. `TukeyHSD` remains the only function in the package where a
  caller-driven timeout is a real concern (tens to hundreds of ms at
  larger group counts).

### Changed
- Cut `studentizedRangeQuantile`'s bisection from 60 to 35 iterations
  (`bisectionIterations` in `studentizedrange.go`). Benchmarking `TukeyHSD`
  showed this single bisection — run once per call, independent of group
  count — dominated its runtime more than the pairwise comparison loop
  that was parallelized in 0.2.1. Empirically, the bisection converges to
  the numerical integration's own resolution by ~15-20 iterations; every
  case in `reference_test.go` still matches its existing tolerance after
  the change. This is a ~25-35% reduction in `TukeyHSD`'s total latency
  depending on group count, with the accuracy the package already had.

### Fixed
- An ineffectual assignment in `FishersExact2x2` (`categorical.go`),
  caught by adding golangci-lint to CI. No behavior change — every branch
  already overwrote the value before it was used.

## [0.2.1] - 2026-09-10

### Added
- Parallelized `TukeyHSD`'s pairwise comparison loop: each comparison's
  p-value requires a numerical integration (`studentizedRangeCDF`) that is
  independent of every other pair, so for enough groups this is now
  spread across goroutines (bounded by `GOMAXPROCS`) instead of computed
  one pair at a time. Benchmarked at k=10 groups (45 pairs): ~81ms → ~57ms;
  the gain grows with group count since the parallelized part dominates
  more of the total time. `DunnTest`'s pairwise loop was deliberately left
  sequential — its per-pair cost is a closed-form normal-CDF evaluation,
  cheap enough that parallelizing it would add complexity with no
  measurable payoff.

## [0.2.0] - 2026-09-09

### Added
- Reference-value verification tests for core statistical methods.
- Verified test coverage for:
  - one-sample t-test
  - two-sample Welch t-test
  - paired t-test
  - one-way ANOVA
  - Mann-Whitney U test
  - Wilcoxon signed-rank test
  - chi-square test of independence
  - Fisher’s exact test
  - p-value adjustment methods
- New tests:
  - Pearson correlation (`PearsonCorrelation`)
  - Spearman rank correlation (`SpearmanCorrelation`)
  - Kruskal-Wallis test (`KruskalWallis`)
  - one-sample and two-sample proportion z-tests (`ProportionOneSample`, `ProportionTwoSample`)
  - chi-square goodness-of-fit test (`ChiSquareGoodnessOfFit`)
  - Levene's test for equal variances, Brown-Forsythe formulation (`LevenesTest`)
  - Shapiro-Wilk normality test (`ShapiroWilk`), via Royston's algorithm AS R94
  - Tukey HSD post-hoc pairwise comparisons following `OneWayANOVA` (`TukeyHSD`)
  - Dunn's post-hoc pairwise comparisons following `KruskalWallis` (`DunnTest`)
  - a from-scratch, numerically-integrated studentized range distribution
    (`studentizedRangeCDF`/`studentizedRangeQuantile`), needed for `TukeyHSD`
    since neither the Go standard library nor gonum provide one
- Reference-value tests for all newly added methods, checked against SciPy, statsmodels, and scikit-posthocs.
- Runnable `Example*` functions for every exported test function, attached to their documentation on pkg.go.dev, plus an unqualified package-level `Example()`.

### Changed
- Improved statistical verification by checking implementation outputs against external reference values and documented package conventions.
- Updated convention-sensitive reference tests for:
  - Mann-Whitney U asymptotic p-value behavior
  - Wilcoxon signed-rank statistic reporting convention (`W+`)
- Moved `Example*` functions from the separate `examples/` package into the root package's test files so they render inline on pkg.go.dev.
- Lowered the minimum Go version from 1.25.1 to 1.21, and gonum from v0.17.0 to v0.15.0 (the newest release that still supports Go 1.21). Nothing in this package needs a newer toolchain; the previous minimum was an unnecessary adoption barrier.

### Removed
- Unused, duplicate validation helpers in `validate.go` (superseded by the helpers already used in `helpers.go`).

### Notes
- Reference tests now distinguish between “implemented” and “verified.”
- Some methods use package-specific reporting conventions where multiple valid statistical conventions exist.
- `ProportionOneSample` and `ProportionTwoSample` use the null-hypothesis (score-test) standard error, i.e. `se = sqrt(p0*(1-p0)/n)`, matching R's `prop.test` rather than a Wald test built from the sample proportion's variance.
- `LevenesTest` uses the Brown-Forsythe formulation (deviations from the group median), matching SciPy's default (`center='median'`), rather than the original Levene formulation (deviations from the group mean).

## [0.1.0] - 2026-04-14

### Added
- Initial release of the `statstest` package.
- Core hypothesis testing support, including:
  - one-sample t-test
  - two-sample t-test
  - paired t-test
  - one-way ANOVA
  - Mann-Whitney U test
  - Wilcoxon signed-rank test
  - chi-square test of independence
  - Fisher’s exact test
- P-value adjustment support.
- Input validation for:
  - empty samples
  - small samples
  - NaN and infinite values
  - invalid alternatives
  - invalid confidence levels
  - invalid contingency tables
  - negative counts