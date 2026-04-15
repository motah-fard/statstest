package statstest

import (
	"testing"
)

func TestAdjustPValuesBonferroni(t *testing.T) {
	p := []float64{0.01, 0.03, 0.2, 0.001}
	got, err := AdjustPValues(p, Bonferroni)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []float64{0.04, 0.12, 0.8, 0.004}
	assertFloatSlicesClose(t, got, want, 1e-12)
}

func TestAdjustPValuesHolm(t *testing.T) {
	p := []float64{0.01, 0.03, 0.2, 0.001}
	got, err := AdjustPValues(p, Holm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []float64{0.03, 0.06, 0.2, 0.004}
	assertFloatSlicesClose(t, got, want, 1e-12)
}

func TestAdjustPValuesBH(t *testing.T) {
	p := []float64{0.01, 0.03, 0.2, 0.001}
	got, err := AdjustPValues(p, BenjaminiHochberg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []float64{0.02, 0.04, 0.2, 0.004}
	assertFloatSlicesClose(t, got, want, 1e-12)
}

func TestAdjustPValuesInvalid(t *testing.T) {
	_, err := AdjustPValues([]float64{0.2, 1.2}, Bonferroni)
	if err == nil {
		t.Fatal("expected error for invalid p-value")
	}
}
