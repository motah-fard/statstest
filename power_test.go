package statstest

import (
	"errors"
	"testing"
)

func TestCohensHKnownValue(t *testing.T) {
	h, err := CohensH(0.5, 0.3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFloatClose(t, h, 0.4115168460674883, 1e-9)
}

func TestCohensHRejectsOutOfRangeProportion(t *testing.T) {
	if _, err := CohensH(0, 0.3); !errors.Is(err, ErrInvalidProportion) {
		t.Fatalf("expected ErrInvalidProportion, got %v", err)
	}
	if _, err := CohensH(0.5, 1); !errors.Is(err, ErrInvalidProportion) {
		t.Fatalf("expected ErrInvalidProportion, got %v", err)
	}
}

func TestPowerTTestTwoSampleIncreasesWithN(t *testing.T) {
	small, err := PowerTTestTwoSample(10, 0.5, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	large, err := PowerTTestTwoSample(200, 0.5, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(large > small) {
		t.Fatalf("expected power to increase with n: n=10 -> %v, n=200 -> %v", small, large)
	}
	if !(small > 0 && large < 1) {
		t.Fatalf("expected power in (0,1), got %v and %v", small, large)
	}
}

func TestPowerTTestTwoSampleRejectsZeroEffectSize(t *testing.T) {
	_, err := PowerTTestTwoSample(50, 0, 0.05)
	if !errors.Is(err, ErrZeroEffectSize) {
		t.Fatalf("expected ErrZeroEffectSize, got %v", err)
	}
}

func TestPowerTTestTwoSampleSignIndependent(t *testing.T) {
	pos, err := PowerTTestTwoSample(50, 0.5, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	neg, err := PowerTTestTwoSample(50, -0.5, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertFloatClose(t, pos, neg, 1e-12)
}

func TestSampleSizeTTestTwoSampleRoundTrips(t *testing.T) {
	n, err := SampleSizeTTestTwoSample(0.5, 0.05, 0.8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := PowerTTestTwoSample(int(n+1), 0.5, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(got >= 0.8) {
		t.Fatalf("expected the computed sample size (rounded up) to reach the target power, got %v", got)
	}
}

func TestSampleSizeTTestTwoSampleRejectsInvalidPower(t *testing.T) {
	_, err := SampleSizeTTestTwoSample(0.5, 0.05, 1.5)
	if !errors.Is(err, ErrInvalidConfidenceLevel) {
		t.Fatalf("expected ErrInvalidConfidenceLevel, got %v", err)
	}
}

func TestPowerProportionTwoSampleIncreasesWithN(t *testing.T) {
	small, err := PowerProportionTwoSample(20, 0.5, 0.3, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	large, err := PowerProportionTwoSample(500, 0.5, 0.3, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(large > small) {
		t.Fatalf("expected power to increase with n: n=20 -> %v, n=500 -> %v", small, large)
	}
}

func TestPowerProportionTwoSampleRejectsEqualProportions(t *testing.T) {
	_, err := PowerProportionTwoSample(50, 0.4, 0.4, 0.05)
	if !errors.Is(err, ErrZeroEffectSize) {
		t.Fatalf("expected ErrZeroEffectSize, got %v", err)
	}
}

func TestSampleSizeProportionTwoSampleRoundTrips(t *testing.T) {
	n, err := SampleSizeProportionTwoSample(0.5, 0.3, 0.05, 0.8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := PowerProportionTwoSample(int(n+1), 0.5, 0.3, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(got >= 0.8) {
		t.Fatalf("expected the computed sample size (rounded up) to reach the target power, got %v", got)
	}
}
