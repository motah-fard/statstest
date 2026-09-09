package statstest

import (
	"math"
	"math/rand"
	"testing"
)

// These fuzz tests target the numerically involved code in this package —
// Royston's Shapiro-Wilk approximation and the from-scratch, numerically
// integrated studentized range distribution behind TukeyHSD — where a
// stray edge case (near-duplicate values, extreme sample sizes, unusual
// tie patterns) is more likely to produce a NaN, an out-of-range
// probability, or a panic than in the closed-form tests elsewhere in the
// package.
//
// Run with, e.g.: go test -fuzz=FuzzShapiroWilk -fuzztime=30s

func FuzzShapiroWilk(f *testing.F) {
	f.Add(3, int64(1))
	f.Add(5, int64(2))
	f.Add(11, int64(3))
	f.Add(12, int64(4))
	f.Add(50, int64(5))
	f.Add(500, int64(6))

	f.Fuzz(func(t *testing.T, n int, seed int64) {
		if n < 3 || n > 500 {
			t.Skip()
		}

		rng := rand.New(rand.NewSource(seed))
		x := make([]float64, n)
		for i := range x {
			x[i] = rng.NormFloat64() * 10
		}

		res, err := ShapiroWilk(x)
		if err != nil {
			return
		}

		if math.IsNaN(res.W) || math.IsInf(res.W, 0) {
			t.Fatalf("W is not finite: %v (n=%d seed=%d)", res.W, n, seed)
		}
		if res.W < 0 || res.W > 1.0000001 {
			t.Fatalf("W out of expected range: %v (n=%d seed=%d)", res.W, n, seed)
		}
		if math.IsNaN(res.PValue) || res.PValue < 0 || res.PValue > 1 {
			t.Fatalf("p-value out of range: %v (n=%d seed=%d)", res.PValue, n, seed)
		}
	})
}

func FuzzStudentizedRangeCDF(f *testing.F) {
	f.Add(0.5, 2, 1.0)
	f.Add(3.0, 3, 10.0)
	f.Add(10.0, 8, 1.0)
	f.Add(0.001, 5, 500.0)

	f.Fuzz(func(t *testing.T, q float64, k int, df float64) {
		if math.IsNaN(q) || math.IsInf(q, 0) || q < 0 || q > 50 {
			t.Skip()
		}
		if k < 2 || k > 20 {
			t.Skip()
		}
		if math.IsNaN(df) || math.IsInf(df, 0) || df < 1 || df > 1000 {
			t.Skip()
		}

		got := studentizedRangeCDF(q, k, df)
		if math.IsNaN(got) {
			t.Fatalf("CDF is NaN: q=%v k=%d df=%v", q, k, df)
		}
		if got < 0 || got > 1 {
			t.Fatalf("CDF out of [0,1]: got %v for q=%v k=%d df=%v", got, q, k, df)
		}
	})
}

func FuzzTukeyHSD(f *testing.F) {
	f.Add(2, 3, int64(1))
	f.Add(5, 10, int64(2))
	f.Add(8, 2, int64(3))

	f.Fuzz(func(t *testing.T, k, n int, seed int64) {
		if k < 2 || k > 8 || n < 2 || n > 50 {
			t.Skip()
		}

		rng := rand.New(rand.NewSource(seed))
		groups := make([][]float64, k)
		for i := range groups {
			g := make([]float64, n)
			for j := range g {
				g[j] = rng.NormFloat64()*5 + float64(i)
			}
			groups[i] = g
		}

		res, err := TukeyHSD(0.95, groups...)
		if err != nil {
			return
		}

		for _, c := range res.Comparisons {
			if math.IsNaN(c.PValue) || c.PValue < 0 || c.PValue > 1 {
				t.Fatalf("p-value out of range: %v for (%d,%d) k=%d n=%d seed=%d", c.PValue, c.GroupI, c.GroupJ, k, n, seed)
			}
			if math.IsNaN(c.CI.Low) || math.IsNaN(c.CI.High) {
				t.Fatalf("CI contains NaN for (%d,%d) k=%d n=%d seed=%d", c.GroupI, c.GroupJ, k, n, seed)
			}
			if c.CI.Low > c.CI.High {
				t.Fatalf("CI.Low > CI.High for (%d,%d): [%v, %v]", c.GroupI, c.GroupJ, c.CI.Low, c.CI.High)
			}
		}
	})
}

func FuzzDunnTest(f *testing.F) {
	f.Add(2, 4, int64(1))
	f.Add(5, 8, int64(2))

	f.Fuzz(func(t *testing.T, k, n int, seed int64) {
		if k < 2 || k > 8 || n < 1 || n > 50 {
			t.Skip()
		}

		rng := rand.New(rand.NewSource(seed))
		groups := make([][]float64, k)
		for i := range groups {
			g := make([]float64, n)
			for j := range g {
				// Bias toward small integers so repeated fuzz runs exercise
				// the tie-correction path, not just distinct-valued data.
				g[j] = float64(rng.Intn(5))
			}
			groups[i] = g
		}

		res, err := DunnTest(Bonferroni, groups...)
		if err != nil {
			return
		}

		for _, c := range res.Comparisons {
			if math.IsNaN(c.PValue) || c.PValue < 0 || c.PValue > 1 {
				t.Fatalf("p-value out of range: %v for (%d,%d) k=%d n=%d seed=%d", c.PValue, c.GroupI, c.GroupJ, k, n, seed)
			}
			if math.IsNaN(c.AdjustedPValue) || c.AdjustedPValue < 0 || c.AdjustedPValue > 1 {
				t.Fatalf("adjusted p-value out of range: %v for (%d,%d) k=%d n=%d seed=%d", c.AdjustedPValue, c.GroupI, c.GroupJ, k, n, seed)
			}
		}
	})
}
