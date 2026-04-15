package statstest

import "math"

// TTestOneSample performs a one-sample t-test for the null hypothesis
// that the population mean equals mu.
//
// The returned confidence interval is a two-sided confidence interval
// for the sample mean, regardless of the alternative hypothesis.
func TTestOneSample(x []float64, mu float64, alt Alternative, conf float64) (TTestOneSampleResult, error) {
	if err := validateSample(x, 2); err != nil {
		return TTestOneSampleResult{}, err
	}
	if err := validateAlternative(alt); err != nil {
		return TTestOneSampleResult{}, err
	}
	if err := validateConfidence(conf); err != nil {
		return TTestOneSampleResult{}, err
	}

	n := len(x)
	m := mean(x)
	s := sampleStdDev(x)
	if s == 0 {
		return TTestOneSampleResult{}, ErrZeroVariance
	}

	se := s / math.Sqrt(float64(n))
	diff := m - mu
	t := diff / se
	df := float64(n - 1)
	p := tPValue(t, df, alt)

	tcrit := tCriticalTwoSided(conf, df)
	ci := ConfidenceInterval{
		Level: conf,
		Low:   m - tcrit*se,
		High:  m + tcrit*se,
	}

	return TTestOneSampleResult{
		Statistic:   t,
		PValue:      p,
		DF:          df,
		Mean:        m,
		NullMean:    mu,
		MeanDiff:    diff,
		CI:          ci,
		Method:      "One-sample t-test",
		Alternative: alt,
	}, nil
}

// TTestTwoSample performs a two-sample t-test.
//
// If opts.EqualVariance is true, the pooled-variance t-test is used.
// Otherwise, Welch's t-test is used.
//
// The returned confidence interval is a two-sided confidence interval
// for the difference in means (mean(x) - mean(y)).
func TTestTwoSample(x, y []float64, opts TwoSampleTOptions) (TTestTwoSampleResult, error) {
	if err := validateSample(x, 2); err != nil {
		return TTestTwoSampleResult{}, err
	}
	if err := validateSample(y, 2); err != nil {
		return TTestTwoSampleResult{}, err
	}
	if err := validateAlternative(opts.Alternative); err != nil {
		return TTestTwoSampleResult{}, err
	}
	if err := validateConfidence(opts.ConfidenceLevel); err != nil {
		return TTestTwoSampleResult{}, err
	}

	n1 := float64(len(x))
	n2 := float64(len(y))
	m1 := mean(x)
	m2 := mean(y)
	v1 := sampleVariance(x)
	v2 := sampleVariance(y)

	diff := m1 - m2

	var se float64
	var df float64
	method := "Two-sample t-test (Welch)"

	if opts.EqualVariance {
		pv := pooledVariance(x, y)
		if pv <= 0 {
			return TTestTwoSampleResult{}, ErrZeroVariance
		}
		se = math.Sqrt(pv * (1/n1 + 1/n2))
		df = n1 + n2 - 2
		method = "Two-sample t-test (pooled)"
	} else {
		if v1 == 0 && v2 == 0 {
			return TTestTwoSampleResult{}, ErrZeroVariance
		}
		se2 := v1/n1 + v2/n2
		if se2 <= 0 {
			return TTestTwoSampleResult{}, ErrZeroVariance
		}
		se = math.Sqrt(se2)

		num := se2 * se2
		den := 0.0
		if v1 > 0 {
			den += (v1 / n1) * (v1 / n1) / (n1 - 1)
		}
		if v2 > 0 {
			den += (v2 / n2) * (v2 / n2) / (n2 - 1)
		}
		if den <= 0 {
			return TTestTwoSampleResult{}, ErrZeroVariance
		}
		df = num / den
	}

	t := diff / se
	p := tPValue(t, df, opts.Alternative)
	tcrit := tCriticalTwoSided(opts.ConfidenceLevel, df)

	ci := ConfidenceInterval{
		Level: opts.ConfidenceLevel,
		Low:   diff - tcrit*se,
		High:  diff + tcrit*se,
	}

	return TTestTwoSampleResult{
		Statistic:   t,
		PValue:      p,
		DF:          df,
		Mean1:       m1,
		Mean2:       m2,
		MeanDiff:    diff,
		CI:          ci,
		Method:      method,
		Alternative: opts.Alternative,
	}, nil
}

// TTestPaired performs a paired t-test by testing whether the mean
// of element-wise differences x[i] - y[i] equals zero.
//
// The returned confidence interval is a two-sided confidence interval
// for the mean paired difference.
func TTestPaired(x, y []float64, opts PairedTOptions) (TTestPairedResult, error) {
	if err := validateSample(x, 2); err != nil {
		return TTestPairedResult{}, err
	}
	if err := validateSample(y, 2); err != nil {
		return TTestPairedResult{}, err
	}
	if len(x) != len(y) {
		return TTestPairedResult{}, ErrMismatchedLengths
	}
	if err := validateAlternative(opts.Alternative); err != nil {
		return TTestPairedResult{}, err
	}
	if err := validateConfidence(opts.ConfidenceLevel); err != nil {
		return TTestPairedResult{}, err
	}

	diffs := make([]float64, len(x))
	for i := range x {
		diffs[i] = x[i] - y[i]
	}

	m := mean(diffs)
	s := sampleStdDev(diffs)
	if s == 0 {
		return TTestPairedResult{}, ErrZeroVariance
	}

	n := len(diffs)
	se := s / math.Sqrt(float64(n))
	t := m / se
	df := float64(n - 1)
	p := tPValue(t, df, opts.Alternative)

	tcrit := tCriticalTwoSided(opts.ConfidenceLevel, df)
	ci := ConfidenceInterval{
		Level: opts.ConfidenceLevel,
		Low:   m - tcrit*se,
		High:  m + tcrit*se,
	}

	return TTestPairedResult{
		Statistic:   t,
		PValue:      p,
		DF:          df,
		MeanDiff:    m,
		CI:          ci,
		Method:      "Paired t-test",
		Alternative: opts.Alternative,
	}, nil
}
