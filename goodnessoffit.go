package statstest

import (
	"gonum.org/v1/gonum/stat/distuv"
)

// ChiSquareGoodnessOfFit performs a chi-square goodness-of-fit test,
// comparing observed category counts to expected counts.
//
// If expected is nil, the expected counts are a uniform distribution of
// the total observed count across all categories. Otherwise expected must
// have the same length as observed and contain only positive values.
func ChiSquareGoodnessOfFit(observed, expected []float64) (ChiSquareGoodnessOfFitResult, error) {
	if err := validateSample(observed, 2); err != nil {
		return ChiSquareGoodnessOfFitResult{}, err
	}
	for _, v := range observed {
		if v < 0 {
			return ChiSquareGoodnessOfFitResult{}, ErrNegativeCount
		}
	}

	k := len(observed)

	exp := expected
	if exp == nil {
		total := 0.0
		for _, v := range observed {
			total += v
		}
		exp = make([]float64, k)
		for i := range exp {
			exp[i] = total / float64(k)
		}
	} else {
		if len(exp) != k {
			return ChiSquareGoodnessOfFitResult{}, ErrMismatchedLengths
		}
		for _, v := range exp {
			if v <= 0 {
				return ChiSquareGoodnessOfFitResult{}, ErrInvalidExpectedCounts
			}
		}
	}

	var chi2 float64
	for i := range observed {
		diff := observed[i] - exp[i]
		chi2 += diff * diff / exp[i]
	}

	df := k - 1
	if df <= 0 {
		return ChiSquareGoodnessOfFitResult{}, ErrInsufficientDF
	}

	dist := distuv.ChiSquared{K: float64(df)}
	p := 1 - dist.CDF(chi2)

	return ChiSquareGoodnessOfFitResult{
		Statistic: chi2,
		PValue:    p,
		DF:        df,
		Expected:  exp,
		Method:    "Chi-square goodness-of-fit test",
	}, nil
}
