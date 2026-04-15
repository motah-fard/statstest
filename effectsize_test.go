package statstest

import (
	"math"
	"testing"
)

func TestCohensDPooled(t *testing.T) {
	x := []float64{2, 4, 6, 8, 10}
	y := []float64{1, 3, 5, 7, 9}

	got, err := CohensD(x, y, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// means differ by 1, pooled sd = sqrt(10/1?) no, let's compute exactly:
	// sample variance for each = 10
	// pooled sd = sqrt(10) = 3.1622776601683795
	want := 1.0 / math.Sqrt(10)

	assertFloatClose(t, got, want, 1e-12)
}

func TestCohensDUnpooled(t *testing.T) {
	x := []float64{10, 12, 14, 16}
	y := []float64{8, 9, 10, 11}

	got, err := CohensD(x, y, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	meanX := 13.0
	meanY := 9.5
	diff := meanX - meanY

	varX := 20.0 / 3.0
	varY := 5.0 / 3.0
	want := diff / math.Sqrt((varX+varY)/2.0)

	assertFloatClose(t, got, want, 1e-12)
}

func TestHedgesG(t *testing.T) {
	x := []float64{2, 4, 6, 8, 10}
	y := []float64{1, 3, 5, 7, 9}

	got, err := HedgesG(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d := 1.0 / math.Sqrt(10)
	df := float64(len(x) + len(y) - 2)
	j := 1.0 - (3.0 / (4.0*df - 1.0))
	want := j * d

	assertFloatClose(t, got, want, 1e-12)
}

func TestCohensDErrorTooSmall(t *testing.T) {
	_, err := CohensD([]float64{1}, []float64{1, 2}, true)
	if err == nil {
		t.Fatal("expected error for small sample")
	}
}

func TestCohensDErrorZeroVariance(t *testing.T) {
	_, err := CohensD([]float64{1, 1, 1}, []float64{2, 2, 2}, true)
	if err == nil {
		t.Fatal("expected zero variance error")
	}
}

func TestHedgesGErrorZeroVariance(t *testing.T) {
	_, err := HedgesG([]float64{1, 1, 1}, []float64{2, 2, 2})
	if err == nil {
		t.Fatal("expected zero variance error")
	}
}

// func assertFloatClose(t *testing.T, got, want, tol float64) {
// 	t.Helper()
// 	if math.Abs(got-want) > tol {
// 		t.Fatalf("got %.12f want %.12f", got, want)
// 	}
// }
