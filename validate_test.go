package statstest

import (
	"errors"
	"math"
	"testing"
)

func TestTTestOneSampleRejectsNaN(t *testing.T) {
	_, err := TTestOneSample([]float64{1, 2, math.NaN()}, 0, TwoSided, 0.95)
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}

func TestTTestOneSampleRejectsInf(t *testing.T) {
	_, err := TTestOneSample([]float64{1, 2, math.Inf(1)}, 0, TwoSided, 0.95)
	if !errors.Is(err, ErrContainsInf) {
		t.Fatalf("expected ErrContainsInf, got %v", err)
	}
}

func TestTTestOneSampleRejectsEmpty(t *testing.T) {
	_, err := TTestOneSample([]float64{}, 0, TwoSided, 0.95)
	if !errors.Is(err, ErrEmptySample) {
		t.Fatalf("expected ErrEmptySample, got %v", err)
	}
}

func TestTTestOneSampleRejectsTooSmall(t *testing.T) {
	_, err := TTestOneSample([]float64{1}, 0, TwoSided, 0.95)
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestTTestOneSampleRejectsInvalidAlternative(t *testing.T) {
	_, err := TTestOneSample([]float64{1, 2, 3}, 0, Alternative("bad"), 0.95)
	if !errors.Is(err, ErrInvalidAlternative) {
		t.Fatalf("expected ErrInvalidAlternative, got %v", err)
	}
}

func TestTTestOneSampleConfidenceValidation(t *testing.T) {
	tests := []struct {
		name      string
		confLevel float64
		wantErr   error
	}{
		{"valid 0.95", 0.95, nil},
		{"reject zero", 0, ErrInvalidConfidenceLevel},
		{"reject one", 1, ErrInvalidConfidenceLevel},
		{"reject negative", -0.5, ErrInvalidConfidenceLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TTestOneSample([]float64{1, 2, 3}, 0, TwoSided, tt.confLevel)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected nil, got %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
func TestTTestTwoSampleRejectsInf(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{1, 2, 3},
		[]float64{4, 5, math.Inf(1)},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     TwoSided,
			ConfidenceLevel: 0.95,
		},
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTTestTwoSampleRejectsNaN(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{1, 2, math.NaN()},
		[]float64{2, 3, 4},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     TwoSided,
			ConfidenceLevel: 0.95,
		},
	)
	if !errors.Is(err, ErrContainsNaN) {
		t.Fatalf("expected ErrContainsNaN, got %v", err)
	}
}

func TestTTestTwoSampleRejectsEmpty(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{},
		[]float64{1, 2},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     TwoSided,
			ConfidenceLevel: 0.95,
		},
	)
	if !errors.Is(err, ErrEmptySample) {
		t.Fatalf("expected ErrEmptySample, got %v", err)
	}
}

func TestTTestTwoSampleRejectsTooSmall(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{1},
		[]float64{2},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     TwoSided,
			ConfidenceLevel: 0.95,
		},
	)
	if !errors.Is(err, ErrSampleTooSmall) {
		t.Fatalf("expected ErrSampleTooSmall, got %v", err)
	}
}

func TestTTestTwoSampleRejectsInvalidAlternative(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{1, 2},
		[]float64{3, 4},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     Alternative("bad"),
			ConfidenceLevel: 0.95,
		},
	)
	if !errors.Is(err, ErrInvalidAlternative) {
		t.Fatalf("expected ErrInvalidAlternative, got %v", err)
	}
}

func TestTTestTwoSampleRejectsZeroVariance(t *testing.T) {
	_, err := TTestTwoSample(
		[]float64{5, 5, 5},
		[]float64{5, 5, 5},
		TwoSampleTOptions{
			EqualVariance:   false,
			Alternative:     TwoSided,
			ConfidenceLevel: 0.95,
		},
	)
	if !errors.Is(err, ErrZeroVariance) {
		t.Fatalf("expected ErrZeroVariance, got %v", err)
	}
}
