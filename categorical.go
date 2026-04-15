package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// ChiSquareIndependence performs a chi-square test of independence
// on a contingency table of observed counts.
func ChiSquareIndependence(observed [][]int) (ChiSquareResult, error) {
	if err := validateContingencyTable(observed); err != nil {
		return ChiSquareResult{}, err
	}

	rows := len(observed)
	cols := len(observed[0])

	rowSums := make([]float64, rows)
	colSums := make([]float64, cols)
	var total float64

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := float64(observed[i][j])
			rowSums[i] += v
			colSums[j] += v
			total += v
		}
	}

	if total <= 0 {
		return ChiSquareResult{}, ErrInvalidTable
	}

	expected := make([][]float64, rows)
	for i := range expected {
		expected[i] = make([]float64, cols)
	}

	var chi2 float64
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			e := rowSums[i] * colSums[j] / total
			expected[i][j] = e
			if e <= 0 {
				return ChiSquareResult{}, ErrZeroVariance
			}

			o := float64(observed[i][j])
			diff := o - e
			chi2 += diff * diff / e
		}
	}

	df := (rows - 1) * (cols - 1)
	if df <= 0 {
		return ChiSquareResult{}, ErrInsufficientDF
	}

	dist := distuv.ChiSquared{K: float64(df)}
	p := 1 - dist.CDF(chi2)

	return ChiSquareResult{
		Statistic: chi2,
		PValue:    p,
		DF:        df,
		Expected:  expected,
		Method:    "Chi-square test of independence",
	}, nil
}

// FishersExact2x2 performs Fisher's exact test for a 2x2 contingency table.
//
// The reported odds ratio is the sample odds ratio a*d / (b*c). If b*c == 0
// and a*d > 0, the odds ratio is reported as +Inf. If both numerator and
// denominator are zero, the odds ratio is reported as NaN.
func FishersExact2x2(table [2][2]int, alt Alternative) (FishersExactResult, error) {
	if err := validate2x2Table(table); err != nil {
		return FishersExactResult{}, err
	}
	if err := validateAlternative(alt); err != nil {
		return FishersExactResult{}, err
	}

	a := table[0][0]
	b := table[0][1]
	c := table[1][0]
	d := table[1][1]

	r1 := a + b
	r2 := c + d
	c1 := a + c
	_ = r2
	n := r1 + r2

	minA := 0
	if r1-(n-c1) > minA {
		minA = r1 - (n - c1)
	}
	maxA := r1
	if c1 < maxA {
		maxA = c1
	}

	pObs := hypergeomProb(a, r1, c1, n)

	var p float64
	switch alt {
	case Less:
		for x := minA; x <= a; x++ {
			p += hypergeomProb(x, r1, c1, n)
		}
	case Greater:
		for x := a; x <= maxA; x++ {
			p += hypergeomProb(x, r1, c1, n)
		}
	case TwoSided:
		const eps = 1e-12
		for x := minA; x <= maxA; x++ {
			px := hypergeomProb(x, r1, c1, n)
			if px <= pObs+eps {
				p += px
			}
		}
	default:
		return FishersExactResult{}, ErrInvalidAlternative
	}

	if p > 1 {
		p = 1
	}

	num := float64(a * d)
	den := float64(b * c)

	oddsRatio := math.NaN()
	switch {
	case den == 0 && num > 0:
		oddsRatio = math.Inf(1)
	case den == 0 && num == 0:
		oddsRatio = math.NaN()
	default:
		oddsRatio = num / den
	}

	return FishersExactResult{
		OddsRatio:   oddsRatio,
		PValue:      p,
		Method:      "Fisher's exact test",
		Alternative: alt,
	}, nil
}
