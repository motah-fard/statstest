package statstest_test

import (
	"fmt"

	"github.com/motah-fard/statstest"
)

func ExampleDunnTest() {
	g1 := []float64{2.9, 3.0, 2.5, 2.6, 3.2}
	g2 := []float64{3.8, 2.7, 4.0, 2.4, 3.9}
	g3 := []float64{2.8, 3.4, 3.7, 2.2, 2.0}

	res, err := statstest.DunnTest(statstest.Bonferroni, g1, g2, g3)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	c := res.Comparisons[0]
	fmt.Printf("group %d vs %d: z = %.3f, adjusted p = %.4f\n", c.GroupI, c.GroupJ, c.Statistic, c.AdjustedPValue)

	// Output:
	// group 0 vs 1: z = -1.061, adjusted p = 0.8665
}
