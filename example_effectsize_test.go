package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleCohensD() {
	treatment := []float64{5, 6, 7, 8, 9}
	control := []float64{3, 4, 5, 6, 7}

	d, err := statstest.CohensD(treatment, control, true)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("d = %.3f\n", d)

	// Output:
	// d = 1.265
}

func ExampleHedgesG() {
	treatment := []float64{5, 6, 7, 8, 9}
	control := []float64{3, 4, 5, 6, 7}

	g, err := statstest.HedgesG(treatment, control)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("g = %.3f\n", g)

	// Output:
	// g = 1.143
}
