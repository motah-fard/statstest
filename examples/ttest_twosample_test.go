package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleTTestTwoSample() {
	x := []float64{12.1, 13.4, 11.9, 14.0, 12.8}
	y := []float64{10.2, 10.8, 11.3, 9.9, 10.5}

	res, err := statstest.TTestTwoSample(x, y, statstest.TwoSampleTOptions{
		EqualVariance:   false,
		Alternative:     statstest.TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("t = %.3f, p = %.3f\n", res.Statistic, res.PValue)

	// Output:
	// t = 4.981, p = 0.002
}
