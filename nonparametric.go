package statstest

import "math"

// MannWhitneyU performs the Mann-Whitney U test for two independent samples.
//
// The test uses average ranks for ties and a normal approximation for the p-value.
// For the reported U statistic, this function returns U for the first sample x.
func MannWhitneyU(x, y []float64, opts MannWhitneyOptions) (MannWhitneyResult, error) {
	if err := validateSample(x, 1); err != nil {
		return MannWhitneyResult{}, err
	}
	if err := validateSample(y, 1); err != nil {
		return MannWhitneyResult{}, err
	}
	if err := validateAlternative(opts.Alternative); err != nil {
		return MannWhitneyResult{}, err
	}

	n1 := float64(len(x))
	n2 := float64(len(y))
	if n1 == 0 || n2 == 0 {
		return MannWhitneyResult{}, ErrEmptySample
	}

	ranked, tieCounts := rankWithAverageTies(x, y)
	r1 := sumRanksForX(ranked)

	u1 := r1 - (n1*(n1+1))/2.0
	mu := n1 * n2 / 2.0

	N := n1 + n2
	tieCorrection := 0.0
	for _, t := range tieCounts {
		tf := float64(t)
		tieCorrection += tf*tf*tf - tf
	}

	varU := (n1 * n2 / 12.0) * (N + 1.0)
	if tieCorrection > 0 && N > 1 {
		varU -= (n1 * n2 / 12.0) * (tieCorrection / (N * (N - 1.0)))
	}

	if varU <= 0 {
		return MannWhitneyResult{}, ErrZeroVariance
	}

	sdU := math.Sqrt(varU)

	cc := 0.0
	if opts.UseContinuityCorrection {
		switch opts.Alternative {
		case TwoSided:
			if u1 > mu {
				cc = 0.5
			} else if u1 < mu {
				cc = -0.5
			}
		case Greater:
			cc = 0.5
		case Less:
			cc = -0.5
		}
	}

	z := (u1 - mu - cc) / sdU
	p := normalPValue(z, opts.Alternative)

	return MannWhitneyResult{
		U:           u1,
		PValue:      p,
		Z:           z,
		Method:      "Mann-Whitney U test",
		Alternative: opts.Alternative,
	}, nil
}

// WilcoxonSignedRank performs the Wilcoxon signed-rank test for paired samples.
//
// Zero paired differences are removed before ranking. The test uses average ranks
// for ties in absolute differences and a normal approximation for the p-value.
//
// The reported W statistic is the sum of positive ranks.
func WilcoxonSignedRank(x, y []float64, alt Alternative) (WilcoxonSignedRankResult, error) {
	if err := validateSample(x, 2); err != nil {
		return WilcoxonSignedRankResult{}, err
	}
	if err := validateSample(y, 2); err != nil {
		return WilcoxonSignedRankResult{}, err
	}
	if len(x) != len(y) {
		return WilcoxonSignedRankResult{}, ErrMismatchedLengths
	}
	if err := validateAlternative(alt); err != nil {
		return WilcoxonSignedRankResult{}, err
	}

	diffs := make([]float64, len(x))
	for i := range x {
		diffs[i] = x[i] - y[i]
	}

	items, tieCounts := signedRanksFromDifferences(diffs)
	if len(items) == 0 {
		return WilcoxonSignedRankResult{}, ErrNoNonZeroDifferences
	}

	n := float64(len(items))
	wPlus, _ := sumPositiveNegativeRanks(items)

	mu := n * (n + 1) / 4.0
	varW := n * (n + 1) * (2*n + 1) / 24.0

	tieCorrection := 0.0
	for _, t := range tieCounts {
		tf := float64(t)
		tieCorrection += tf*tf*tf - tf
	}
	if tieCorrection > 0 {
		varW -= tieCorrection / 48.0
	}

	if varW <= 0 {
		return WilcoxonSignedRankResult{}, ErrZeroVariance
	}

	z := (wPlus - mu) / math.Sqrt(varW)
	p := normalPValue(z, alt)

	return WilcoxonSignedRankResult{
		W:           wPlus,
		PValue:      p,
		Z:           z,
		Method:      "Wilcoxon signed-rank test",
		Alternative: alt,
	}, nil
}
