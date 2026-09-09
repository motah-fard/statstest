package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleChiSquareGoodnessOfFit() {
	// nil expected counts means "uniform across categories".
	observed := []float64{18, 22, 16, 24, 20}

	res, err := statstest.ChiSquareGoodnessOfFit(observed, nil)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("chi2 = %.3f, p = %.4f, df = %d\n", res.Statistic, res.PValue, res.DF)

	// Output:
	// chi2 = 2.000, p = 0.7358, df = 4
}
