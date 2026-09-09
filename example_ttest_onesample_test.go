package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleTTestOneSample() {
	x := []float64{2.3, 2.5, 2.1, 2.7, 2.4}

	res, err := statstest.TTestOneSample(x, 2.0, statstest.TwoSided, 0.95)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("t = %.3f, p = %.3f, CI = [%.3f, %.3f]\n", res.Statistic, res.PValue, res.CI.Low, res.CI.High)

	// Output:
	// t = 4.000, p = 0.016, CI = [2.122, 2.678]
}
