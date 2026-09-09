package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestKruskalWallisSameDistribution(t *testing.T) {
	g1 := []float64{1, 2, 3, 4}
	g2 := []float64{1, 2, 3, 4}
	g3 := []float64{1, 2, 3, 4}

	got, err := KruskalWallis(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.H, 0, 1e-9)
	assertFloatClose(t, got.PValue, 1, 1e-9)
	if got.DF != 2 {
		t.Fatalf("got df %d want 2", got.DF)
	}
	if got.Method != "Kruskal-Wallis test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestKruskalWallisLargeDifference(t *testing.T) {
	g1 := []float64{1, 2, 3, 4}
	g2 := []float64{11, 12, 13, 14}
	g3 := []float64{21, 22, 23, 24}

	got, err := KruskalWallis(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.H > 0) {
		t.Fatalf("expected positive H, got %v", got.H)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
}

func TestKruskalWallisRejectsTooFewGroups(t *testing.T) {
	_, err := KruskalWallis([]float64{1, 2, 3})
	if !errors.Is(err, ErrTooFewGroups) {
		t.Fatalf("expected ErrTooFewGroups, got %v", err)
	}
}

func TestKruskalWallisRejectsEmptyGroup(t *testing.T) {
	_, err := KruskalWallis([]float64{1, 2, 3}, []float64{})
	if !errors.Is(err, ErrEmptySample) {
		t.Fatalf("expected ErrEmptySample, got %v", err)
	}
}

func TestKruskalWallisRejectsNaN(t *testing.T) {
	_, err := KruskalWallis([]float64{1, 2, 3}, []float64{4, math.NaN(), 6})
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}
