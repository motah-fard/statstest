# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project follows Semantic Versioning while the API is still evolving.

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