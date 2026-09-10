package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleSampleSizeTTestTwoSample() {
	// How many observations per group do I need to reliably detect a
	// medium effect size (d=0.5) at the 95% significance level with 80%
	// power?
	n, err := statstest.SampleSizeTTestTwoSample(0.5, 0.05, 0.8)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("n per group = %.1f\n", n)

	// Output:
	// n per group = 62.8
}

func ExamplePowerTTestTwoSample() {
	// With 64 observations per group, how likely am I to detect a medium
	// effect size (d=0.5) at the 95% significance level?
	power, err := statstest.PowerTTestTwoSample(64, 0.5, 0.05)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("power = %.3f\n", power)

	// Output:
	// power = 0.807
}

func ExampleSampleSizeProportionTwoSample() {
	// How many users per variant does an A/B test need to reliably detect
	// a conversion-rate lift from 30% to 50%?
	n, err := statstest.SampleSizeProportionTwoSample(0.5, 0.3, 0.05, 0.8)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("n per variant = %.1f\n", n)

	// Output:
	// n per variant = 92.7
}
