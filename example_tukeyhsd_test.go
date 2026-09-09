package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleTukeyHSD() {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2}
	g2 := []float64{8.8, 9.0, 8.7, 8.9, 9.1}
	g3 := []float64{7.5, 7.6, 7.4, 7.7, 7.5}

	res, err := statstest.TukeyHSD(0.95, g1, g2, g3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// res.Comparisons holds one entry per pair of group indices; the first
	// pair compares group 0 (g1) against group 1 (g2).
	c := res.Comparisons[0]
	fmt.Printf("group %d vs %d: diff = %.3f, p = %.2e, CI = [%.3f, %.3f]\n",
		c.GroupI, c.GroupJ, c.MeanDiff, c.PValue, c.CI.Low, c.CI.High)

	// Output:
	// group 0 vs 1: diff = -0.800, p = 3.55e-06, CI = [-1.045, -0.555]
}
