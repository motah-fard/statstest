package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestLevenesTestEqualVariance(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{2, 3, 4, 5, 6}
	g3 := []float64{3, 4, 5, 6, 7}

	got, err := LevenesTest(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 0, 1e-9)
	assertFloatClose(t, got.PValue, 1, 1e-9)
	if got.Method != "Levene's test (Brown-Forsythe)" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestLevenesTestUnequalVariance(t *testing.T) {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2, 8.5}
	g2 := []float64{7.5, 9.0, 6.8, 9.5, 7.0, 8.8}
	g3 := []float64{8.0, 8.1, 7.9, 8.0, 8.1, 7.95}

	got, err := LevenesTest(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.Statistic > 0) {
		t.Fatalf("expected positive statistic, got %v", got.Statistic)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
	if got.DFBetween != 2 {
		t.Fatalf("got DFBetween %d want 2", got.DFBetween)
	}
	if got.DFWithin != 15 {
		t.Fatalf("got DFWithin %d want 15", got.DFWithin)
	}
}

func TestLevenesTestRejectsTooFewGroups(t *testing.T) {
	_, err := LevenesTest([]float64{1, 2, 3})
	if !errors.Is(err, ErrTooFewGroups) {
		t.Fatalf("expected ErrTooFewGroups, got %v", err)
	}
}

func TestLevenesTestRejectsNaN(t *testing.T) {
	_, err := LevenesTest([]float64{1, 2, 3}, []float64{4, math.NaN(), 6})
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}

func TestMedianOddEven(t *testing.T) {
	assertFloatClose(t, median([]float64{3, 1, 2}), 2, 1e-12)
	assertFloatClose(t, median([]float64{4, 1, 3, 2}), 2.5, 1e-12)
}
