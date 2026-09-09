package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// studentizedRangeCDF returns P(Q <= q) for the studentized range
// distribution with k means and df error degrees of freedom, computed by
// numerical integration over the distribution of the pooled standard
// deviation estimate.
func studentizedRangeCDF(q float64, k int, df float64) float64 {
	if q <= 0 {
		return 0
	}
	if k < 2 {
		return 1
	}
	if math.IsInf(df, 1) || df > 1e6 {
		return normalRangeCDF(q, k)
	}

	chi := distuv.ChiSquared{K: df}
	upper := chi.Quantile(1 - 1e-12)

	f := func(v float64) float64 {
		if v <= 0 {
			return 0
		}
		s := math.Sqrt(v / df)
		return chi.Prob(v) * normalRangeCDF(q*s, k)
	}

	return clampProbability(simpsonIntegrate(f, 1e-10, upper, 128))
}

// clampProbability clamps a numerically integrated probability to [0, 1],
// guarding against small overshoot from quadrature error near the tails.
func clampProbability(p float64) float64 {
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}

// normalRangeCDF returns P(R <= r), the CDF of the range of k iid standard
// normal variables: k * integral of phi(z) * [Phi(z) - Phi(z-r)]^(k-1) dz.
func normalRangeCDF(r float64, k int) float64 {
	if r <= 0 {
		return 0
	}
	norm := distuv.Normal{Mu: 0, Sigma: 1}
	f := func(z float64) float64 {
		diff := norm.CDF(z) - norm.CDF(z-r)
		if diff <= 0 {
			return 0
		}
		return norm.Prob(z) * math.Pow(diff, float64(k-1))
	}
	return clampProbability(float64(k) * simpsonIntegrate(f, -12, 12, 128))
}

// studentizedRangeQuantile inverts studentizedRangeCDF by bisection.
func studentizedRangeQuantile(p float64, k int, df float64) float64 {
	lo, hi := 0.0, 1.0
	for studentizedRangeCDF(hi, k, df) < p {
		hi *= 2
		if hi > 1e6 {
			break
		}
	}
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if studentizedRangeCDF(mid, k, df) < p {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

// simpsonIntegrate approximates the integral of f over [a, b] using
// composite Simpson's rule with n subintervals (rounded up to even).
func simpsonIntegrate(f func(float64) float64, a, b float64, n int) float64 {
	if n%2 != 0 {
		n++
	}
	h := (b - a) / float64(n)
	sum := f(a) + f(b)
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * f(x)
		} else {
			sum += 4 * f(x)
		}
	}
	return sum * h / 3
}
