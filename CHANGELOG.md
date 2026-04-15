# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog,
and this project follows Semantic Versioning while the API is still evolving.

## [Unreleased]

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

### Changed
- Improved statistical verification by checking implementation outputs against external reference values and documented package conventions.
- Updated convention-sensitive reference tests for:
  - Mann-Whitney U asymptotic p-value behavior
  - Wilcoxon signed-rank statistic reporting convention (`W+`)

### Notes
- Reference tests now distinguish between “implemented” and “verified.”
- Some methods use package-specific reporting conventions where multiple valid statistical conventions exist.

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