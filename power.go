package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// CohensH computes Cohen's h, the arcsine-transformed effect size for
// comparing two proportions. It is the proportion analogue of Cohen's d
// and is what PowerProportionTwoSample and SampleSizeProportionTwoSample
// use internally.
func CohensH(p1, p2 float64) (float64, error) {
	if p1 <= 0 || p1 >= 1 {
		return 0, ErrInvalidProportion
	}
	if p2 <= 0 || p2 >= 1 {
		return 0, ErrInvalidProportion
	}
	return 2*math.Asin(math.Sqrt(p1)) - 2*math.Asin(math.Sqrt(p2)), nil
}

// PowerTTestTwoSample computes the achieved statistical power of a
// two-sided two-sample t-test with n observations per group, given
// Cohen's d and significance level alpha.
//
// This uses the normal approximation to the sampling distribution
// (matching statsmodels' NormalIndPower), which is standard practice and
// close to exact once n is roughly 20 or more per group. For an exact
// small-sample calculation via the noncentral t-distribution, use a tool
// built around that distribution instead — this package does not
// implement one.
func PowerTTestTwoSample(n int, d, alpha float64) (float64, error) {
	if n < 2 {
		return 0, ErrSampleTooSmall
	}
	if d == 0 {
		return 0, ErrZeroEffectSize
	}
	if err := validateConfidence(alpha); err != nil {
		return 0, err
	}

	ncp := math.Abs(d) * math.Sqrt(float64(n)/2)
	return normalApproxPower(ncp, alpha), nil
}

// SampleSizeTTestTwoSample computes the required sample size per group
// for a two-sided two-sample t-test to achieve the target power, given
// Cohen's d and significance level alpha. See PowerTTestTwoSample for the
// approximation this is built on.
//
// The result is a continuous value; round up to the next integer to
// choose an actual sample size.
func SampleSizeTTestTwoSample(d, alpha, power float64) (float64, error) {
	if d == 0 {
		return 0, ErrZeroEffectSize
	}
	if err := validateConfidence(alpha); err != nil {
		return 0, err
	}
	if err := validateConfidence(power); err != nil {
		return 0, err
	}

	ncp := solveNCPForPower(power, alpha)
	absD := math.Abs(d)
	return 2 * (ncp / absD) * (ncp / absD), nil
}

// PowerProportionTwoSample computes the achieved statistical power of a
// two-sided two-sample proportion (z) test with n observations per group,
// given the two population proportions and significance level alpha.
//
// Like PowerTTestTwoSample, this uses the normal approximation to the
// sampling distribution, applied to Cohen's h (see CohensH) in place of
// Cohen's d.
func PowerProportionTwoSample(n int, p1, p2, alpha float64) (float64, error) {
	if n < 2 {
		return 0, ErrSampleTooSmall
	}
	h, err := CohensH(p1, p2)
	if err != nil {
		return 0, err
	}
	if h == 0 {
		return 0, ErrZeroEffectSize
	}
	if err := validateConfidence(alpha); err != nil {
		return 0, err
	}

	ncp := math.Abs(h) * math.Sqrt(float64(n)/2)
	return normalApproxPower(ncp, alpha), nil
}

// SampleSizeProportionTwoSample computes the required sample size per
// group for a two-sided two-sample proportion (z) test to achieve the
// target power, given the two population proportions and significance
// level alpha.
//
// The result is a continuous value; round up to the next integer to
// choose an actual sample size.
func SampleSizeProportionTwoSample(p1, p2, alpha, power float64) (float64, error) {
	h, err := CohensH(p1, p2)
	if err != nil {
		return 0, err
	}
	if h == 0 {
		return 0, ErrZeroEffectSize
	}
	if err := validateConfidence(alpha); err != nil {
		return 0, err
	}
	if err := validateConfidence(power); err != nil {
		return 0, err
	}

	ncp := solveNCPForPower(power, alpha)
	absH := math.Abs(h)
	return 2 * (ncp / absH) * (ncp / absH), nil
}

// normalApproxPower returns the power of a two-sided test with
// noncentrality parameter ncp (effect size scaled by sample size) and
// significance level alpha, under the normal approximation:
//
//	power = Phi(ncp - z) + Phi(-ncp - z),  z = Phi^-1(1 - alpha/2)
func normalApproxPower(ncp, alpha float64) float64 {
	norm := distuv.Normal{Mu: 0, Sigma: 1}
	zCrit := norm.Quantile(1 - alpha/2)
	return clampProbability(norm.CDF(ncp-zCrit) + norm.CDF(-ncp-zCrit))
}

// solveNCPForPower inverts normalApproxPower for ncp by bisection. Unlike
// studentizedRangeQuantile, normalApproxPower is a closed-form expression
// (no numerical integration), so a generous iteration count costs
// essentially nothing.
func solveNCPForPower(power, alpha float64) float64 {
	lo, hi := 0.0, 1.0
	for normalApproxPower(hi, alpha) < power {
		hi *= 2
		if hi > 1e6 {
			break
		}
	}
	for i := 0; i < 100; i++ {
		mid := (lo + hi) / 2
		if normalApproxPower(mid, alpha) < power {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}
