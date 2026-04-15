package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleMannWhitneyU() {
	x := []float64{1, 2, 3}
	y := []float64{10, 11, 12}

	res, err := statstest.MannWhitneyU(x, y, statstest.MannWhitneyOptions{
		Alternative:             statstest.Less,
		UseContinuityCorrection: false,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("U = %.1f\n", res.U)

	// Output:
	// U = 0.0
}
