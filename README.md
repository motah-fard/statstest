# statstest

[![Go Reference](https://pkg.go.dev/badge/github.com/motah-fard/statstest.svg)](https://pkg.go.dev/github.com/motah-fard/statstest)
[![Go Report Card](https://goreportcard.com/badge/github.com/motah-fard/statstest)](https://goreportcard.com/report/github.com/motah-fard/statstest)
[![License](https://img.shields.io/github/license/motah-fard/statstest)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/motah-fard/statstest)](go.mod)

`statstest` is a focused Go library for practical statistical inference.

It provides common hypothesis tests, effect sizes, confidence intervals,
and multiple-testing corrections using idiomatic Go APIs and standard slice types.

## Why this library exists

Go has strong numerical foundations, but practical inference workflows are still
more fragmented than they should be.

`statstest` fills that gap with a small, focused package for everyday
statistical testing in production, analytics, experimentation, and research tools.

## Goals

- simple APIs using `[]float64` and standard Go types
- reliable statistical test implementations
- structured results instead of raw tuples
- clear documentation of assumptions and interpretation
- no dataframe dependency

## Features

- one-sample t-test
- two-sample t-test
- Welch t-test
- paired t-test
- one-way ANOVA
- Mann-Whitney U test
- Wilcoxon signed-rank test
- chi-square test of independence
- Fisher’s exact test
- p-value adjustment methods
- confidence intervals and effect-size support where applicable
- input validation for invalid samples and tables

## Verification

This package includes reference-value tests for core statistical methods.
Expected values are checked against external statistical references and
documented package conventions for methods where multiple valid reporting
conventions exist.

## Installation

```bash
go get github.com/motah-fard/statstest


