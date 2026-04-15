package statstest

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

func assertFloatClose(t *testing.T, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Fatalf("got %.12f want %.12f", got, want)
	}
}

func assertFloatSlicesClose(t *testing.T, got, want []float64, tol float64) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d want %d", len(got), len(want))
	}

	for i := range got {
		if math.Abs(got[i]-want[i]) > tol {
			t.Fatalf("at index %d: got %.12f want %.12f", i, got[i], want[i])
		}
	}
}
func TestValidateSample(t *testing.T) {
	tests := []struct {
		name   string
		x      []float64
		minLen int
		want   error
	}{
		{
			name:   "empty sample",
			x:      []float64{},
			minLen: 1,
			want:   ErrEmptySample,
		},
		{
			name:   "too small",
			x:      []float64{1.0},
			minLen: 2,
			want:   ErrSampleTooSmall,
		},
		{
			name:   "contains nan",
			x:      []float64{1.0, math.NaN()},
			minLen: 2,
			want:   ErrContainsNaN,
		},
		{
			name:   "contains inf",
			x:      []float64{1.0, math.Inf(1)},
			minLen: 2,
			want:   ErrContainsInf,
		},
		{
			name:   "valid sample",
			x:      []float64{1.0, 2.0, 3.0},
			minLen: 2,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateSample(tt.x, tt.minLen)
			if got != tt.want {
				t.Fatalf("got error %v want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAlternative(t *testing.T) {
	tests := []struct {
		name string
		alt  Alternative
		want error
	}{
		{name: "two sided", alt: TwoSided, want: nil},
		{name: "less", alt: Less, want: nil},
		{name: "greater", alt: Greater, want: nil},
		{name: "invalid", alt: Alternative("bad"), want: ErrInvalidAlternative},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateAlternative(tt.alt)
			if got != tt.want {
				t.Fatalf("got error %v want %v", got, tt.want)
			}
		})
	}
}

func TestValidateConfidence(t *testing.T) {
	tests := []struct {
		name string
		conf float64
		want error
	}{
		{name: "zero", conf: 0, want: ErrInvalidConfidenceLevel},
		{name: "one", conf: 1, want: ErrInvalidConfidenceLevel},
		{name: "negative", conf: -0.5, want: ErrInvalidConfidenceLevel},
		{name: "greater than one", conf: 1.2, want: ErrInvalidConfidenceLevel},
		{name: "valid", conf: 0.95, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateConfidence(tt.conf)
			if got != tt.want {
				t.Fatalf("got error %v want %v", got, tt.want)
			}
		})
	}
}

func TestMean(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	got := mean(x)
	want := 2.5
	assertFloatClose(t, got, want, 1e-12)
}

func TestSampleVariance(t *testing.T) {
	t.Run("len less than 2", func(t *testing.T) {
		got := sampleVariance([]float64{5})
		want := 0.0
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("normal case", func(t *testing.T) {
		x := []float64{1, 2, 3}
		// sample variance = ((1-2)^2 + (2-2)^2 + (3-2)^2) / (3-1) = 1
		got := sampleVariance(x)
		want := 1.0
		assertFloatClose(t, got, want, 1e-12)
	})
}

func TestSampleStdDev(t *testing.T) {
	t.Run("len less than 2", func(t *testing.T) {
		got := sampleStdDev([]float64{5})
		want := 0.0
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("normal case", func(t *testing.T) {
		x := []float64{1, 2, 3}
		got := sampleStdDev(x)
		want := 1.0
		assertFloatClose(t, got, want, 1e-12)
	})
}

func TestPooledVariance(t *testing.T) {
	t.Run("sample too small", func(t *testing.T) {
		got := pooledVariance([]float64{1}, []float64{2, 3})
		want := 0.0
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("normal case", func(t *testing.T) {
		x := []float64{1, 2, 3}
		y := []float64{2, 3, 4}

		// both sample variances are 1, so pooled variance should also be 1
		got := pooledVariance(x, y)
		want := 1.0
		assertFloatClose(t, got, want, 1e-12)
	})
}

func TestTPValue(t *testing.T) {
	df := 10.0
	dist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    df,
	}

	t.Run("two sided", func(t *testing.T) {
		stat := 1.5
		got := tPValue(stat, df, TwoSided)
		want := 2 * (1 - dist.CDF(math.Abs(stat)))
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("less", func(t *testing.T) {
		stat := 1.5
		got := tPValue(stat, df, Less)
		want := dist.CDF(stat)
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("greater", func(t *testing.T) {
		stat := 1.5
		got := tPValue(stat, df, Greater)
		want := 1 - dist.CDF(stat)
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("two sided negative equals positive", func(t *testing.T) {
		gotNeg := tPValue(-2.0, df, TwoSided)
		gotPos := tPValue(2.0, df, TwoSided)
		assertFloatClose(t, gotNeg, gotPos, 1e-12)
	})

	t.Run("invalid alternative returns nan", func(t *testing.T) {
		got := tPValue(1.0, df, Alternative("bad"))
		if !math.IsNaN(got) {
			t.Fatalf("expected NaN, got %v", got)
		}
	})
}

func TestTCriticalTwoSided(t *testing.T) {
	conf := 0.95
	df := 10.0

	dist := distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
		Nu:    df,
	}
	want := dist.Quantile(1 - (1-conf)/2)

	got := tCriticalTwoSided(conf, df)
	assertFloatClose(t, got, want, 1e-12)
}

func TestNormalPValue(t *testing.T) {
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	t.Run("two sided", func(t *testing.T) {
		z := 1.5
		got := normalPValue(z, TwoSided)
		want := 2 * (1 - dist.CDF(math.Abs(z)))
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("less", func(t *testing.T) {
		z := 1.5
		got := normalPValue(z, Less)
		want := dist.CDF(z)
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("greater", func(t *testing.T) {
		z := 1.5
		got := normalPValue(z, Greater)
		want := 1 - dist.CDF(z)
		assertFloatClose(t, got, want, 1e-12)
	})

	t.Run("two sided negative equals positive", func(t *testing.T) {
		gotNeg := normalPValue(-2.0, TwoSided)
		gotPos := normalPValue(2.0, TwoSided)
		assertFloatClose(t, gotNeg, gotPos, 1e-12)
	})

	t.Run("invalid alternative returns nan", func(t *testing.T) {
		got := normalPValue(1.0, Alternative("bad"))
		if !math.IsNaN(got) {
			t.Fatalf("expected NaN, got %v", got)
		}
	})
}
