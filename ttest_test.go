package statstest

import (
	"math"
	"testing"
)

func TestTTestOneSample(t *testing.T) {
	x := []float64{2, 4, 6, 8, 10}

	got, err := TTestOneSample(x, 5, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// mean = 6, sd = sqrt(10), se = sqrt(2), t = 1/sqrt(2)
	wantT := 1.0 / math.Sqrt(2.0)
	assertFloatClose(t, got.Statistic, wantT, 1e-12)

	if got.DF != 4 {
		t.Fatalf("got df %.0f want 4", got.DF)
	}
	if got.Method != "One-sample t-test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
	if !(got.CI.Low < got.Mean && got.CI.High > got.Mean) {
		t.Fatal("expected CI to contain sample mean")
	}
}

func TestTTestTwoSamplePooled(t *testing.T) {
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

	// mean diff = 1, pooled variance = 10
	// se = sqrt(10*(1/5+1/5)) = sqrt(4) = 2
	// t = 0.5
	wantT := 0.5
	assertFloatClose(t, got.Statistic, wantT, 1e-12)

	if got.DF != 8 {
		t.Fatalf("got df %.0f want 8", got.DF)
	}
	if got.Method != "Two-sample t-test (pooled)" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.CI.Low < got.MeanDiff && got.CI.High > got.MeanDiff) {
		t.Fatal("expected CI to contain mean difference")
	}
}

func TestTTestTwoSampleWelch(t *testing.T) {
	x := []float64{10, 12, 14, 16}
	y := []float64{8, 9, 10, 11}

	got, err := TTestTwoSample(x, y, TwoSampleTOptions{
		EqualVariance:   false,
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Method != "Two-sample t-test (Welch)" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if !(got.DF > 0) {
		t.Fatalf("expected positive df, got %v", got.DF)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
}

func TestTTestPaired(t *testing.T) {
	x := []float64{10, 12, 14, 16}
	y := []float64{9, 11, 13, 15}

	_, err := TTestPaired(x, y, PairedTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err == nil {
		t.Fatal("expected zero variance error")
	}
}

func TestTTestPairedNonZeroVariance(t *testing.T) {
	x := []float64{10, 12, 14, 16}
	y := []float64{9, 10, 13, 17}

	got, err := TTestPaired(x, y, PairedTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Method != "Paired t-test" {
		t.Fatalf("unexpected method: %s", got.Method)
	}
	if got.DF != 3 {
		t.Fatalf("got df %.0f want 3", got.DF)
	}
	if !(got.PValue > 0 && got.PValue < 1) {
		t.Fatalf("unexpected p-value: %v", got.PValue)
	}
	if !(got.CI.Low < got.MeanDiff && got.CI.High > got.MeanDiff) {
		t.Fatal("expected CI to contain mean paired difference")
	}
}

func TestTTestPairedMismatchedLengths(t *testing.T) {
	_, err := TTestPaired(
		[]float64{1, 2, 3},
		[]float64{1, 2},
		PairedTOptions{Alternative: TwoSided, ConfidenceLevel: 0.95},
	)
	if err == nil {
		t.Fatal("expected mismatched lengths error")
	}
}

func TestTTestOneSampleInvalidAlternative(t *testing.T) {
	_, err := TTestOneSample([]float64{1, 2, 3}, 0, Alternative("weird"), 0.95)
	if err == nil {
		t.Fatal("expected invalid alternative error")
	}
}

func TestTTestTwoSampleInvalidConfidence(t *testing.T) {
	_, err := TTestTwoSample([]float64{1, 2}, []float64{3, 4}, TwoSampleTOptions{
		EqualVariance:   true,
		Alternative:     TwoSided,
		ConfidenceLevel: 1.0,
	})
	if err == nil {
		t.Fatal("expected invalid confidence error")
	}
}

func TestTTestPairedZeroVariance(t *testing.T) {
	_, err := TTestPaired(
		[]float64{10, 12, 14, 16},
		[]float64{9, 11, 13, 15},
		PairedTOptions{Alternative: TwoSided, ConfidenceLevel: 0.95},
	)
	if err == nil {
		t.Fatal("expected zero variance error")
	}
}

func TestTPValueDirectional(t *testing.T) {
	df := 10.0
	tstat := 2.0

	pTwo := tPValue(tstat, df, TwoSided)
	pGreater := tPValue(tstat, df, Greater)
	pLess := tPValue(tstat, df, Less)

	if math.Abs(pTwo-2*pGreater) > 1e-12 {
		t.Fatalf("expected two-sided p to be twice greater-tail p, got %v vs %v", pTwo, pGreater)
	}
	if !(pLess > 0.5) {
		t.Fatalf("expected less-tail p to be > 0.5 for positive t, got %v", pLess)
	}
}
