package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestRankWithAverageTies(t *testing.T) {
	x := []float64{1, 2}
	y := []float64{2, 3}

	ranked, tieCounts := rankWithAverageTies(x, y)

	if len(ranked) != 4 {
		t.Fatalf("got %d ranked values, want 4", len(ranked))
	}
	if len(tieCounts) != 1 || tieCounts[0] != 2 {
		t.Fatalf("unexpected tie counts: %#v", tieCounts)
	}

	// Sorted pooled values: 1, 2, 2, 3
	// Ranks should be 1, 2.5, 2.5, 4
	assertFloatClose(t, ranked[0].rank, 1.0, 1e-12)
	assertFloatClose(t, ranked[1].rank, 2.5, 1e-12)
	assertFloatClose(t, ranked[2].rank, 2.5, 1e-12)
	assertFloatClose(t, ranked[3].rank, 4.0, 1e-12)
}

func TestMannWhitneyRejectsInf(t *testing.T) {
	_, err := MannWhitneyU(
		[]float64{1, 2, math.Inf(1)},
		[]float64{3, 4, 5},
		MannWhitneyOptions{Alternative: TwoSided},
	)
	if !errors.Is(err, ErrContainsInf) {
		t.Fatalf("expected ErrContainsInf, got %v", err)
	}
}

func TestMannWhitneyUSimpleSeparation(t *testing.T) {
	x := []float64{1, 2, 3}
	y := []float64{10, 11, 12}

	got, err := MannWhitneyU(x, y, MannWhitneyOptions{
		Alternative:             Less,
		UseContinuityCorrection: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// x values are all smaller, so U for x should be 0
	assertFloatClose(t, got.U, 0, 1e-12)

	if !(got.Z < 0) {
		t.Fatalf("expected negative z, got %v", got.Z)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
	if got.Method != "Mann-Whitney U test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestMannWhitneyUGreaterAlternative(t *testing.T) {
	x := []float64{10, 11, 12}
	y := []float64{1, 2, 3}

	got, err := MannWhitneyU(x, y, MannWhitneyOptions{
		Alternative:             Greater,
		UseContinuityCorrection: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// x values are all larger, so U for x should be max = n1*n2 = 9
	assertFloatClose(t, got.U, 9, 1e-12)

	if !(got.Z > 0) {
		t.Fatalf("expected positive z, got %v", got.Z)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
}

func TestMannWhitneyUWithTies(t *testing.T) {
	x := []float64{1, 2, 2}
	y := []float64{2, 3, 4}

	got, err := MannWhitneyU(x, y, MannWhitneyOptions{
		Alternative:             TwoSided,
		UseContinuityCorrection: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestMannWhitneyURejectsInvalidAlternative(t *testing.T) {
	_, err := MannWhitneyU(
		[]float64{1, 2},
		[]float64{3, 4},
		MannWhitneyOptions{Alternative: Alternative("bad")},
	)
	if !errors.Is(err, ErrInvalidAlternative) {
		t.Fatalf("expected ErrInvalidAlternative, got %v", err)
	}
}

func TestMannWhitneyRejectsNaN(t *testing.T) {
	_, err := MannWhitneyU(
		[]float64{1, 2},
		[]float64{3, math.NaN()},
		MannWhitneyOptions{Alternative: TwoSided},
	)
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}
func TestMannWhitneyRejectsEmpty(t *testing.T) {
	_, err := MannWhitneyU(
		[]float64{},
		[]float64{1, 2},
		MannWhitneyOptions{Alternative: TwoSided},
	)
	if !errors.Is(err, ErrEmptySample) {
		t.Fatalf("expected ErrEmptySample, got %v", err)
	}
}

func TestSignedRanksFromDifferences(t *testing.T) {
	diffs := []float64{2, -1, 0, -2, 1}

	items, tieCounts := signedRanksFromDifferences(diffs)

	// zero should be removed, so 4 items remain
	if len(items) != 4 {
		t.Fatalf("got %d items, want 4", len(items))
	}

	// abs diffs are 1,1,2,2 so tie groups should be [2,2]
	if len(tieCounts) != 2 || tieCounts[0] != 2 || tieCounts[1] != 2 {
		t.Fatalf("unexpected tieCounts: %#v", tieCounts)
	}

	// ranks should be 1.5, 1.5, 3.5, 3.5 in sorted order
	assertFloatClose(t, items[0].rank, 1.5, 1e-12)
	assertFloatClose(t, items[1].rank, 1.5, 1e-12)
	assertFloatClose(t, items[2].rank, 3.5, 1e-12)
	assertFloatClose(t, items[3].rank, 3.5, 1e-12)
}

func TestWilcoxonSignedRankGreater(t *testing.T) {
	x := []float64{10, 12, 14, 16, 18}
	y := []float64{8, 11, 13, 14, 15}

	got, err := WilcoxonSignedRank(x, y, Greater)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.W > 0) {
		t.Fatalf("expected positive W, got %v", got.W)
	}
	if !(got.Z > 0) {
		t.Fatalf("expected positive z, got %v", got.Z)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
	if got.Method != "Wilcoxon signed-rank test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestWilcoxonSignedRankLess(t *testing.T) {
	x := []float64{8, 11, 13, 14, 15}
	y := []float64{10, 12, 14, 16, 18}

	got, err := WilcoxonSignedRank(x, y, Less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.Z < 0) {
		t.Fatalf("expected negative z, got %v", got.Z)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestWilcoxonSignedRankWithZeroDifferences(t *testing.T) {
	x := []float64{10, 12, 14, 16}
	y := []float64{10, 11, 14, 13}

	got, err := WilcoxonSignedRank(x, y, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestWilcoxonRejectsNaN(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{1, 2, math.NaN()},
		[]float64{0, 1, 2},
		TwoSided,
	)
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}

func TestWilcoxonRejectsInf(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{1, 2, math.Inf(1)},
		[]float64{0, 1, 2},
		TwoSided,
	)
	if !errors.Is(err, ErrContainsInf) {
		t.Fatalf("expected ErrContainsInf, got %v", err)
	}
}

func TestWilcoxonRejectsEmpty(t *testing.T) {
	_, err := WilcoxonSignedRank([]float64{}, []float64{}, TwoSided)
	if !errors.Is(err, ErrEmptySample) {
		t.Fatalf("expected ErrEmptySample, got %v", err)
	}
}

func TestWilcoxonRejectsMismatchedLengths(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{1, 2, 3},
		[]float64{1, 2},
		TwoSided,
	)
	if !errors.Is(err, ErrMismatchedLengths) {
		t.Fatalf("expected ErrMismatchedLengths, got %v", err)
	}
}

func TestWilcoxonRejectsAllZeroDifferences(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{1, 2, 3},
		[]float64{1, 2, 3},
		TwoSided,
	)
	if !errors.Is(err, ErrNoNonZeroDifferences) {
		t.Fatalf("expected ErrNoNonZeroDifferences, got %v", err)
	}
}

func TestWilcoxonHandlesTiesAndZeros(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{5, 7, 7, 10, 10},
		[]float64{5, 6, 8, 10, 8},
		TwoSided,
	)
	if err != nil {
		t.Fatalf("expected ties/zeros to be handled, got %v", err)
	}
}
func TestWilcoxonRejectsInvalidAlternative(t *testing.T) {
	_, err := WilcoxonSignedRank(
		[]float64{1, 2, 3},
		[]float64{0, 1, 2},
		Alternative("bad"),
	)
	if !errors.Is(err, ErrInvalidAlternative) {
		t.Fatalf("expected ErrInvalidAlternative, got %v", err)
	}
}
