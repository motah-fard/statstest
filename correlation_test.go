package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestPearsonCorrelationPerfectPositive(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10}

	got, err := PearsonCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.R, 1.0, 1e-12)
	assertFloatClose(t, got.CI.Low, 1.0, 1e-9)
	assertFloatClose(t, got.CI.High, 1.0, 1e-9)
	if got.DF != 3 {
		t.Fatalf("got df %.0f want 3", got.DF)
	}
	if got.Method != "Pearson correlation" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestPearsonCorrelationNoAssociation(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6}
	y := []float64{3, 1, 4, 1, 5, 9}

	got, err := PearsonCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
	if !(got.CI.Low < got.R && got.CI.High > got.R) {
		t.Fatal("expected CI to contain r")
	}
}

func TestPearsonCorrelationRejectsShortSample(t *testing.T) {
	_, err := PearsonCorrelation([]float64{1, 2}, []float64{1, 2}, TwoSided, 0.95)
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestPearsonCorrelationRejectsMismatchedLengths(t *testing.T) {
	_, err := PearsonCorrelation([]float64{1, 2, 3, 4}, []float64{1, 2, 3}, TwoSided, 0.95)
	if !errors.Is(err, ErrMismatchedLengths) {
		t.Fatalf("expected ErrMismatchedLengths, got %v", err)
	}
}

func TestPearsonCorrelationRejectsZeroVariance(t *testing.T) {
	_, err := PearsonCorrelation([]float64{1, 1, 1}, []float64{1, 2, 3}, TwoSided, 0.95)
	if !errors.Is(err, ErrZeroVariance) {
		t.Fatalf("expected ErrZeroVariance, got %v", err)
	}
}

func TestPearsonCorrelationRejectsInvalidAlternative(t *testing.T) {
	_, err := PearsonCorrelation([]float64{1, 2, 3}, []float64{1, 2, 3}, Alternative("weird"), 0.95)
	if !errors.Is(err, ErrInvalidAlternative) {
		t.Fatalf("expected ErrInvalidAlternative, got %v", err)
	}
}

func TestPearsonCorrelationRejectsInvalidConfidence(t *testing.T) {
	_, err := PearsonCorrelation([]float64{1, 2, 3}, []float64{1, 2, 3}, TwoSided, 1.5)
	if !errors.Is(err, ErrInvalidConfidenceLevel) {
		t.Fatalf("expected ErrInvalidConfidenceLevel, got %v", err)
	}
}

func TestSpearmanCorrelationMonotonic(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 20, 15, 40, 50}

	got, err := SpearmanCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Method != "Spearman rank correlation" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.R > 0.5) {
		t.Fatalf("expected strong positive rank correlation, got %v", got.R)
	}
}

func TestSpearmanCorrelationRejectsShortSample(t *testing.T) {
	_, err := SpearmanCorrelation([]float64{1, 2}, []float64{1, 2}, TwoSided, 0.95)
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestCorrelationDirectionalAlternatives(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{2, 3, 5, 4, 6, 7, 9, 8}

	two, err := PearsonCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	greater, err := PearsonCorrelation(x, y, Greater, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	less, err := PearsonCorrelation(x, y, Less, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if math.Abs(two.PValue-2*greater.PValue) > 1e-9 {
		t.Fatalf("expected two-sided p to be twice greater-tail p, got %v vs %v", two.PValue, greater.PValue)
	}
	if !(less.PValue > 0.5) {
		t.Fatalf("expected less-tail p > 0.5 for positive r, got %v", less.PValue)
	}
}
