package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

// Example shows the shape of a typical statstest call: pass plain
// []float64 samples, get back a structured result with the statistic,
// p-value, and confidence interval instead of a bare tuple.
func Example() {
	treatment := []float64{12.1, 13.4, 11.9, 14.0, 12.8}
	control := []float64{10.2, 10.8, 11.3, 9.9, 10.5}

	res, err := statstest.TTestTwoSample(treatment, control, statstest.TwoSampleTOptions{
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
