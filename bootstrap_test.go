package statstest

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"testing"
)

func TestPercentileKnownValues(t *testing.T) {
	sorted := []float64{1, 2, 3, 4, 5}
	assertFloatClose(t, percentile(sorted, 0), 1, 1e-12)
	assertFloatClose(t, percentile(sorted, 1), 5, 1e-12)
	assertFloatClose(t, percentile(sorted, 0.5), 3, 1e-12)
	assertFloatClose(t, percentile(sorted, 0.25), 2, 1e-12)
	assertFloatClose(t, percentile([]float64{7}, 0.5), 7, 1e-12)
}

func TestMixSeedIsDeterministicAndVaries(t *testing.T) {
	a := mixSeed(42, 0)
	b := mixSeed(42, 0)
	if a != b {
		t.Fatalf("expected mixSeed to be deterministic, got %d and %d", a, b)
	}
	c := mixSeed(42, 1)
	if a == c {
		t.Fatalf("expected different indices to produce different seeds")
	}
}

func meanStatistic(x []float64) float64 {
	return mean(x)
}

func TestBootstrapCIDeterministicWithSeed(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	x := make([]float64, 40)
	for i := range x {
		x[i] = rng.NormFloat64()*2 + 10
	}

	seed := int64(123)
	opts := BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95, Seed: &seed}

	a, err := BootstrapCI(x, meanStatistic, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := BootstrapCI(x, meanStatistic, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a != b {
		t.Fatalf("expected identical results for the same seed, got %+v vs %+v", a, b)
	}
}

func TestBootstrapCIDifferentSeedsDiffer(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	x := make([]float64, 40)
	for i := range x {
		x[i] = rng.NormFloat64()*2 + 10
	}

	seed1, seed2 := int64(1), int64(2)
	a, err := BootstrapCI(x, meanStatistic, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95, Seed: &seed1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := BootstrapCI(x, meanStatistic, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95, Seed: &seed2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if a.CI == b.CI {
		t.Fatalf("expected different seeds to (almost certainly) produce different CIs")
	}
}

// TestBootstrapCIConvergesToClassicalCI is the closest thing to a
// reference-value check available for a randomized method: with a large
// sample and many resamples, the bootstrap CI for the mean should land
// close to the classical t-based CI from TTestOneSample.
func TestBootstrapCIConvergesToClassicalCI(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	x := make([]float64, 500)
	for i := range x {
		x[i] = rng.NormFloat64()*3 + 50
	}

	classical, err := TTestOneSample(x, 0, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	seed := int64(99)
	boot, err := BootstrapCI(x, meanStatistic, BootstrapOptions{
		NumResamples:    5000,
		ConfidenceLevel: 0.95,
		Seed:            &seed,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, boot.Estimate, classical.Mean, 1e-9)

	// The two CIs should be close, not identical — allow a generous
	// tolerance relative to the interval's own width rather than an
	// absolute one.
	tol := (classical.CI.High - classical.CI.Low) * 0.15
	if math.Abs(boot.CI.Low-classical.CI.Low) > tol {
		t.Fatalf("bootstrap CI low %v too far from classical CI low %v (tol %v)", boot.CI.Low, classical.CI.Low, tol)
	}
	if math.Abs(boot.CI.High-classical.CI.High) > tol {
		t.Fatalf("bootstrap CI high %v too far from classical CI high %v (tol %v)", boot.CI.High, classical.CI.High, tol)
	}
}

func TestBootstrapCIRejectsTooFewResamples(t *testing.T) {
	_, err := BootstrapCI([]float64{1, 2, 3}, meanStatistic, BootstrapOptions{NumResamples: 10, ConfidenceLevel: 0.95})
	if !errors.Is(err, ErrTooFewResamples) {
		t.Fatalf("expected ErrTooFewResamples, got %v", err)
	}
}

func TestBootstrapCIRejectsNilStatistic(t *testing.T) {
	_, err := BootstrapCI([]float64{1, 2, 3}, nil, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95})
	if !errors.Is(err, ErrNilStatistic) {
		t.Fatalf("expected ErrNilStatistic, got %v", err)
	}
}

func TestBootstrapCIRejectsTooSmallSample(t *testing.T) {
	_, err := BootstrapCI([]float64{1}, meanStatistic, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95})
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestBootstrapCIRejectsInvalidConfidence(t *testing.T) {
	_, err := BootstrapCI([]float64{1, 2, 3}, meanStatistic, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 1.5})
	if !errors.Is(err, ErrInvalidConfidenceLevel) {
		t.Fatalf("expected ErrInvalidConfidenceLevel, got %v", err)
	}
}

func TestBootstrapCIContextCancellationBeforeCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := BootstrapCIContext(ctx, []float64{1, 2, 3}, meanStatistic, BootstrapOptions{NumResamples: 500, ConfidenceLevel: 0.95})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestBootstrapCIWorksWithNonMeanStatistic(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	x := make([]float64, 60)
	for i := range x {
		x[i] = rng.NormFloat64()
	}

	medianStatistic := func(s []float64) float64 { return median(s) }

	res, err := BootstrapCI(x, medianStatistic, BootstrapOptions{NumResamples: 1000, ConfidenceLevel: 0.9})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(res.CI.Low < res.Estimate && res.CI.High > res.Estimate) {
		t.Fatalf("expected CI to contain the point estimate: CI=[%v,%v] estimate=%v", res.CI.Low, res.CI.High, res.Estimate)
	}
	if res.Method != "Percentile bootstrap" {
		t.Fatalf("unexpected method: %s", res.Method)
	}
}
