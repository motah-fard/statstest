package statstest

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// PearsonCorrelation tests for a linear association between two paired
// samples using the Pearson product-moment correlation coefficient.
//
// The null hypothesis is that the population correlation is zero. The
// returned confidence interval for r is two-sided and computed from the
// Fisher z-transformation, regardless of the alternative hypothesis.
func PearsonCorrelation(x, y []float64, alt Alternative, conf float64) (CorrelationResult, error) {
	if err := validateSample(x, 3); err != nil {
		return CorrelationResult{}, err
	}
	if err := validateSample(y, 3); err != nil {
		return CorrelationResult{}, err
	}
	if len(x) != len(y) {
		return CorrelationResult{}, ErrMismatchedLengths
	}
	if err := validateAlternative(alt); err != nil {
		return CorrelationResult{}, err
	}
	if err := validateConfidence(conf); err != nil {
		return CorrelationResult{}, err
	}

	r, err := pearsonR(x, y)
	if err != nil {
		return CorrelationResult{}, err
	}

	return correlationFromR(r, len(x), alt, conf, "Pearson correlation")
}

// SpearmanCorrelation tests for a monotonic association between two paired
// samples using Spearman's rank correlation coefficient.
//
// The p-value and confidence interval use the same t-distribution and
// Fisher z-transformation approximations as PearsonCorrelation, applied to
// the ranked data. This is standard practice but is less accurate for very
// small samples with many tied ranks.
func SpearmanCorrelation(x, y []float64, alt Alternative, conf float64) (CorrelationResult, error) {
	if err := validateSample(x, 3); err != nil {
		return CorrelationResult{}, err
	}
	if err := validateSample(y, 3); err != nil {
		return CorrelationResult{}, err
	}
	if len(x) != len(y) {
		return CorrelationResult{}, ErrMismatchedLengths
	}
	if err := validateAlternative(alt); err != nil {
		return CorrelationResult{}, err
	}
	if err := validateConfidence(conf); err != nil {
		return CorrelationResult{}, err
	}

	rx := rankValues(x)
	ry := rankValues(y)

	r, err := pearsonR(rx, ry)
	if err != nil {
		return CorrelationResult{}, err
	}

	res, err := correlationFromR(r, len(x), alt, conf, "Spearman rank correlation")
	return res, err
}

func pearsonR(x, y []float64) (float64, error) {
	mx := mean(x)
	my := mean(y)

	var sxy, sxx, syy float64
	for i := range x {
		dx := x[i] - mx
		dy := y[i] - my
		sxy += dx * dy
		sxx += dx * dx
		syy += dy * dy
	}

	if sxx == 0 || syy == 0 {
		return 0, ErrZeroVariance
	}

	r := sxy / math.Sqrt(sxx*syy)
	if r > 1 {
		r = 1
	} else if r < -1 {
		r = -1
	}
	return r, nil
}

func correlationFromR(r float64, n int, alt Alternative, conf float64, method string) (CorrelationResult, error) {
	df := float64(n - 2)
	if df <= 0 {
		return CorrelationResult{}, ErrInsufficientDF
	}

	var t float64
	if math.Abs(r) >= 1 {
		t = math.Inf(int(math.Copysign(1, r)))
	} else {
		t = r * math.Sqrt(df/(1-r*r))
	}
	p := tPValue(t, df, alt)

	var ci ConfidenceInterval
	if math.Abs(r) >= 1 || n < 4 {
		ci = ConfidenceInterval{Level: conf, Low: r, High: r}
	} else {
		z := math.Atanh(r)
		se := 1 / math.Sqrt(float64(n)-3)
		zcrit := distuv.Normal{Mu: 0, Sigma: 1}.Quantile(1 - (1-conf)/2)
		ci = ConfidenceInterval{
			Level: conf,
			Low:   math.Tanh(z - zcrit*se),
			High:  math.Tanh(z + zcrit*se),
		}
	}

	return CorrelationResult{
		R:           r,
		PValue:      p,
		DF:          df,
		N:           n,
		CI:          ci,
		Method:      method,
		Alternative: alt,
	}, nil
}
