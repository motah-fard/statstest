package statstest

import (
	"gonum.org/v1/gonum/stat/distuv"
)

// KruskalWallis performs the Kruskal-Wallis H test, the nonparametric
// analogue of one-way ANOVA, across two or more independent groups.
//
// The null hypothesis is that all groups are drawn from the same
// distribution. The test uses average ranks for ties and a chi-square
// approximation, corrected for ties, for the p-value.
func KruskalWallis(groups ...[]float64) (KruskalWallisResult, error) {
	if len(groups) < 2 {
		return KruskalWallisResult{}, ErrTooFewGroups
	}
	for _, g := range groups {
		if err := validateSample(g, 1); err != nil {
			return KruskalWallisResult{}, err
		}
	}

	ranks, tieCounts, n := rankGroups(groups)
	if n <= len(groups) {
		return KruskalWallisResult{}, ErrInsufficientDF
	}

	var h float64
	for gi, g := range groups {
		var rankSum float64
		for _, r := range ranks[gi] {
			rankSum += r
		}
		h += (rankSum * rankSum) / float64(len(g))
	}

	nf := float64(n)
	h = (12/(nf*(nf+1)))*h - 3*(nf+1)

	tieCorrection := 0.0
	for _, tc := range tieCounts {
		tf := float64(tc)
		tieCorrection += tf*tf*tf - tf
	}
	if denom := nf*nf*nf - nf; tieCorrection > 0 && denom > 0 {
		h /= 1 - tieCorrection/denom
	}

	df := len(groups) - 1
	dist := distuv.ChiSquared{K: float64(df)}
	p := 1 - dist.CDF(h)

	return KruskalWallisResult{
		H:      h,
		PValue: p,
		DF:     df,
		Method: "Kruskal-Wallis test",
	}, nil
}
