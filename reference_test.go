// -----------------------------------------------------------------------------
// Reference-value tests
//
// These tests verify results against external statistical references or
// package-documented conventions, rather than only checking that functions run.
//
// Convention notes:
// - Mann-Whitney expected values follow this package's current asymptotic
//   p-value convention.
// - WilcoxonSignedRank reports W as the sum of positive ranks (W+).
// - OneWayANOVA reports FStatistic, DFBetween, and DFWithin.
// -----------------------------------------------------------------------------

package statstest

import (
	"math"
	"testing"
)

func TestReferencePAdjustBH(t *testing.T) {
	p := []float64{0.001, 0.01, 0.03, 0.2}

	got, err := AdjustPValues(p, BenjaminiHochberg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []float64{0.004, 0.02, 0.04, 0.2}
	requireSliceAlmostEqual(t, "BH adjusted p-values", got, want, 1e-12)
}

func almostEqual(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol
}

func requireAlmostEqual(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if !almostEqual(got, want, tol) {
		t.Fatalf("%s: got %.12g, want %.12g (tol=%.12g)", name, got, want, tol)
	}
}

func requireSliceAlmostEqual(t *testing.T, name string, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: length mismatch: got %d, want %d", name, len(got), len(want))
	}
	for i := range got {
		if !almostEqual(got[i], want[i], tol) {
			t.Fatalf("%s[%d]: got %.12g, want %.12g (tol=%.12g)", name, i, got[i], want[i], tol)
		}
	}
}

func TestReferenceTTestOneSample(t *testing.T) {
	x := []float64{2.3, 2.5, 2.1, 2.7, 2.4}

	wantStat := 4.000000000000003
	wantP := 0.016130089900092518
	wantDF := 4.0
	tol := 1e-12

	res, err := TTestOneSample(x, 2.0, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceTTestWelch(t *testing.T) {
	x := []float64{12.1, 11.8, 12.5, 12.0, 11.9}
	y := []float64{10.2, 10.4, 10.1, 10.3}

	wantStat := 13.21250013653373
	wantP := 0.000012203818575066068
	wantDF := 5.96150036361814
	tol := 1e-10

	res, err := TTestTwoSample(x, y, TwoSampleTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
		EqualVariance:   false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceTTestPaired(t *testing.T) {
	x := []float64{10.2, 9.8, 10.5, 10.1, 9.9, 10.3}
	y := []float64{10.0, 9.7, 10.1, 9.8, 9.5, 10.0}

	wantStat := 5.936657514041429
	wantP := 0.0019358364032480116
	wantDF := 5.0
	tol := 1e-12

	res, err := TTestPaired(x, y, PairedTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceOneWayANOVA(t *testing.T) {
	groups := [][]float64{
		{8.1, 8.3, 7.9, 8.0, 8.2},
		{8.8, 9.0, 8.7, 8.9, 9.1},
		{7.5, 7.6, 7.4, 7.7, 7.5},
	}

	wantStat := 111.23809523809491
	wantP := 0.00000001796780307526363
	wantDFBetween := 2.0
	wantDFWithin := 12.0
	tol := 1e-9

	res, err := OneWayANOVA(groups...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "F statistic", res.FStatistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df_between", float64(res.DFBetween), wantDFBetween, tol)
	requireAlmostEqual(t, "df_within", float64(res.DFWithin), wantDFWithin, tol)
}

func TestReferenceMannWhitney(t *testing.T) {
	x := []float64{14, 15, 16, 17, 18, 19, 20, 21}
	y := []float64{8, 9, 10, 11, 12, 13, 14, 15}

	wantStat := 62.0
	wantP := 0.00160347596869
	tol := 1e-9

	res, err := MannWhitneyU(x, y, MannWhitneyOptions{
		Alternative: TwoSided,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "U statistic", res.U, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceWilcoxonSignedRank(t *testing.T) {
	x := []float64{120, 118, 121, 119, 117, 122, 116, 123, 118, 120}
	y := []float64{115, 117, 119, 118, 116, 120, 114, 121, 117, 119}

	wantStat := 55.0
	wantP := 0.004245585155691297
	tol := 1e-9

	res, err := WilcoxonSignedRank(x, y, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "W statistic", res.W, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceChiSquare(t *testing.T) {
	table := [][]int{
		{20, 30, 25},
		{22, 28, 35},
	}

	wantStat := 1.2105991822016162
	wantP := 0.5459108521050702
	wantDF := 2.0
	tol := 1e-12

	res, err := ChiSquareIndependence(table) // false = no Yates correction
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "chi-square statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", float64(res.DF), wantDF, tol)
}

func TestReferenceFisherExact(t *testing.T) {
	table := [2][2]int{
		{8, 2},
		{1, 5},
	}

	wantOddsRatio := 20.0
	wantP := 0.034965034965034975
	tol := 1e-12

	res, err := FishersExact2x2(table, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "odds ratio", res.OddsRatio, wantOddsRatio, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "odds ratio", res.OddsRatio, wantOddsRatio, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferencePAdjustHolm(t *testing.T) {
	p := []float64{0.001, 0.01, 0.03, 0.2, 0.5}

	want := []float64{0.005, 0.04, 0.09, 0.4, 0.5}
	tol := 1e-12

	got, err := AdjustPValues(p, Holm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireSliceAlmostEqual(t, "adjusted p-values", got, want, tol)
}
