package statstest

import "math"

// CohensD computes Cohen's d for two independent samples.
//
// If pooled is true, the pooled standard deviation is used.
// If pooled is false, the denominator is the square root of the average
// of the two sample variances.
func CohensD(x, y []float64, pooled bool) (float64, error) {
	if err := validateSample(x, 2); err != nil {
		return 0, err
	}
	if err := validateSample(y, 2); err != nil {
		return 0, err
	}

	meanX := mean(x)
	meanY := mean(y)
	varX := sampleVariance(x)
	varY := sampleVariance(y)

	var denom float64
	if pooled {
		pooledVar := pooledVariance(x, y)
		if pooledVar <= 0 {
			return 0, ErrZeroVariance
		}
		denom = math.Sqrt(pooledVar)
	} else {
		avgVar := (varX + varY) / 2.0
		if avgVar <= 0 {
			return 0, ErrZeroVariance
		}
		denom = math.Sqrt(avgVar)
	}

	return (meanX - meanY) / denom, nil
}

// HedgesG computes Hedges' g for two independent samples using
// Cohen's d with pooled standard deviation and a small-sample correction.
func HedgesG(x, y []float64) (float64, error) {
	if err := validateSample(x, 2); err != nil {
		return 0, err
	}
	if err := validateSample(y, 2); err != nil {
		return 0, err
	}

	d, err := CohensD(x, y, true)
	if err != nil {
		return 0, err
	}

	df := float64(len(x) + len(y) - 2)
	if df <= 0 {
		return 0, ErrSampleTooSmall
	}

	j := 1.0 - (3.0 / (4.0*df - 1.0))
	return j * d, nil
}
