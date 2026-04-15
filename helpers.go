package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

func validateSample(x []float64, minLen int) error {
	if len(x) == 0 {
		return ErrEmptySample
	}
	if len(x) < minLen {
		return ErrSampleTooSmall
	}
	for _, v := range x {
		if math.IsNaN(v) {
			return ErrContainsNaN
		}
		if math.IsInf(v, 0) {
			return ErrContainsInf
		}
	}
	return nil
}

func validateAlternative(alt Alternative) error {
	switch alt {
	case TwoSided, Less, Greater:
		return nil
	default:
		return ErrInvalidAlternative
	}
}

func validateConfidence(conf float64) error {
	if conf <= 0 || conf >= 1 {
		return ErrInvalidConfidenceLevel
	}
	return nil
}

func mean(x []float64) float64 {
	var sum float64
	for _, v := range x {
		sum += v
	}
	return sum / float64(len(x))
}

func sampleVariance(x []float64) float64 {
	if len(x) < 2 {
		return 0
	}
	m := mean(x)
	var ss float64
	for _, v := range x {
		d := v - m
		ss += d * d
	}
	return ss / float64(len(x)-1)
}

func sampleStdDev(x []float64) float64 {
	return math.Sqrt(sampleVariance(x))
}

func pooledVariance(x, y []float64) float64 {
	n1 := len(x)
	n2 := len(y)
	if n1 < 2 || n2 < 2 {
		return 0
	}

	v1 := sampleVariance(x)
	v2 := sampleVariance(y)

	num := float64(n1-1)*v1 + float64(n2-1)*v2
	den := float64(n1 + n2 - 2)
	if den == 0 {
		return 0
	}
	return num / den
}

func tPValue(t float64, df float64, alt Alternative) float64 {
	dist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    df,
	}

	switch alt {
	case TwoSided:
		if t < 0 {
			t = -t
		}
		return 2 * (1 - dist.CDF(t))
	case Less:
		return dist.CDF(t)
	case Greater:
		return 1 - dist.CDF(t)
	default:
		return math.NaN()
	}
}

func tCriticalTwoSided(conf float64, df float64) float64 {
	alpha := 1 - conf
	dist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    df,
	}
	return dist.Quantile(1 - alpha/2)
}

func normalPValue(z float64, alt Alternative) float64 {
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	switch alt {
	case TwoSided:
		if z < 0 {
			z = -z
		}
		return 2 * (1 - dist.CDF(z))
	case Less:
		return dist.CDF(z)
	case Greater:
		return 1 - dist.CDF(z)
	default:
		return math.NaN()
	}
}

func validateContingencyTable(observed [][]int) error {
	if len(observed) < 2 {
		return ErrInvalidTable
	}
	if len(observed[0]) < 2 {
		return ErrInvalidTable
	}

	cols := len(observed[0])
	for i := range observed {
		if len(observed[i]) != cols {
			return ErrInvalidTable
		}
		for _, v := range observed[i] {
			if v < 0 {
				return ErrNegativeCount
			}
		}
	}

	return nil
}

func validate2x2Table(table [2][2]int) error {
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if table[i][j] < 0 {
				return ErrNegativeCount
			}
		}
	}

	total := 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			total += table[i][j]
		}
	}
	if total == 0 {
		return ErrInvalidTable
	}

	return nil
}

func logChoose(n, k int) float64 {
	if k < 0 || k > n {
		return math.Inf(-1)
	}
	lnN, _ := math.Lgamma(float64(n + 1))
	lnK, _ := math.Lgamma(float64(k + 1))
	lnNK, _ := math.Lgamma(float64(n - k + 1))
	return lnN - lnK - lnNK
}

// hypergeomProb returns the hypergeometric probability for a 2x2 table
// defined by cell count a and fixed margins.
func hypergeomProb(a, r1, c1, n int) float64 {
	c2 := n - c1
	logP := logChoose(c1, a) + logChoose(c2, r1-a) - logChoose(n, r1)
	return math.Exp(logP)
}
