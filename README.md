# statstest

[![Go Reference](https://pkg.go.dev/badge/github.com/motah-fard/statstest.svg)](https://pkg.go.dev/github.com/motah-fard/statstest)
[![Go Report Card](https://goreportcard.com/badge/github.com/motah-fard/statstest)](https://goreportcard.com/report/github.com/motah-fard/statstest)
[![License](https://img.shields.io/github/license/motah-fard/statstest)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/motah-fard/statstest)](go.mod)

`statstest` brings R's `t.test`, `wilcox.test`, `chisq.test`, `cor.test`,
and `aov` to Go: hypothesis tests, correlations, effect sizes, confidence
intervals, and multiple-testing corrections, with structured results
instead of raw tuples.

```go
treatment := []float64{12.1, 13.4, 11.9, 14.0, 12.8}
control := []float64{10.2, 10.8, 11.3, 9.9, 10.5}

res, _ := statstest.TTestTwoSample(treatment, control, statstest.TwoSampleTOptions{
    EqualVariance:   false,
    Alternative:     statstest.TwoSided,
    ConfidenceLevel: 0.95,
})

fmt.Printf("t = %.3f, p = %.3f, 95%% CI = [%.3f, %.3f]\n",
    res.Statistic, res.PValue, res.CI.Low, res.CI.High)
// t = 4.981, p = 0.002, 95% CI = [1.196, 3.404]
```

## Why statstest

Go's `gonum.org/v1/gonum/stat` gives you distributions and summary
statistics; it does not give you a hypothesis test. Calling `TTestTwoSample`
here returns a `TTestTwoSampleResult{Statistic, PValue, DF, CI, Method, ...}`
instead of a bare `float64`, and `TTestOneSample(x, mu, alt, 1.5)` returns
`ErrInvalidConfidenceLevel` instead of silently producing `NaN`. If you're
porting an analysis out of R or Python, `statstest`'s functions map
directly onto their equivalents:

| Task                        | R                  | statstest                |
| --------------------------- | ------------------ | ------------------------- |
| One/two-sample t-test       | `t.test()`         | `TTestOneSample`, `TTestTwoSample` |
| Paired t-test                | `t.test(paired=TRUE)` | `TTestPaired`          |
| One-way ANOVA                | `aov()`             | `OneWayANOVA`             |
| Tukey HSD post-hoc            | `TukeyHSD()`        | `TukeyHSD`                |
| Mann-Whitney / Wilcoxon rank-sum | `wilcox.test()` | `MannWhitneyU`             |
| Wilcoxon signed-rank          | `wilcox.test(paired=TRUE)` | `WilcoxonSignedRank` |
| Kruskal-Wallis                | `kruskal.test()`    | `KruskalWallis`           |
| Dunn's post-hoc test          | `dunn.test()`       | `DunnTest`                |
| Chi-square test of independence | `chisq.test()`   | `ChiSquareIndependence`    |
| Chi-square goodness-of-fit     | `chisq.test(x, p=...)` | `ChiSquareGoodnessOfFit` |
| Fisher's exact test            | `fisher.test()`    | `FishersExact2x2`         |
| Pearson / Spearman correlation | `cor.test()`       | `PearsonCorrelation`, `SpearmanCorrelation` |
| One/two-sample proportion test | `prop.test()`     | `ProportionOneSample`, `ProportionTwoSample` |
| Levene's test (equal variances) | `car::leveneTest()` | `LevenesTest`            |
| Shapiro-Wilk normality test    | `shapiro.test()`   | `ShapiroWilk`             |
| Multiple-testing correction    | `p.adjust()`       | `AdjustPValues`           |

## Which test should I use?

**Comparing two groups' means?**
- Paired / before-after measurements on the same subjects → check normality of the differences with `ShapiroWilk`, then use `TTestPaired` (normal) or `WilcoxonSignedRank` (not normal).
- Independent groups → check normality (`ShapiroWilk`) and equal variance (`LevenesTest`) first.
  - Normal, equal variance → `TTestTwoSample` with `EqualVariance: true`.
  - Normal, unequal variance → `TTestTwoSample` with `EqualVariance: false` (Welch, the safer default when unsure).
  - Not normal → `MannWhitneyU`.

**Comparing three or more groups' means?**
- Normal data, roughly equal variances → `OneWayANOVA`. If the overall p-value is significant, follow up with `TukeyHSD` to see which specific pairs of groups differ.
- Not normal, or ordinal data → `KruskalWallis`. Follow up a significant result with `DunnTest`.

**Is a sample normally distributed?** → `ShapiroWilk`. Most parametric tests above (t-tests, ANOVA) assume this; check it before trusting them on small samples.

**Do two continuous variables move together?**
- Roughly linear relationship, no extreme outliers → `PearsonCorrelation`.
- Monotonic but not necessarily linear, or outliers/ranks → `SpearmanCorrelation`.

**Working with counts or categories?**
- Is a categorical variable associated with another? → `ChiSquareIndependence` (or `FishersExact2x2` for a 2x2 table with small counts, where the chi-square approximation is unreliable).
- Do observed category counts match an expected distribution? → `ChiSquareGoodnessOfFit`.
- Comparing a conversion rate/proportion to a target, or two proportions against each other (e.g. A/B test results) → `ProportionOneSample` / `ProportionTwoSample`.

**Running many tests at once (e.g. many A/B metrics, or post-hoc pairs)?** → adjust the resulting p-values with `AdjustPValues` before deciding what's significant. `TukeyHSD` and `DunnTest` already build this in for their pairwise comparisons.

## Goals

- simple APIs using `[]float64` and standard Go types
- reliable statistical test implementations
- structured results instead of raw tuples
- clear documentation of assumptions and interpretation
- no dataframe dependency

## Features

- one-sample, two-sample, Welch, and paired t-tests
- one-way ANOVA, with Tukey HSD post-hoc pairwise comparisons
- Mann-Whitney U test, Wilcoxon signed-rank test, Kruskal-Wallis test,
  with Dunn's post-hoc pairwise comparisons
- chi-square test of independence, chi-square goodness-of-fit, Fisher's exact test
- Pearson and Spearman correlation
- one-sample and two-sample proportion (z) tests
- Levene's test for equal variances, Shapiro-Wilk test for normality
- p-value adjustment (Bonferroni, Holm, Benjamini-Hochberg)
- effect sizes (Cohen's d, Hedges' g) and confidence intervals
- input validation for invalid samples, counts, and tables

## Verification

This package includes reference-value tests for core statistical methods.
Expected values are checked against external statistical references
(R, SciPy, statsmodels, scikit-posthocs) and documented package
conventions for methods where multiple valid reporting conventions exist.
Tukey HSD's underlying studentized range distribution is implemented from
scratch (via numerical integration, since neither Go's standard library
nor gonum provide it) and independently verified against
`scipy.stats.studentized_range`.

## Installation

```bash
go get github.com/motah-fard/statstest
```
