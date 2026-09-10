package statstest_test

import (
	"fmt"
	"sort"

	"github.com/motah-fard/statstest"
)

func ExampleBootstrapCI() {
	x := []float64{
		12.1, 13.4, 11.9, 14.0, 12.8, 13.1, 12.5, 14.2, 11.7, 13.9,
		12.3, 13.6, 12.0, 14.1, 12.9, 13.3, 12.6, 13.8, 12.2, 13.5,
	}

	median := func(s []float64) float64 {
		sorted := append([]float64(nil), s...)
		sort.Float64s(sorted)
		n := len(sorted)
		if n%2 == 1 {
			return sorted[n/2]
		}
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}

	seed := int64(42) // fixed for a reproducible doc example; omit in real use
	res, err := statstest.BootstrapCI(x, median, statstest.BootstrapOptions{
		NumResamples:    2000,
		ConfidenceLevel: 0.95,
		Seed:            &seed,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("median = %.2f, 95%% CI = [%.2f, %.2f]\n", res.Estimate, res.CI.Low, res.CI.High)

	// Output:
	// median = 13.00, 95% CI = [12.40, 13.55]
}
