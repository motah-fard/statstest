package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestOneWayANOVAExactValues(t *testing.T) {
	g1 := []float64{1, 2, 3}
	g2 := []float64{4, 5, 6}
	g3 := []float64{7, 8, 9}

	got, err := OneWayANOVA(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// grand mean = 5
	// group means = 2, 5, 8
	// ssBetween = 3*(2-5)^2 + 3*(5-5)^2 + 3*(8-5)^2 = 54
	// ssWithin = 2 + 2 + 2 = 6
	// ssTotal = 60
	// dfBetween = 2
	// dfWithin = 6
	// msBetween = 27
	// msWithin = 1
	// F = 27
	wantSSBetween := 54.0
	wantSSWithin := 6.0
	wantSSTotal := 60.0
	wantF := 27.0
	wantEta := 54.0 / 60.0

	assertFloatClose(t, got.SSBetween, wantSSBetween, 1e-12)
	assertFloatClose(t, got.SSWithin, wantSSWithin, 1e-12)
	assertFloatClose(t, got.SSTotal, wantSSTotal, 1e-12)
	assertFloatClose(t, got.FStatistic, wantF, 1e-12)
	assertFloatClose(t, got.EtaSquared, wantEta, 1e-12)

	if got.DFBetween != 2 {
		t.Fatalf("got DFBetween %d want 2", got.DFBetween)
	}
	if got.DFWithin != 6 {
		t.Fatalf("got DFWithin %d want 6", got.DFWithin)
	}
	if got.Method != "One-way ANOVA" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestOneWayANOVASameMeans(t *testing.T) {
	g1 := []float64{1, 2, 3}
	g2 := []float64{1, 2, 3}
	g3 := []float64{1, 2, 3}

	got, err := OneWayANOVA(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.SSBetween, 0, 1e-12)
	assertFloatClose(t, got.FStatistic, 0, 1e-12)
	assertFloatClose(t, got.EtaSquared, 0, 1e-12)
	assertFloatClose(t, got.PValue, 1, 1e-12)
}

func TestOneWayANOVARejectsTooFewGroups(t *testing.T) {
	_, err := OneWayANOVA([]float64{1, 2, 3})
	if err == nil {
		t.Fatal("expected too few groups error")
	}
}

func TestOneWayANOVARejectsSmallGroup(t *testing.T) {
	_, err := OneWayANOVA(
		[]float64{1, 2, 3},
		[]float64{4},
	)
	if err == nil {
		t.Fatal("expected error for small group")
	}
}

func TestOneWayANOVARejectsZeroVarianceEverywhere(t *testing.T) {
	_, err := OneWayANOVA(
		[]float64{5, 5, 5},
		[]float64{5, 5, 5},
	)
	if err == nil {
		t.Fatal("expected zero variance error")
	}
}

func TestOneWayANOVALargeDifference(t *testing.T) {
	g1 := []float64{10, 11, 9, 10}
	g2 := []float64{20, 19, 21, 20}
	g3 := []float64{30, 31, 29, 30}

	got, err := OneWayANOVA(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.FStatistic > 0) {
		t.Fatalf("expected positive F statistic, got %v", got.FStatistic)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
	if !(got.EtaSquared > 0.8) {
		t.Fatalf("expected large eta squared, got %v", got.EtaSquared)
	}
}

func TestOneWayANOVARejectsNaN(t *testing.T) {
	_, err := OneWayANOVA(
		[]float64{1, 2, 3},
		[]float64{4, math.NaN(), 6},
	)
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}

func TestOneWayANOVARejectsInf(t *testing.T) {
	_, err := OneWayANOVA(
		[]float64{1, 2, 3},
		[]float64{4, math.Inf(1), 6},
	)
	if !errors.Is(err, ErrContainsInf) {
		t.Fatalf("expected ErrContainsInf, got %v", err)
	}
}

func TestOneWayANOVARejectsZeroVarianceGroups(t *testing.T) {
	_, err := OneWayANOVA(
		[]float64{5, 5, 5},
		[]float64{5, 5, 5},
	)
	if !errors.Is(err, ErrZeroVariance) {
		t.Fatalf("expected ErrZeroVariance, got %v", err)
	}
}
