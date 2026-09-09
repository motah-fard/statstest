package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleProportionOneSample() {
	res, err := statstest.ProportionOneSample(45, 100, 0.5, statstest.ProportionOneSampleOptions{
		Alternative:     statstest.TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("z = %.3f, p = %.4f\n", res.Statistic, res.PValue)

	// Output:
	// z = -1.000, p = 0.3173
}

func ExampleProportionTwoSample() {
	res, err := statstest.ProportionTwoSample(45, 100, 30, 90, statstest.ProportionTwoSampleOptions{
		Alternative:     statstest.TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("z = %.3f, p = %.4f, diff = %.3f\n", res.Statistic, res.PValue, res.Diff)

	// Output:
	// z = 1.643, p = 0.1004, diff = 0.117
}
