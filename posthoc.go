package statstest

import "math"

// TukeyHSD performs Tukey's Honestly Significant Difference test: a
// post-hoc pairwise comparison of group means, typically run after a
// significant OneWayANOVA to find out which groups differ.
//
// Every pair of groups is compared using the Tukey-Kramer method, which
// handles unequal group sizes. conf is the simultaneous (family-wise)
// confidence level applied to every reported interval.
func TukeyHSD(conf float64, groups ...[]float64) (TukeyHSDResult, error) {
	if len(groups) < 2 {
		return TukeyHSDResult{}, ErrTooFewGroups
	}
	if err := validateConfidence(conf); err != nil {
		return TukeyHSDResult{}, err
	}

	means := make([]float64, len(groups))
	totalN := 0
	for i, g := range groups {
		if err := validateSample(g, 2); err != nil {
			return TukeyHSDResult{}, err
		}
		means[i] = mean(g)
		totalN += len(g)
	}

	k := len(groups)
	dfWithin := totalN - k
	if dfWithin <= 0 {
		return TukeyHSDResult{}, ErrInsufficientDF
	}

	var ssWithin float64
	for i, g := range groups {
		for _, v := range g {
			d := v - means[i]
			ssWithin += d * d
		}
	}
	msWithin := ssWithin / float64(dfWithin)
	if msWithin <= 0 {
		return TukeyHSDResult{}, ErrZeroVariance
	}

	qCrit := studentizedRangeQuantile(conf, k, float64(dfWithin))

	var comparisons []PairwiseComparison
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			se := math.Sqrt(msWithin / 2 * (1/float64(len(groups[i])) + 1/float64(len(groups[j]))))
			diff := means[i] - means[j]
			q := math.Abs(diff) / se
			p := 1 - studentizedRangeCDF(q, k, float64(dfWithin))

			comparisons = append(comparisons, PairwiseComparison{
				GroupI:    i,
				GroupJ:    j,
				MeanDiff:  diff,
				Statistic: q,
				PValue:    p,
				CI: ConfidenceInterval{
					Level: conf,
					Low:   diff - qCrit*se,
					High:  diff + qCrit*se,
				},
			})
		}
	}

	return TukeyHSDResult{
		Comparisons: comparisons,
		Method:      "Tukey HSD",
	}, nil
}

// DunnTest performs Dunn's post-hoc test: pairwise rank-sum comparisons
// between groups, typically run after a significant KruskalWallis test to
// find out which groups differ.
//
// Ranks are computed once over all groups pooled together, matching the
// ranking used by KruskalWallis, with an average-rank tie correction.
// Raw two-sided p-values are adjusted for multiple comparisons using
// method.
func DunnTest(method PAdjustMethod, groups ...[]float64) (DunnTestResult, error) {
	if len(groups) < 2 {
		return DunnTestResult{}, ErrTooFewGroups
	}
	for _, g := range groups {
		if err := validateSample(g, 1); err != nil {
			return DunnTestResult{}, err
		}
	}

	ranks, tieCounts, n := rankGroups(groups)
	if n <= len(groups) {
		return DunnTestResult{}, ErrInsufficientDF
	}

	nf := float64(n)
	tieCorrection := 0.0
	for _, tc := range tieCounts {
		tf := float64(tc)
		tieCorrection += tf*tf*tf - tf
	}
	sigma2 := nf*(nf+1)/12 - tieCorrection/(12*(nf-1))
	if sigma2 <= 0 {
		return DunnTestResult{}, ErrZeroVariance
	}

	rankMeans := make([]float64, len(groups))
	for gi, g := range groups {
		var sum float64
		for _, r := range ranks[gi] {
			sum += r
		}
		rankMeans[gi] = sum / float64(len(g))
	}

	var comparisons []DunnComparison
	var rawP []float64
	for i := 0; i < len(groups); i++ {
		for j := i + 1; j < len(groups); j++ {
			ni := float64(len(groups[i]))
			nj := float64(len(groups[j]))
			se := math.Sqrt(sigma2 * (1/ni + 1/nj))
			z := (rankMeans[i] - rankMeans[j]) / se
			p := normalPValue(z, TwoSided)

			comparisons = append(comparisons, DunnComparison{
				GroupI:    i,
				GroupJ:    j,
				Statistic: z,
				PValue:    p,
			})
			rawP = append(rawP, p)
		}
	}

	adjusted, err := AdjustPValues(rawP, method)
	if err != nil {
		return DunnTestResult{}, err
	}
	for idx := range comparisons {
		comparisons[idx].AdjustedPValue = adjusted[idx]
	}

	return DunnTestResult{
		Comparisons:  comparisons,
		AdjustMethod: method,
		Method:       "Dunn's test",
	}, nil
}
