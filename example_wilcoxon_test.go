package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleWilcoxonSignedRank() {
	before := []float64{120, 118, 121, 119, 117, 122, 116, 123, 118, 120}
	after := []float64{115, 117, 119, 118, 116, 120, 114, 121, 117, 119}

	res, err := statstest.WilcoxonSignedRank(before, after, statstest.TwoSided)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("W = %.1f, p = %.4f\n", res.W, res.PValue)

	// Output:
	// W = 55.0, p = 0.0042
}
