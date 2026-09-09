package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// LevenesTest performs Levene's test for equality of variances across two
// or more groups, using the Brown-Forsythe formulation (deviations from
// each group's median), which is more robust to non-normal data than the
// classic mean-based formulation.
//
// The null hypothesis is that all groups have equal variance. This is
// useful for choosing between the pooled and Welch forms of
// TTestTwoSample, or for checking the equal-variance assumption before a
// one-way ANOVA.
func LevenesTest(groups ...[]float64) (LeveneResult, error) {
	if len(groups) < 2 {
		return LeveneResult{}, ErrTooFewGroups
	}

	deviations := make([][]float64, len(groups))
	totalN := 0
	for i, g := range groups {
		if err := validateSample(g, 2); err != nil {
			return LeveneResult{}, err
		}
		center := median(g)
		d := make([]float64, len(g))
		for j, v := range g {
			d[j] = math.Abs(v - center)
		}
		deviations[i] = d
		totalN += len(g)
	}

	k := len(groups)
	dfBetween := k - 1
	dfWithin := totalN - k
	if dfBetween <= 0 || dfWithin <= 0 {
		return LeveneResult{}, ErrInsufficientDF
	}

	grandMean := 0.0
	for _, d := range deviations {
		for _, v := range d {
			grandMean += v
		}
	}
	grandMean /= float64(totalN)

	var ssBetween, ssWithin float64
	for _, d := range deviations {
		groupMean := mean(d)
		n := float64(len(d))

		diffGroup := groupMean - grandMean
		ssBetween += n * diffGroup * diffGroup

		for _, v := range d {
			diffWithin := v - groupMean
			ssWithin += diffWithin * diffWithin
		}
	}

	msBetween := ssBetween / float64(dfBetween)
	msWithin := ssWithin / float64(dfWithin)

	var w float64
	switch {
	case msWithin > 0:
		w = msBetween / msWithin
	case msBetween > 0:
		w = math.Inf(1)
	}

	dist := distuv.F{D1: float64(dfBetween), D2: float64(dfWithin)}
	p := 1 - dist.CDF(w)
	if math.IsNaN(p) {
		p = 0
	}

	return LeveneResult{
		Statistic: w,
		PValue:    p,
		DFBetween: dfBetween,
		DFWithin:  dfWithin,
		Method:    "Levene's test (Brown-Forsythe)",
	}, nil
}
