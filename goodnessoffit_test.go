package statstest

import (
	"errors"
	"testing"
)

func TestChiSquareGoodnessOfFitUniform(t *testing.T) {
	got, err := ChiSquareGoodnessOfFit([]float64{25, 25, 25, 25}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 0, 1e-12)
	assertFloatClose(t, got.PValue, 1, 1e-12)
	if got.DF != 3 {
		t.Fatalf("got df %d want 3", got.DF)
	}
	if got.Method != "Chi-square goodness-of-fit test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestChiSquareGoodnessOfFitCustomExpected(t *testing.T) {
	got, err := ChiSquareGoodnessOfFit([]float64{50, 30, 20}, []float64{40, 40, 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 5.0, 1e-12)
	if got.DF != 2 {
		t.Fatalf("got df %d want 2", got.DF)
	}
}

func TestChiSquareGoodnessOfFitRejectsMismatchedLengths(t *testing.T) {
	_, err := ChiSquareGoodnessOfFit([]float64{10, 20, 30}, []float64{10, 20})
	if !errors.Is(err, ErrMismatchedLengths) {
		t.Fatalf("expected ErrMismatchedLengths, got %v", err)
	}
}

func TestChiSquareGoodnessOfFitRejectsNonPositiveExpected(t *testing.T) {
	_, err := ChiSquareGoodnessOfFit([]float64{10, 20, 30}, []float64{10, 0, 30})
	if !errors.Is(err, ErrInvalidExpectedCounts) {
		t.Fatalf("expected ErrInvalidExpectedCounts, got %v", err)
	}
}

func TestChiSquareGoodnessOfFitRejectsNegativeObserved(t *testing.T) {
	_, err := ChiSquareGoodnessOfFit([]float64{10, -5, 30}, nil)
	if !errors.Is(err, ErrNegativeCount) {
		t.Fatalf("expected ErrNegativeCount, got %v", err)
	}
}

func TestChiSquareGoodnessOfFitRejectsTooFewCategories(t *testing.T) {
	_, err := ChiSquareGoodnessOfFit([]float64{10}, nil)
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}
