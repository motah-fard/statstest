package statstest

import (
	"context"
	"math/rand"
	"sort"
	"time"
)

// minResamples is the smallest NumResamples BootstrapCI will accept.
// Below this, a percentile confidence interval is too noisy to be
// meaningful; 1000+ is recommended in practice.
const minResamples = 100

// BootstrapCI computes a percentile bootstrap confidence interval for an
// arbitrary statistic, by resampling x with replacement opts.NumResamples
// times and taking the (alpha/2, 1-alpha/2) percentiles of the resampled
// statistic's distribution.
//
// Unlike every other confidence interval in this package, a bootstrap CI
// has no closed form to check against R or SciPy bit-for-bit — its
// correctness is instead validated in this package's tests by checking
// that it converges to the classical CI for the sample mean, and that the
// same seed always reproduces the same result.
//
// statistic must not mutate its argument. It is called concurrently
// across resamples once there is enough work to be worth it, so it must
// be safe to call from multiple goroutines at once — this holds
// automatically for any pure function of its input, which covers every
// statistic in this package (mean, median, a t-statistic, Cohen's d,
// and so on).
//
// This is equivalent to calling BootstrapCIContext with
// context.Background().
func BootstrapCI(x []float64, statistic func([]float64) float64, opts BootstrapOptions) (BootstrapResult, error) {
	return BootstrapCIContext(context.Background(), x, statistic, opts)
}

// BootstrapCIContext is BootstrapCI with cancellation support. Because
// statistic is arbitrary, unbounded-cost code supplied by the caller,
// this is the function in the package most likely to need a deadline. If
// ctx is done before every resample finishes, it returns a zero
// BootstrapResult and ctx.Err().
func BootstrapCIContext(ctx context.Context, x []float64, statistic func([]float64) float64, opts BootstrapOptions) (BootstrapResult, error) {
	if err := ctx.Err(); err != nil {
		return BootstrapResult{}, err
	}
	if err := validateSample(x, 2); err != nil {
		return BootstrapResult{}, err
	}
	if statistic == nil {
		return BootstrapResult{}, ErrNilStatistic
	}
	if opts.NumResamples < minResamples {
		return BootstrapResult{}, ErrTooFewResamples
	}
	if err := validateConfidence(opts.ConfidenceLevel); err != nil {
		return BootstrapResult{}, err
	}

	estimate := statistic(x)

	baseSeed := time.Now().UnixNano()
	if opts.Seed != nil {
		baseSeed = *opts.Seed
	}

	n := len(x)
	replicates := make([]float64, opts.NumResamples)
	err := parallelComputeEachContext(ctx, opts.NumResamples, func(i int) {
		rng := rand.New(rand.NewSource(mixSeed(baseSeed, i)))
		resample := make([]float64, n)
		for j := range resample {
			resample[j] = x[rng.Intn(n)]
		}
		replicates[i] = statistic(resample)
	})
	if err != nil {
		return BootstrapResult{}, err
	}

	sort.Float64s(replicates)

	alpha := 1 - opts.ConfidenceLevel
	ci := ConfidenceInterval{
		Level: opts.ConfidenceLevel,
		Low:   percentile(replicates, alpha/2),
		High:  percentile(replicates, 1-alpha/2),
	}

	return BootstrapResult{
		Estimate:     estimate,
		StdError:     sampleStdDev(replicates),
		CI:           ci,
		NumResamples: opts.NumResamples,
		Method:       "Percentile bootstrap",
	}, nil
}

// percentile returns the p-quantile (p in [0, 1]) of an already-sorted
// slice, using linear interpolation between the two nearest ranks — the
// same convention numpy's default "linear" method uses.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 1 {
		return sorted[0]
	}

	pos := p * float64(len(sorted)-1)
	lo := int(pos)
	if lo >= len(sorted)-1 {
		return sorted[len(sorted)-1]
	}
	frac := pos - float64(lo)
	return sorted[lo] + frac*(sorted[lo+1]-sorted[lo])
}

// mixSeed derives an independent seed for resample i from a base seed,
// using the splitmix64 finalizer so that seeds for adjacent i are not
// correlated the way naively adding i to base could be with some PRNGs.
func mixSeed(base int64, i int) int64 {
	h := uint64(base) + uint64(i)*0x9E3779B97F4A7C15
	h = (h ^ (h >> 30)) * 0xBF58476D1CE4E5B9
	h = (h ^ (h >> 27)) * 0x94D049BB133111EB
	h ^= h >> 31
	return int64(h)
}
