package statstest

import (
	"errors"
	"testing"
)

func TestProportionOneSampleBasic(t *testing.T) {
	got, err := ProportionOneSample(45, 100, 0.5, ProportionOneSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Proportion, 0.45, 1e-12)
	if got.Method != "One-sample proportion z-test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.CI.Low < got.Proportion && got.CI.High > got.Proportion) {
		t.Fatal("expected CI to contain sample proportion")
	}
}

func TestProportionOneSampleRejectsInvalidCount(t *testing.T) {
	_, err := ProportionOneSample(120, 100, 0.5, ProportionOneSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if !errors.Is(err, ErrInvalidCount) {
		t.Fatalf("expected ErrInvalidCount, got %v", err)
	}
}

func TestProportionOneSampleRejectsInvalidNullProportion(t *testing.T) {
	_, err := ProportionOneSample(5, 10, 1.5, ProportionOneSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if !errors.Is(err, ErrInvalidProportion) {
		t.Fatalf("expected ErrInvalidProportion, got %v", err)
	}
}

func TestProportionOneSampleRejectsZeroTrials(t *testing.T) {
	_, err := ProportionOneSample(0, 0, 0.5, ProportionOneSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestProportionOneSampleContinuityCorrectionShrinksStatistic(t *testing.T) {
	opts := ProportionOneSampleOptions{Alternative: TwoSided, ConfidenceLevel: 0.95}
	withoutCC, err := ProportionOneSample(45, 100, 0.5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	opts.UseContinuityCorrection = true
	withCC, err := ProportionOneSample(45, 100, 0.5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(absFloat(withCC.Statistic) < absFloat(withoutCC.Statistic)) {
		t.Fatalf("expected continuity correction to shrink |z|, got %v vs %v", withCC.Statistic, withoutCC.Statistic)
	}
}

func TestProportionTwoSampleBasic(t *testing.T) {
	got, err := ProportionTwoSample(45, 100, 30, 90, ProportionTwoSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Proportion1, 0.45, 1e-12)
	assertFloatClose(t, got.Proportion2, 30.0/90.0, 1e-12)
	assertFloatClose(t, got.Diff, 0.45-30.0/90.0, 1e-12)
	if got.Method != "Two-sample proportion z-test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestProportionTwoSampleRejectsInvalidCount(t *testing.T) {
	_, err := ProportionTwoSample(-1, 100, 30, 90, ProportionTwoSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if !errors.Is(err, ErrInvalidCount) {
		t.Fatalf("expected ErrInvalidCount, got %v", err)
	}
}

func TestProportionTwoSampleIdenticalProportionsHaveLargePValue(t *testing.T) {
	got, err := ProportionTwoSample(50, 100, 50, 100, ProportionTwoSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 0, 1e-12)
	assertFloatClose(t, got.PValue, 1, 1e-9)
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
