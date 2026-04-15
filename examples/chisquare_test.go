package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleChiSquareIndependence() {
	observed := [][]int{
		{90, 60, 104, 95},
		{30, 50, 51, 20},
		{30, 40, 45, 35},
	}

	res, err := statstest.ChiSquareIndependence(observed)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("X^2 = %.3f, df = %d\n", res.Statistic, res.DF)

	// Output:
	// X^2 = 24.571, df = 6
}
