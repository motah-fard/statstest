package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleOneWayANOVA() {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2}
	g2 := []float64{8.8, 9.0, 8.7, 8.9, 9.1}
	g3 := []float64{7.5, 7.6, 7.4, 7.7, 7.5}

	res, err := statstest.OneWayANOVA(g1, g2, g3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("F = %.3f, p = %.2e, etaSq = %.3f\n", res.FStatistic, res.PValue, res.EtaSquared)

	// Output:
	// F = 111.238, p = 1.80e-08, etaSq = 0.949
}
