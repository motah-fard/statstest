package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExamplePearsonCorrelation() {
	x := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7}
	y := []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82}

	res, err := statstest.PearsonCorrelation(x, y, statstest.TwoSided, 0.95)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("r = %.3f, p = %.4f\n", res.R, res.PValue)

	// Output:
	// r = 0.797, p = 0.0058
}

func ExampleSpearmanCorrelation() {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	y := []float64{2, 1, 4, 3, 7, 5, 6, 9, 8, 10}

	res, err := statstest.SpearmanCorrelation(x, y, statstest.TwoSided, 0.95)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("rho = %.3f, p = %.5f\n", res.R, res.PValue)

	// Output:
	// rho = 0.927, p = 0.00011
}
