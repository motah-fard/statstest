package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleTTestPaired() {
	before := []float64{82, 75, 90, 68, 77, 84}
	after := []float64{86, 80, 92, 71, 82, 88}

	res, err := statstest.TTestPaired(before, after, statstest.PairedTOptions{
		Alternative:     statstest.TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("t = %.3f, p = %.4f, meanDiff = %.3f\n", res.Statistic, res.PValue, res.MeanDiff)

	// Output:
	// t = -8.032, p = 0.0005, meanDiff = -3.833
}
