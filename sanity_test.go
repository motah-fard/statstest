package statstest

import (
	"testing"
)

func TestReferenceTTestTwoSamplePooled(t *testing.T) {
	x := []float64{2, 4, 6, 8, 10}
	y := []float64{1, 3, 5, 7, 9}

	got, err := TTestTwoSample(x, y, TwoSampleTOptions{
		EqualVariance:   true,
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 0.5, 1e-12)
	assertFloatClose(t, got.DF, 8, 1e-12)
	assertFloatClose(t, got.MeanDiff, 1, 1e-12)
}
func TestReferenceChiSquareIndependence(t *testing.T) {
	observed := [][]int{
		{10, 20},
		{20, 40},
	}

	got, err := ChiSquareIndependence(observed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.Statistic, 0, 1e-12)
	assertFloatClose(t, got.PValue, 1, 1e-12)
	assertFloatClose(t, float64(got.DF), 1, 1e-12)
}
func TestReferenceFishersExactSymmetric(t *testing.T) {
	table := [2][2]int{
		{5, 5},
		{5, 5},
	}

	got, err := FishersExact2x2(table, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFloatClose(t, got.OddsRatio, 1, 1e-12)
	if !(got.PValue > 0 && got.PValue <= 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}
