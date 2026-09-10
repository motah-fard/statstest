package statstest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

// These benchmarks exist to catch performance regressions and to make an
// honest, reproducible claim about where this package's costs actually
// are. Run them with:
//
//	go test -bench=. -benchmem ./...
//
// TukeyHSD is the one function in the package expensive enough (tens to
// hundreds of milliseconds at larger group counts) that its parallel and
// context-aware variants are worth benchmarking on their own — see
// BenchmarkTukeyHSD and BenchmarkTukeyHSDContext.

func benchNormalSample(rng *rand.Rand, n int, mean float64) []float64 {
	x := make([]float64, n)
	for i := range x {
		x[i] = rng.NormFloat64() + mean
	}
	return x
}

func benchNormalGroups(rng *rand.Rand, k, n int) [][]float64 {
	groups := make([][]float64, k)
	for i := range groups {
		groups[i] = benchNormalSample(rng, n, float64(i)*0.3)
	}
	return groups
}

func BenchmarkTTestTwoSample(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	x := benchNormalSample(rng, 100, 0)
	y := benchNormalSample(rng, 100, 0.5)
	opts := TwoSampleTOptions{Alternative: TwoSided, ConfidenceLevel: 0.95}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TTestTwoSample(x, y, opts)
	}
}

func BenchmarkMannWhitneyU(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	x := benchNormalSample(rng, 100, 0)
	y := benchNormalSample(rng, 100, 0.5)
	opts := MannWhitneyOptions{Alternative: TwoSided}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MannWhitneyU(x, y, opts)
	}
}

func BenchmarkOneWayANOVA(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	groups := benchNormalGroups(rng, 5, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = OneWayANOVA(groups...)
	}
}

func BenchmarkKruskalWallis(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	groups := benchNormalGroups(rng, 5, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = KruskalWallis(groups...)
	}
}

func BenchmarkPearsonCorrelation(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	x := benchNormalSample(rng, 200, 0)
	y := benchNormalSample(rng, 200, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PearsonCorrelation(x, y, TwoSided, 0.95)
	}
}

func BenchmarkChiSquareIndependence(b *testing.B) {
	table := [][]int{
		{90, 60, 104, 95},
		{30, 50, 51, 20},
		{30, 40, 45, 35},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ChiSquareIndependence(table)
	}
}

func BenchmarkProportionTwoSample(b *testing.B) {
	opts := ProportionTwoSampleOptions{Alternative: TwoSided, ConfidenceLevel: 0.95}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ProportionTwoSample(450, 1000, 300, 900, opts)
	}
}

func BenchmarkAdjustPValues(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{100, 10000} {
		p := make([]float64, n)
		for i := range p {
			p[i] = rng.Float64()
		}
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = AdjustPValues(p, BenjaminiHochberg)
			}
		})
	}
}

func BenchmarkShapiroWilk(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{30, 500, 5000} {
		x := benchNormalSample(rng, n, 0)
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ShapiroWilk(x)
			}
		})
	}
}

// BenchmarkTukeyHSD covers a range of group counts because its pairwise
// comparison loop is the one place in the package where cost grows
// noticeably (roughly quadratically in the number of groups) and where
// parallelComputeEachContext's goroutine fan-out actually matters.
func BenchmarkTukeyHSD(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	for _, k := range []int{3, 10, 30} {
		groups := benchNormalGroups(rng, k, 20)
		b.Run(fmt.Sprintf("k=%d", k), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = TukeyHSD(0.95, groups...)
			}
		})
	}
}

// BenchmarkTukeyHSDContext measures the overhead of the context-aware
// variant against BenchmarkTukeyHSD's k=10 case, with a context that is
// never canceled — the two should be close, since TukeyHSD is a thin
// wrapper around TukeyHSDContext.
func BenchmarkTukeyHSDContext(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	groups := benchNormalGroups(rng, 10, 20)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TukeyHSDContext(ctx, 0.95, groups...)
	}
}

func BenchmarkDunnTest(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	for _, k := range []int{3, 20} {
		groups := benchNormalGroups(rng, k, 20)
		b.Run(fmt.Sprintf("k=%d", k), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = DunnTest(Bonferroni, groups...)
			}
		})
	}
}
