package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestShapiroWilkNormalLikeSample(t *testing.T) {
	x := []float64{5.1, 4.9, 5.3, 5.0, 4.8, 5.2, 5.4, 4.7, 5.1, 5.0}

	got, err := ShapiroWilk(x)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.W > 0 && got.W <= 1) {
		t.Fatalf("expected W in (0, 1], got %v", got.W)
	}
	if !(got.PValue > 0.5) {
		t.Fatalf("expected a large p-value for a normal-like sample, got %v", got.PValue)
	}
	if got.Method != "Shapiro-Wilk test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestShapiroWilkSkewedSampleSmallPValue(t *testing.T) {
	x := []float64{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 5, 6, 7, 8, 10, 12, 20}

	got, err := ShapiroWilk(x)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.PValue < 0.05) {
		t.Fatalf("expected a small p-value for a skewed sample, got %v", got.PValue)
	}
}

func TestShapiroWilkRejectsTooSmallSample(t *testing.T) {
	_, err := ShapiroWilk([]float64{1, 2})
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestShapiroWilkRejectsConstantSample(t *testing.T) {
	_, err := ShapiroWilk([]float64{5, 5, 5, 5})
	if !errors.Is(err, ErrZeroVariance) {
		t.Fatalf("expected ErrZeroVariance, got %v", err)
	}
}

func TestShapiroWilkRejectsNaN(t *testing.T) {
	_, err := ShapiroWilk([]float64{1, 2, math.NaN()})
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}
