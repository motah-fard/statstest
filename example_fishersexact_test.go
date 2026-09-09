package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleFishersExact2x2() {
	table := [2][2]int{
		{8, 2},
		{1, 5},
	}

	res, err := statstest.FishersExact2x2(table, statstest.TwoSided)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("oddsRatio = %.3f, p = %.4f\n", res.OddsRatio, res.PValue)

	// Output:
	// oddsRatio = 20.000, p = 0.0350
}
