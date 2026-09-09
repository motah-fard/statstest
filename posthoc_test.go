package statstest

import (
	"errors"
	"testing"
)

func TestTukeyHSDPairwiseComparisons(t *testing.T) {
	g1 := []float64{1, 2, 3}
	g2 := []float64{4, 5, 6}
	g3 := []float64{7, 8, 9}

	got, err := TukeyHSD(0.95, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Comparisons) != 3 {
		t.Fatalf("got %d comparisons, want 3", len(got.Comparisons))
	}
	if got.Method != "Tukey HSD" {
		t.Fatalf("unexpected method: %s", got.Method)
	}

	for _, c := range got.Comparisons {
		if !(c.PValue >= 0 && c.PValue <= 1) {
			t.Fatalf("p-value out of range: %v", c.PValue)
		}
		if !(c.CI.Low < c.MeanDiff && c.CI.High > c.MeanDiff) {
			t.Fatalf("expected CI to contain mean diff for (%d,%d)", c.GroupI, c.GroupJ)
		}
	}
}

func TestTukeyHSDSameMeansLargePValues(t *testing.T) {
	g1 := []float64{1, 2, 3, 4}
	g2 := []float64{1, 2, 3, 4}
	g3 := []float64{1, 2, 3, 4}

	got, err := TukeyHSD(0.95, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, c := range got.Comparisons {
		assertFloatClose(t, c.PValue, 1, 1e-6)
	}
}

func TestTukeyHSDRejectsTooFewGroups(t *testing.T) {
	_, err := TukeyHSD(0.95, []float64{1, 2, 3})
	if !errors.Is(err, ErrTooFewGroups) {
		t.Fatalf("expected ErrTooFewGroups, got %v", err)
	}
}

func TestTukeyHSDRejectsInvalidConfidence(t *testing.T) {
	_, err := TukeyHSD(1.5, []float64{1, 2, 3}, []float64{4, 5, 6})
	if !errors.Is(err, ErrInvalidConfidenceLevel) {
		t.Fatalf("expected ErrInvalidConfidenceLevel, got %v", err)
	}
}

func TestDunnTestPairwiseComparisons(t *testing.T) {
	g1 := []float64{1, 2, 3, 4}
	g2 := []float64{5, 6, 7, 8}
	g3 := []float64{9, 10, 11, 12}

	got, err := DunnTest(Bonferroni, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.Comparisons) != 3 {
		t.Fatalf("got %d comparisons, want 3", len(got.Comparisons))
	}
	if got.Method != "Dunn's test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if got.AdjustMethod != Bonferroni {
		t.Fatalf("unexpected adjust method: %s", got.AdjustMethod)
	}

	for _, c := range got.Comparisons {
		if !(c.AdjustedPValue >= c.PValue) {
			t.Fatalf("expected adjusted p-value >= raw p-value, got %v < %v", c.AdjustedPValue, c.PValue)
		}
	}
}

func TestDunnTestRejectsTooFewGroups(t *testing.T) {
	_, err := DunnTest(Bonferroni, []float64{1, 2, 3})
	if !errors.Is(err, ErrTooFewGroups) {
		t.Fatalf("expected ErrTooFewGroups, got %v", err)
	}
}

func TestDunnTestRejectsInvalidAdjustMethod(t *testing.T) {
	_, err := DunnTest(PAdjustMethod("bogus"), []float64{1, 2, 3}, []float64{4, 5, 6})
	if !errors.Is(err, ErrInvalidPAdjustMethod) {
		t.Fatalf("expected ErrInvalidPAdjustMethod, got %v", err)
	}
}
