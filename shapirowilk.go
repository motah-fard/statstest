package statstest

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat/distuv"
)

// Coefficient polynomials for Royston's Shapiro-Wilk approximation
// (Algorithm AS R94, Applied Statistics 1995, vol. 44, no. 4, 547-551),
// the same algorithm used by R's shapiro.test.
var (
	swG  = []float64{-2.273, 0.459}
	swC1 = []float64{0, 0.221157, -0.147981, -2.071190, 4.434685, -2.706056}
	swC2 = []float64{0, 0.042981, -0.293762, -1.752461, 5.682633, -3.582633}
	swC3 = []float64{0.544, -0.39978, 0.025054, -6.714e-4}
	swC4 = []float64{1.3822, -0.77857, 0.062767, -0.0020322}
	swC5 = []float64{-1.5861, -0.31082, -0.083751, 0.0038915}
	swC6 = []float64{-0.4803, -0.082676, 0.0030302}
)

// ShapiroWilk performs the Shapiro-Wilk test for normality.
//
// The null hypothesis is that the sample was drawn from a normal
// distribution. This implementation follows Royston's algorithm AS R94
// (Applied Statistics, 1995), the same algorithm used by R's
// shapiro.test, and supports sample sizes from 3 to 5000.
func ShapiroWilk(x []float64) (ShapiroWilkResult, error) {
	if err := validateSample(x, 3); err != nil {
		return ShapiroWilkResult{}, err
	}
	n := len(x)
	if n > 5000 {
		return ShapiroWilkResult{}, ErrSampleTooLarge
	}

	sorted := append([]float64(nil), x...)
	sort.Float64s(sorted)

	rng := sorted[n-1] - sorted[0]
	if rng <= 0 {
		return ShapiroWilkResult{}, ErrZeroVariance
	}

	w := shapiroWilkStatistic(sorted, rng)
	p := shapiroWilkPValue(w, n)

	return ShapiroWilkResult{
		W:      w,
		PValue: p,
		Method: "Shapiro-Wilk test",
	}, nil
}

// shapiroWilkCoefficients computes the 1-based (index 0 unused) vector of
// the first n/2 Shapiro-Wilk weights for a sample of size n.
func shapiroWilkCoefficients(n int) []float64 {
	nn2 := n / 2
	a := make([]float64, nn2+1)

	if n == 3 {
		a[1] = 0.70710678
		return a
	}

	an := float64(n)
	an25 := an + 0.25
	norm := distuv.Normal{Mu: 0, Sigma: 1}

	summ2 := 0.0
	for i := 1; i <= nn2; i++ {
		a[i] = norm.Quantile((float64(i) - 0.375) / an25)
		summ2 += a[i] * a[i]
	}
	summ2 *= 2
	ssumm2 := math.Sqrt(summ2)
	rsn := 1 / math.Sqrt(an)
	a1 := polyEval(swC1, rsn) - a[1]/ssumm2

	var i1 int
	var fac float64
	if n > 5 {
		i1 = 3
		rawA2 := a[2]
		a2 := -rawA2/ssumm2 + polyEval(swC2, rsn)
		fac = math.Sqrt((summ2 - 2*a[1]*a[1] - 2*rawA2*rawA2) / (1 - 2*a1*a1 - 2*a2*a2))
		a[2] = a2
	} else {
		i1 = 2
		fac = math.Sqrt((summ2 - 2*a[1]*a[1]) / (1 - 2*a1*a1))
	}
	a[1] = a1
	for i := i1; i <= nn2; i++ {
		a[i] /= -fac
	}

	return a
}

// shapiroWilkStatistic computes W as the squared correlation between the
// sorted, range-scaled sample and the Shapiro-Wilk coefficients.
func shapiroWilkStatistic(sorted []float64, rng float64) float64 {
	n := len(sorted)
	a := shapiroWilkCoefficients(n)

	// sa and sx are near-zero correction terms (the a[] and x[] values are
	// each nearly antisymmetric/centered) used to reduce rounding error in
	// the correlation computed below, following the reference algorithm.
	sx := sorted[0] / rng
	sa := -a[1]

	i, j := 1, n-1
	for i < n {
		xi := sorted[i] / rng
		sx += xi
		i++
		if i != j {
			sa += swSign(i-j) * a[minInt(i, j)]
		}
		j--
	}
	sa /= float64(n)
	sx /= float64(n)

	var ssa, ssx, sax float64
	i, j = 0, n-1
	for i < n {
		var asa float64
		if i != j {
			asa = swSign(i-j)*a[1+minInt(i, j)] - sa
		} else {
			asa = -sa
		}
		xsx := sorted[i]/rng - sx
		ssa += asa * asa
		ssx += xsx * xsx
		sax += asa * xsx
		i++
		j--
	}

	ssassx := math.Sqrt(ssa * ssx)
	w1 := (ssassx - sax) * (ssassx + sax) / (ssa * ssx)
	return 1 - w1
}

// shapiroWilkPValue converts the W statistic to a p-value using Royston's
// approximation, which differs for n == 3, 4 <= n <= 11, and n >= 12.
func shapiroWilkPValue(w float64, n int) float64 {
	if n == 3 {
		const pi6 = 1.90985931710274  // 6 / pi
		const stqr = 1.04719755119660 // asin(sqrt(3/4))
		p := pi6 * (math.Asin(math.Sqrt(w)) - stqr)
		if p < 0 {
			p = 0
		}
		return p
	}

	an := float64(n)
	y := math.Log(1 - w)
	norm := distuv.Normal{Mu: 0, Sigma: 1}

	var m, s float64
	if n <= 11 {
		gamma := polyEval(swG, an)
		if y >= gamma {
			return 1e-99
		}
		y = -math.Log(gamma - y)
		m = polyEval(swC3, an)
		s = math.Exp(polyEval(swC4, an))
	} else {
		xx := math.Log(an)
		m = polyEval(swC5, xx)
		s = math.Exp(polyEval(swC6, xx))
	}

	return 1 - norm.CDF((y-m)/s)
}

// polyEval evaluates cc[0] + cc[1]*x + cc[2]*x^2 + ... via Horner's method.
func polyEval(cc []float64, x float64) float64 {
	result := 0.0
	for i := len(cc) - 1; i >= 0; i-- {
		result = result*x + cc[i]
	}
	return result
}

func swSign(v int) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
