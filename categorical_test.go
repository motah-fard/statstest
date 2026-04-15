package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestChiSquareIndependence2x2Exact(t *testing.T) {
	observed := [][]int{
		{10, 20},
		{20, 40},
	}

	got, err := ChiSquareIndependence(observed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// These rows are proportional, so statistic should be 0.
	assertFloatClose(t, got.Statistic, 0, 1e-12)
	assertFloatClose(t, got.PValue, 1, 1e-12)

	if got.DF != 1 {
		t.Fatalf("got DF %d want 1", got.DF)
	}
	if got.Method != "Chi-square test of independence" {
		t.Fatalf("unexpected method: %s", got.Method)
	}

	wantExpected := [][]float64{
		{10, 20},
		{20, 40},
	}
	for i := range wantExpected {
		assertFloatSlicesClose(t, got.Expected[i], wantExpected[i], 1e-12)
	}
}

func TestChiSquareIndependenceNonZeroStatistic(t *testing.T) {
	observed := [][]int{
		{90, 60, 104, 95},
		{30, 50, 51, 20},
		{30, 40, 45, 35},
	}

	got, err := ChiSquareIndependence(observed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.Statistic > 0) {
		t.Fatalf("expected positive statistic, got %v", got.Statistic)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
	if got.DF != 6 {
		t.Fatalf("got DF %d want 6", got.DF)
	}
	if len(got.Expected) != 3 || len(got.Expected[0]) != 4 {
		t.Fatal("unexpected expected matrix dimensions")
	}
}

func TestChiSquareIndependenceRejectsSmallTable(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{{1, 2}})
	if err == nil {
		t.Fatal("expected invalid table error")
	}
}

func TestChiSquareIndependenceRejectsJaggedTable(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{
		{10, 20, 30},
		{5, 15},
	})
	if !errors.Is(err, ErrInvalidTable) {
		t.Fatalf("expected ErrInvalidTable, got %v", err)
	}
}

func TestChiSquareIndependenceRejectsNegativeCountsExact(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{
		{10, -2},
		{5, 7},
	})
	if !errors.Is(err, ErrNegativeCount) {
		t.Fatalf("expected ErrNegativeCount, got %v", err)
	}
}

func TestChiSquareIndependenceRejectsEmptyTable(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{})
	if !errors.Is(err, ErrInvalidTable) {
		t.Fatalf("expected ErrInvalidTable, got %v", err)
	}
}

func TestChiSquareIndependenceRejectsEmptyRow(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{
		{1, 2},
		{},
	})
	if !errors.Is(err, ErrInvalidTable) {
		t.Fatalf("expected ErrInvalidTable, got %v", err)
	}
}

func TestChiSquareIndependenceRejectsNegativeCounts(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{
		{1, -2},
		{3, 4},
	})
	if err == nil {
		t.Fatal("expected negative count error")
	}
}

func TestChiSquareIndependenceRejectsZeroTotal(t *testing.T) {
	_, err := ChiSquareIndependence([][]int{
		{0, 0},
		{0, 0},
	})
	if err == nil {
		t.Fatal("expected invalid table error")
	}
}

func TestFishersExact2x2SymmetricTable(t *testing.T) {
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
	if got.Method != "Fisher's exact test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
}

func TestFishersExact2x2StrongPositiveAssociation(t *testing.T) {
	table := [2][2]int{
		{8, 1},
		{1, 8},
	}

	got, err := FishersExact2x2(table, Greater)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.OddsRatio > 1) {
		t.Fatalf("expected odds ratio > 1, got %v", got.OddsRatio)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
}

func TestFishersExact2x2StrongNegativeAssociation(t *testing.T) {
	table := [2][2]int{
		{1, 8},
		{8, 1},
	}

	got, err := FishersExact2x2(table, Less)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.OddsRatio < 1) {
		t.Fatalf("expected odds ratio < 1, got %v", got.OddsRatio)
	}
	if !(got.PValue < 0.05) {
		t.Fatalf("expected small p-value, got %v", got.PValue)
	}
}

func TestFishersExact2x2TwoSided(t *testing.T) {
	table := [2][2]int{
		{8, 1},
		{1, 8},
	}

	got, err := FishersExact2x2(table, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !(got.PValue > 0 && got.PValue <= 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestFishersExact2x2InfiniteOddsRatio(t *testing.T) {
	table := [2][2]int{
		{4, 0},
		{1, 3},
	}

	got, err := FishersExact2x2(table, Greater)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !math.IsInf(got.OddsRatio, 1) {
		t.Fatalf("expected +Inf odds ratio, got %v", got.OddsRatio)
	}
}

func TestFishersExact2x2RejectsNegativeCountsExact(t *testing.T) {
	_, err := FishersExact2x2(
		[2][2]int{
			{1, -1},
			{2, 3},
		},
		TwoSided,
	)
	if err == nil {
		t.Fatal("expected negative count error")
	}
}

func TestFishersExact2x2RejectsZeroTotal(t *testing.T) {
	_, err := FishersExact2x2(
		[2][2]int{
			{0, 0},
			{0, 0},
		},
		TwoSided,
	)
	if err == nil {
		t.Fatal("expected invalid table error")
	}
}

func TestFishersExact2x2RejectsInvalidAlternative(t *testing.T) {
	_, err := FishersExact2x2(
		[2][2]int{
			{1, 2},
			{3, 4},
		},
		Alternative("bad"),
	)
	if err == nil {
		t.Fatal("expected invalid alternative error")
	}
}

func TestFishersExact2x2AllowsZeroCells(t *testing.T) {
	_, err := FishersExact2x2(
		[2][2]int{
			{0, 5},
			{3, 7},
		},
		TwoSided,
	)
	if err != nil {
		t.Fatalf("expected zero cells to be allowed, got %v", err)
	}
}
