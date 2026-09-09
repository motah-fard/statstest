package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// ProportionOneSample performs a one-sample z-test for the null hypothesis
// that a population proportion equals p0, given x successes out of n
// trials.
//
// The confidence interval is the Wald interval for the sample proportion,
// clipped to [0, 1].
func ProportionOneSample(x, n int, p0 float64, opts ProportionOneSampleOptions) (ProportionResult, error) {
	if err := validateCount(x, n); err != nil {
		return ProportionResult{}, err
	}
	if p0 <= 0 || p0 >= 1 {
		return ProportionResult{}, ErrInvalidProportion
	}
	if err := validateAlternative(opts.Alternative); err != nil {
		return ProportionResult{}, err
	}
	if err := validateConfidence(opts.ConfidenceLevel); err != nil {
		return ProportionResult{}, err
	}

	nf := float64(n)
	phat := float64(x) / nf

	se0 := math.Sqrt(p0 * (1 - p0) / nf)
	if se0 == 0 {
		return ProportionResult{}, ErrZeroVariance
	}

	d := phat - p0
	if opts.UseContinuityCorrection {
		d = continuityAdjust(d, 0.5/nf)
	}

	z := d / se0
	p := normalPValue(z, opts.Alternative)

	se := math.Sqrt(phat * (1 - phat) / nf)
	zcrit := distuv.Normal{Mu: 0, Sigma: 1}.Quantile(1 - (1-opts.ConfidenceLevel)/2)
	ci := ConfidenceInterval{
		Level: opts.ConfidenceLevel,
		Low:   math.Max(0, phat-zcrit*se),
		High:  math.Min(1, phat+zcrit*se),
	}

	return ProportionResult{
		Proportion:  phat,
		NullValue:   p0,
		Statistic:   z,
		PValue:      p,
		CI:          ci,
		Method:      "One-sample proportion z-test",
		Alternative: opts.Alternative,
	}, nil
}

// ProportionTwoSample performs a two-sample z-test for the null hypothesis
// that two population proportions are equal, given x1 successes out of n1
// trials and x2 successes out of n2 trials.
//
// The test statistic uses the pooled proportion under the null hypothesis.
// The confidence interval for the difference in proportions uses the
// unpooled (Wald) standard error.
func ProportionTwoSample(x1, n1, x2, n2 int, opts ProportionTwoSampleOptions) (ProportionTwoSampleResult, error) {
	if err := validateCount(x1, n1); err != nil {
		return ProportionTwoSampleResult{}, err
	}
	if err := validateCount(x2, n2); err != nil {
		return ProportionTwoSampleResult{}, err
	}
	if err := validateAlternative(opts.Alternative); err != nil {
		return ProportionTwoSampleResult{}, err
	}
	if err := validateConfidence(opts.ConfidenceLevel); err != nil {
		return ProportionTwoSampleResult{}, err
	}

	n1f := float64(n1)
	n2f := float64(n2)
	p1 := float64(x1) / n1f
	p2 := float64(x2) / n2f
	diff := p1 - p2

	pooled := float64(x1+x2) / (n1f + n2f)
	sePooled := math.Sqrt(pooled * (1 - pooled) * (1/n1f + 1/n2f))
	if sePooled == 0 {
		return ProportionTwoSampleResult{}, ErrZeroVariance
	}

	d := diff
	if opts.UseContinuityCorrection {
		d = continuityAdjust(d, 0.5*(1/n1f+1/n2f))
	}

	z := d / sePooled
	p := normalPValue(z, opts.Alternative)

	seUnpooled := math.Sqrt(p1*(1-p1)/n1f + p2*(1-p2)/n2f)
	zcrit := distuv.Normal{Mu: 0, Sigma: 1}.Quantile(1 - (1-opts.ConfidenceLevel)/2)
	ci := ConfidenceInterval{
		Level: opts.ConfidenceLevel,
		Low:   diff - zcrit*seUnpooled,
		High:  diff + zcrit*seUnpooled,
	}

	return ProportionTwoSampleResult{
		Proportion1: p1,
		Proportion2: p2,
		Diff:        diff,
		Statistic:   z,
		PValue:      p,
		CI:          ci,
		Method:      "Two-sample proportion z-test",
		Alternative: opts.Alternative,
	}, nil
}

// continuityAdjust shrinks d toward zero by amount cc, without crossing zero.
func continuityAdjust(d, cc float64) float64 {
	adj := math.Abs(d) - cc
	if adj < 0 {
		adj = 0
	}
	return math.Copysign(adj, d)
}
